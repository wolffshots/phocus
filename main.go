// Package `main` contains the entrypoint for the `phocus` system
package main

import (
	"errors"  // creating custom errors
	"fmt"     // string formatting
	"log"     // formatted logging
	"os"      // exiting
	"os/exec" // auto restart
	"time"    // for sleeping

	api "github.com/wolffshots/phocus/v2/api"           // api setup
	messages "github.com/wolffshots/phocus/v2/messages" // message structures
	mqtt "github.com/wolffshots/phocus/v2/mqtt"         // comms with mqtt broker
	sensors "github.com/wolffshots/phocus/v2/sensors"   // registering common sensors
	serial "github.com/wolffshots/phocus/v2/serial"     // comms with inverter
)

const version = "v2.4.2"

// main is the entrypoint to the app
func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.Println("Starting up phocus")

	// mqtt
	// wrap in a for loop to retry Setup
	err := mqtt.Setup("192.168.88.14", "go_phocus_client'")
	if err != nil {
		log.Fatalf("Failed to set up mqtt with err: %v", err)
	}
	// reset error
	pubErr := mqtt.Send("phocus/stats/error", 0, false, "", 10)
	if pubErr != nil {
		log.Printf("Failed to clear previous error: %v\n", pubErr)
	}

	// serial
	// wrap in a for loop to retry Setup
	port, err := serial.Setup("/dev/ttyUSB0")
	if err != nil {
		pubErr := mqtt.Error(0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		port.Port.Close()
		log.Fatalf("Failed to set up serial with err: %v", err)
	}

	// spawns a go-routine which handles web requests
	go func() {
		err := api.SetupRouter().Run("0.0.0.0:8080")
		if err != nil {
			pubErr := mqtt.Error(0, false, err, 10)
			if pubErr != nil {
				log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
			}
			port.Port.Close()
			log.Fatalf("Failed to run http routine with err: %v", err)
		}
	}()

	// sensors
	// we only add them once we know the mqtt, serial and http aspects are up
	err = sensors.Register(version)
	if err != nil {
		pubErr := mqtt.Error(0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		port.Port.Close()
		log.Fatalf("Failed to set up sensors with err: %v", err)
	}

	// sleep to make sure web server comes on before polling starts
	time.Sleep(5 * time.Second)

	// spawn go-routine to repeatedly enQueue QPGSn commands
	go api.QueueQPGSn()

	// loop to check Queue and deQueue index 0, run it process result and wait 30 seconds
	for {
		api.QueueMutex.Lock()
		log.Print(".")
		// if there is an entry at [0] then run that command
		if len(api.Queue) > 0 {
			err := messages.Interpret(port, api.Queue[0])
			if err != nil {
				pubErr := mqtt.Error(0, false, err, 10)
				if pubErr != nil {
					log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
				}
				if fmt.Sprint(err) == "read timed out" { // immediately jailed when read timeout
					port.Port.Close()
					pubErr := mqtt.Error(0, false, errors.New("read timed out, waiting 5 minutes then restarting"), 10)
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
			api.Queue = api.Queue[1:]
		} else {
			// min sleep between actual comms with inverter
			time.Sleep(5 * time.Second)
		}
		api.QueueMutex.Unlock()
		// min sleep between Queue checks
		time.Sleep(1 * time.Second)
	}
}
