// Package `main` contains the entrypoint for the `phocus` system
package main

import (
	"errors"  // creating custom errors
	"fmt"     // string formatting
	"log"     // formatted logging
	"os"      // exiting
	"os/exec" // auto restart
	"time"    // for sleeping

	"encoding/json" // for config reading

	"github.com/gin-gonic/gin"
	api "github.com/wolffshots/phocus/v2/api"           // api setup
	messages "github.com/wolffshots/phocus/v2/messages" // message structures
	mqtt "github.com/wolffshots/phocus/v2/mqtt"         // comms with mqtt broker
	sensors "github.com/wolffshots/phocus/v2/sensors"   // registering common sensors
	serial "github.com/wolffshots/phocus/v2/serial"     // comms with inverter
)

const version = "v2.9.10"

type Configuration struct {
	Serial struct {
		Port    string
		Baud    int
		Retries int
	}
	MQTT struct {
		Host   string
		Port   int
		Client struct {
			Name string
		}
		Retries int
	}
	Messages struct {
		Read struct {
			TimeoutSeconds int
		}
	}
	DelaySeconds     int
	RandDelaySeconds int
	MinDelaySeconds  int
}

func ParseConfig(fileName string) (Configuration, error) {
	file, _ := os.Open(fileName)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	return configuration, err
}

func Router(client mqtt.Client) error {
	err := api.SetupRouter(gin.ReleaseMode).Run("0.0.0.0:8080")
	if err != nil {
		pubErr := mqtt.Error(client, 0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Printf("Failed to run http routine with err: %v", err)
		os.Exit(1)
	}
	return err
}

// main is the entrypoint to the app
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.Println("Starting up phocus")

	configuration, err := ParseConfig("config.json")

	if err != nil {
		log.Printf("Error parsing config: %v", err)
		os.Exit(1)
	}

	// mqtt
	client, err := mqtt.Setup(
		configuration.MQTT.Host,
		configuration.MQTT.Port,
		configuration.MQTT.Retries,
		configuration.MQTT.Client.Name,
	)

	if err != nil {
		log.Printf("Failed to set up mqtt %d times with err: %v", configuration.MQTT.Retries, err)
		os.Exit(1)
	}
	// reset error
	pubErr := mqtt.Send(client, "phocus/stats/error", 0, false, "", 10)
	if pubErr != nil {
		log.Printf("Failed to clear previous error: %v\n", pubErr)
	}

	// serial
	port, err := serial.Setup(
		configuration.Serial.Port,
		configuration.Serial.Baud,
		configuration.Serial.Retries,
	)
	if err != nil {
		pubErr := mqtt.Error(client, 0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Printf("Failed to set up serial with err: %v", err)
		os.Exit(1)
	}
	defer port.Port.Close()

	// spawns a go-routine which handles web requests
	go Router(client)

	// sensors
	// we only add them once we know the mqtt, serial and http aspects are up
	err = sensors.Register(client, version)
	if err != nil {
		pubErr := mqtt.Error(client, 0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Printf("Failed to set up sensors with err: %v", err)
		os.Exit(1)
	}

	// sleep to make sure web server comes on before polling starts
	time.Sleep(5 * time.Second)

	// spawn go-routine to repeatedly enQueue QPGSn commands
	go api.QueueQPGSn(configuration.DelaySeconds, configuration.RandDelaySeconds)

	// loop to check Queue and deQueue index 0, run it process result and wait 30 seconds
	for {
		api.QueueMutex.Lock()
		// if there is an entry at [0] then run that command
		if len(api.Queue) > 0 {
			QPGSnResponse, err := messages.Interpret(client, port, api.Queue[0], time.Duration(configuration.Messages.Read.TimeoutSeconds)*time.Second)
			if err != nil {
				pubErr := mqtt.Error(client, 0, false, err, 10)
				if pubErr != nil {
					log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
				}
				if fmt.Sprint(err) == "read returned nothing" { // immediately jailed when read timeout
					port.Port.Close()
					pubErr := mqtt.Error(client, 0, false, errors.New("read timed out, waiting 2 minutes then restarting"), 10)
					if pubErr != nil {
						log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
					}
					time.Sleep(2 * time.Minute)
					cmd, err := exec.Command("bash", "-c", "sudo service phocus restart").Output()
					// it should die here
					log.Printf("cmd=================>%s\n", cmd)
					if err != nil {
						log.Printf("Error execing cmd: %v", err)
					}
					// if it reaches here at all that implies it didn't restart properly
					os.Exit(1)
				}
			}
			if QPGSnResponse != nil {
				api.SetLast(QPGSnResponse)
			}
			api.Queue = api.Queue[1:]
		} else {
			// min sleep between actual comms with inverter
			time.Sleep(time.Duration(configuration.MinDelaySeconds) * time.Second)
		}
		api.QueueMutex.Unlock()
		// min sleep between Queue checks
		time.Sleep(1 * time.Second)
	}
}
