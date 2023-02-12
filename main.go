package main

import (
	"errors"                                // creating custom errors
	"fmt"                                   // string formatting
	"github.com/gin-gonic/gin"              // for web server
	"github.com/google/uuid"                // for generating UUIDs for commands
	"github.com/wolffshots/phocus_messages" // message structures
	"github.com/wolffshots/phocus_mqtt"     // comms with mqtt broker
	"github.com/wolffshots/phocus_sensors"  // registering common sensors
	"github.com/wolffshots/phocus_serial"   // comms with inverter
	"log"                                   // formatted logging
	"math/rand"                             // queue randomisation
	"net/http"                              // for statuses primarily
	"os"                                    // exiting
	"os/exec"                               // auto restart
	"sync"                                  // mutexes for mutating queue
	"time"                                  // for sleeping
)

// queue of messages seeded with QID to run at startup
var queue = []phocus_messages.Message{
	{ID: uuid.New(), Command: "QID", Payload: ""},
}
var queueMutex sync.Mutex

// QueueQPGSn is a simple loop to add QPGSn to the queue as long as it isn't too long
func QueueQPGSn() {
	for {
		queueMutex.Lock()
		if len(queue) < 20 {
			queue = append(
				queue,
				phocus_messages.Message{ID: uuid.New(), Command: "QPGS1", Payload: ""},
			)
		}
		queueMutex.Unlock()
		time.Sleep(time.Duration(15+rand.Intn(5)) * time.Second)
		queueMutex.Lock()
		if len(queue) < 20 {
			queue = append(
				queue,
				phocus_messages.Message{ID: uuid.New(), Command: "QPGS2", Payload: ""},
			)
		}
		time.Sleep(time.Duration(15+rand.Intn(5)) * time.Second)
		queueMutex.Unlock()
	}
}

// PostMessage enqueues a new message manually (requires knowledge of commands and a generated uuid on the request)
func PostMessage(c *gin.Context) {
	var newMessage phocus_messages.Message
	// Call BindJSON to bind the received JSON to
	// newMessage - will throw an error if it can't cast ID to UUID
	if err := c.BindJSON(&newMessage); err != nil || newMessage.Command == "" {
		log.Printf("Error binding to JSON: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Coudln't bind JSON to message"})
	} else {
		queueMutex.Lock()
		// append new message to the queue if there is space
		if len(queue) < 50 {
			queue = append(queue, newMessage)
			queueMutex.Unlock()
			c.IndentedJSON(http.StatusCreated, newMessage)
		} else {
			queueMutex.Unlock()
			c.IndentedJSON(http.StatusInsufficientStorage, gin.H{"message": "Message queue already full!"})
		}
	}
}

// GetQueue is called to view the current queue as JSON
func GetQueue(c *gin.Context) {
	queueMutex.Lock()
	tempQueue := queue
	queueMutex.Unlock()
	c.IndentedJSON(http.StatusOK, tempQueue)
}

func GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "UP")
}

// GetMessage attempts to select a specified message from the queue and returns it or fails
func GetMessage(c *gin.Context) {
	id := c.Param("id")
	queueMutex.Lock()
	if id == "next" && len(queue) > 0 {
		message := queue[0]
		queueMutex.Unlock()
		c.IndentedJSON(http.StatusOK, message)
	} else {
		for _, message := range queue {
			if message.ID.String() == id {
				queueMutex.Unlock()
				c.IndentedJSON(http.StatusOK, message)
				return
			}
		}
		queueMutex.Unlock()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "message not found"})
	}
}

// DeleteQueue clears the current queue
func DeleteQueue(c *gin.Context) {
	queueMutex.Lock()
	queue = []phocus_messages.Message{}
	queueMutex.Unlock()
	c.Status(http.StatusNoContent)
}

// DeleteMessage attempts to delete a specified message from the queue
func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	queueMutex.Lock()
	if len(queue) > 0 {
		for index, a := range queue {
			if a.ID.String() == id {
				queue = append(queue[:index], queue[index+1:]...)
				queueMutex.Unlock()
				c.Status(http.StatusNoContent)
				return
			}
		}
		queueMutex.Unlock()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "message not found"})
	} else {
		queueMutex.Unlock()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "message not found"})
	}
}

func setupRouter() *gin.Engine {
	// router setup for async rest api for queueing
	router := gin.Default()
	router.GET("/health", GetHealth)
	router.GET("/queue", GetQueue)
	router.GET("/queue/:id", GetMessage)
	router.POST("/queue", PostMessage)
	router.DELETE("/queue", DeleteQueue)
	router.DELETE("/queue/:id", DeleteMessage)
	return router
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.Println("Starting up phocus")

	// mqtt
	err := phocus_mqtt.Setup("192.168.88.124", "go_phocus_client'")
	if err != nil {
		log.Fatalf("Failed to set up mqtt with err: %v", err)
	}
	// reset error
	pubErr := phocus_mqtt.Send("phocus/stats/error", 0, false, fmt.Sprint(""), 10)
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
		err := setupRouter().Run("localhost:8080")
		if err != nil {
			pubErr := phocus_mqtt.Error(0, false, err, 10)
			if pubErr != nil {
				log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
			}
			log.Fatalf("Failed to run http routine with err: %v", err)
		}
	}()

	// sensors
	err = phocus_sensors.Register()
	if err != nil {
		pubErr := phocus_mqtt.Error(0, false, err, 10)
		if pubErr != nil {
			log.Printf("Failed to post previous error (%v) to mqtt: %v\n", err, pubErr)
		}
		log.Fatalf("Failed to set up sensors with err: %v", err)
	}

	// sleep to make sure web server comes on before polling starts
	time.Sleep(5 * time.Second)

	// spawn go-routine to repeatedly enqueue QPGSn commands
	go QueueQPGSn()

	// loop to check queue and dequeue index 0, run it process result and wait 30 seconds
	for {
		queueMutex.Lock()
		log.Print(".")
		// if there is an entry at [0] then run that command
		if len(queue) > 0 {
			err := phocus_messages.Interpret(queue[0])
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
			queue = queue[1:]
		}else{
            // min sleep between actual comms with inverter
            time.Sleep(5 * time.Second)
        }
		queueMutex.Unlock()
		// min sleep between queue checks
		time.Sleep(1 * time.Second)
	}
}
