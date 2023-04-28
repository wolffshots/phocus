// Package `main` contains the entrypoint for the `phocus` system
package main

import (
	"errors"  // creating custom errors
	"fmt"     // string formatting
	"log"     // formatted logging
	"os"      // exiting
	"os/exec" // auto restart
	"time"    // for sleeping

	phocus_api "github.com/wolffshots/phocus/v2/api"           // api setup
	phocus_messages "github.com/wolffshots/phocus/v2/messages" // message structures
	phocus_mqtt "github.com/wolffshots/phocus/v2/mqtt"         // comms with mqtt broker
	phocus_sensors "github.com/wolffshots/phocus/v2/sensors"   // registering common sensors
	phocus_serial "github.com/wolffshots/phocus/v2/serial"     // comms with inverter
)

const version = "v2.1.0"

// main is the entrypoint to the app
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.Println("Starting up phocus")

	// mqtt
	err := phocus_mqtt.Setup("192.168.88.14", "go_phocus_client'")
	if err != nil {
		log.Fatalf("Failed to set up mqtt with err: %v", err)
	}
	// reset error
	pubErr := phocus_mqtt.Send("phocus/stats/error", 0, false, "", 10)
	if pubErr != nil {
		log.Printf("Failed to clear previous error: %v\n", pubErr)
	}

	// serial
	err = phocus_serial.Setup()
	if err != nil {
		pubErr := phocus_mqtt.Error(0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Fatalf("Failed to set up serial with err: %v", err)
	}

	// spawns a go-routine which handles web requests
	go func() {
		err := phocus_api.SetupRouter().Run("localhost:8080")
		if err != nil {
			pubErr := phocus_mqtt.Error(0, false, err, 10)
			if pubErr != nil {
				log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
			}
			log.Fatalf("Failed to run http routine with err: %v", err)
		}
	}()

	// sensors
	// we only add them once we know the mqtt, serial and http aspects are up
	err = phocus_sensors.Register(version)
	if err != nil {
		pubErr := phocus_mqtt.Error(0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Fatalf("Failed to set up sensors with err: %v", err)
	}

	// sleep to make sure web server comes on before polling starts
	time.Sleep(5 * time.Second)

	// spawn go-routine to repeatedly enQueue QPGSn commands
	go phocus_api.QueueQPGSn()

	// loop to check Queue and deQueue index 0, run it process result and wait 30 seconds
	for {
		phocus_api.QueueMutex.Lock()
		log.Print(".")
		// if there is an entry at [0] then run that command
		if len(phocus_api.Queue) > 0 {
			err := phocus_messages.Interpret(phocus_api.Queue[0])
			if err != nil {
				pubErr := phocus_mqtt.Error(0, false, err, 10)
				if pubErr != nil {
					log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
				}
				if fmt.Sprint(err) == "read timed out" { // immediately jailed when read timeout
					pubErr := phocus_mqtt.Error(0, false, errors.New("read timed out, waiting 5 minutes then restarting"), 10)
					if pubErr != nil {
						log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
					}
					time.Sleep(3 * time.Minute)
					cmd, err := exec.Command("bash", "-c", "sudo service phocus restart").Output()
					// it should die here
					log.Printf("cmd=================>%s\n", cmd)
					if err != nil {
						log.Fatal(err)
					}
					// if it reaches here at all that implies it didn't restart properly
					os.Exit(1)
				}
			}
			phocus_api.Queue = phocus_api.Queue[1:]
		} else {
			// min sleep between actual comms with inverter
			time.Sleep(5 * time.Second)
		}
		phocus_api.QueueMutex.Unlock()
		// min sleep between Queue checks
		time.Sleep(1 * time.Second)
	}
}
