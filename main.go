package main

import (
	"encoding/json" // encoding to json for mqtt
	"errors"
	"fmt"                        // string formatting
	"github.com/gin-gonic/gin"   // for web server
	"github.com/google/uuid"     // for generating UUIDs for commands
	"log"                        // formatted logging
	"net/http"                   // for statuses primarily
	"sync"                       // mutexes for mutating queue
	"time"                       // for sleeping
	"wolffshots/phocus/crc"      // checksum calculations
	"wolffshots/phocus/messages" // message structures
	"wolffshots/phocus/mqtt"     // comms with mqtt broker
	"wolffshots/phocus/sensors"  // registering common sensors
	"wolffshots/phocus/serial"   // comms with inverter
)

// queue of messages seeded with QID to run at startup
var queue = []messages.Message{
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
				messages.Message{ID: uuid.New(), Command: "QPGSn", Payload: ""},
			)
		}
		queueMutex.Unlock()
		time.Sleep(60 * time.Second)
	}
}

// PostMessage enqueues a new message manually (requires knowledge of commands and a generated uuid on the request)
func PostMessage(c *gin.Context) {
	var newMessage messages.Message
	// Call BindJSON to bind the received JSON to
	// newMessage - will throw an error if it can't cast ID to UUID
	if err := c.BindJSON(&newMessage); err != nil || newMessage.Command == "" {
		log.Printf("Error binding to JSON: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Coudln't bind JSON to message"})
	} else {
		queueMutex.Lock()
		queue = append(queue, newMessage)
		queueMutex.Unlock()
		c.IndentedJSON(http.StatusCreated, newMessage)
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
	queue = []messages.Message{}
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

// HandleQPGS writes the query to the inverter and
// reads the response, deserialises it into a response
// object and sends it to MQTT
func HandleQPGS(inverterNum int) error {
	query := fmt.Sprintf("QPGS%d", inverterNum)
	log.Println(query)
	bytes, err := serial.Write(query)
	log.Printf("Sent %v bytes\n", bytes)
	if err != nil {
		log.Printf("Failed to write to serial with :%v\n", err)
		return err
	}
	response, err := serial.Read(2 * time.Second)
	if err != nil || response == "" {
		log.Printf("Failed to read from serial with :%v\n", err)
		return err
	}
	valid, err := crc.Verify(response)
	if err != nil {
		log.Fatalf("Verification of response from inverter produced an error :%v\n", err)
		return err
	}
	if valid {
		QPGSResponse, err := messages.NewQPGSnResponse(response)
		if err != nil || QPGSResponse == nil {
			log.Fatalf("Failed to create response with :%v", err)
		}
		jsonQPGSResponse, err := json.Marshal(QPGSResponse)
		if err != nil {
			log.Fatalf("Failed to parse response to json with :%v", err)
		}
		err = mqtt.Send(fmt.Sprintf("phocus/stats/qpgs%d", inverterNum), 0, false, string(jsonQPGSResponse), 10)
		if err != nil {
			log.Fatalf("MQTT send of %s failed with: %v\ntype of thing sent was: %T", query, err, jsonQPGSResponse)
		}
		log.Printf("Sent to MQTT:\n%s\n", jsonQPGSResponse)
	} else {
		log.Println("Invalid response from QPGSn")
		err = errors.New("invalid response from QPGSn")
	}
	return err
}

// Interpret converts the generic `phocus` message into a specific inverter message
// TODO add even more generalisation and separated implementation details here
func Interpret(input messages.Message) error {
	switch input.Command {
	case "QPGSn":
		err := HandleQPGS(1)
		if err != nil {
			log.Printf("Failed to handle QPGS1 :%v\n", err)
			return err
		}
		err = HandleQPGS(2)
		if err != nil {
			log.Printf("Failed to handle QPGS2 :%v\n", err)
			return err
		}
		return err
	case "QID":
		log.Println("TODO send QID")
	default:
		log.Println("Unexpected message on queue")
	}
	return nil
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

	// serial
	err := serial.Setup()
	if err != nil {
		log.Fatalf("Failed to set up serial with err: %v", err)
	}

	// spawns a go-routine which handles web requests
	go func() {
		err := setupRouter().Run("localhost:8080")
		if err != nil {
			log.Fatalf("Failed to run http routine with err: %v", err)
		}
	}()

	// mqtt
	mqtt.Setup()

	// sensors
	sensors.Register()

	// sleep to make sure web server comes on before polling starts
	time.Sleep(5 * time.Second)

	// spawn go-routine to repeatedly enqueue QPGSn commands
	go QueueQPGSn()

	// loop to check queue and dequeue index 0, run it process result and wait 30 seconds
	for {
		queueMutex.Lock()
		log.Printf("re-checking queue of length: %d", len(queue))
		// if there is an entry at [0] then run that command
		if len(queue) > 0 {
			err := Interpret(queue[0])
			if err != nil {
				log.Printf("Handle error (retry or not, maybe config or attempts): %v\n", err)
			}
			queue = queue[1:len(queue)] // TODO wrap in a mutex
		}
		queueMutex.Unlock()
		// min sleep between comms with inverter
		time.Sleep(5 * time.Second)
	}
}
