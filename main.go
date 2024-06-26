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

var version = "development"

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
	Profiling        bool
}

func ParseConfig(fileName string) (Configuration, error) {
	file, _ := os.Open(fileName)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	return configuration, err
}

func Router(client mqtt.Client, profiling bool) error {
	err := api.SetupRouter(gin.ReleaseMode, profiling).Run("0.0.0.0:8080")
	if err != nil {
		pubErr := mqtt.Error(client, 0, true, err, 10)
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
	log.Printf("Phocus Version: %s\n\n", version)

	configuration, err := ParseConfig("config.json")

	// TODO log some other useful info here

	// just give a chance to see the http server coming up
	time.Sleep(3 * time.Second)

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
	pubErr := mqtt.Send(client, "phocus/stats/error", 0, true, "", 10)
	if pubErr != nil {
		log.Printf("Failed to clear previous error: %v\n", pubErr)
	}

	// send new version
	pubErr = mqtt.Send(client, "phocus/stats/version", 0, true, version, 10)
	if pubErr != nil {
		log.Printf("Failed to set phocus version: %v\n", pubErr)
	}

	// serial
	port, err := serial.Setup(
		configuration.Serial.Port,
		configuration.Serial.Baud,
		configuration.Serial.Retries,
	)
	if err != nil {
		pubErr := mqtt.Error(client, 0, true, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Printf("Failed to set up serial with err: %v", err)
		os.Exit(1)
	}
	defer port.Port.Close()

	// spawns a go-routine which handles web requests
	go Router(client, configuration.Profiling)

	// sensors
	// we only add them once we know the mqtt, serial and http aspects are up
	err = sensors.Register(client, version)
	if err != nil {
		pubErr := mqtt.Error(client, 0, true, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Printf("Failed to set up sensors with err: %v", err)
		os.Exit(1)
	}

	// sleep to make sure web server comes on before polling starts
	time.Sleep(2 * time.Second)

	// spawn go-routine to repeatedly enQueue QPGSn commands
	go api.QueueQPGSn(configuration.DelaySeconds, configuration.RandDelaySeconds)

	// loop to check Queue and deQueue index 0, run it process result and wait 30 seconds
	for {
		api.QueueMutex.Lock()
		// if there is an entry at [0] then run that command
		if len(api.Queue) > 0 {
			QPGSnResponse, err := messages.Interpret(client, port, api.Queue[0], time.Duration(configuration.Messages.Read.TimeoutSeconds)*time.Second)
			if err != nil {
				pubErr := mqtt.Error(client, 0, true, err, 10)
				if pubErr != nil {
					log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
				}
				if fmt.Sprint(err) == "read returned nothing" { // immediately jailed when read timeout
					port.Port.Close()
					pubErr := mqtt.Error(client, 0, true, errors.New("read timed out, waiting 2 minutes then restarting"), 10)
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
