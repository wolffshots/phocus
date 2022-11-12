package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"   // for web server
	"github.com/google/uuid"     // for generating UUIDs for commands
	"log"                        // formatted logging
	"net/http"                   // for statuses primarily
	"time"                       // for sleeping
	"wolffshots/phocus/messages" // message structures
	"wolffshots/phocus/mqtt"     // comms with mqtt broker
	"wolffshots/phocus/sensors"
	"wolffshots/phocus/serial" // comms with inverter
)

// shape of a message for phocus to interpret and handle queuing of
type message struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Payload string    `json:"payload"`
}

// queue of messages seeded with QID to run at startup
var queue = []message{
	{ID: uuid.New(), Command: "QID", Payload: ""},
}

// loop and add QPGSi x n to the queue as long as it isn't too long
func QueueQPGSn() {
	for {
		if len(queue) < 20 {
			queue = append(
				queue,
				message{ID: uuid.New(), Command: "QPGSn", Payload: ""},
			)
		}
		time.Sleep(30 * time.Second)
	}
}

// enqueue new message manually (requires knowledge of commands and a generated uuid on the request)
func PostMessage(c *gin.Context) {
	var newMessage message
	// Call BindJSON to bind the received JSON to
	// newMessage - will throw an error if it can't cast ID to UUID
	if err := c.BindJSON(&newMessage); err != nil {
		return
	}
	queue = append(queue, newMessage)
	c.IndentedJSON(http.StatusCreated, newMessage)
}

// view current queue
func GetQueue(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, queue)
}

// get specific message
func GetMessage(c *gin.Context) {
	id := c.Param("id")

	if id == "next" && len(queue) > 0 {
		c.IndentedJSON(http.StatusOK, queue[0])
	} else {
		for _, a := range queue {
			if a.ID.String() == id {
				c.IndentedJSON(http.StatusOK, a)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}
}

// clear current queue
func DeleteQueue(c *gin.Context) {
	queue = []message{}
}

// function to interpret message and run relevant action (command or query)
func Interpret(input message) error {
	switch input.Command {
	case "QPGSn":
		log.Println("QPGS0") // if debug

		log.Println("QPGS1")
		serial.Write("QPGS1")
		response, err := serial.Read()
		if err != nil {
			log.Fatalf("failed to read from serial with :%v", err)
		}
		QPGSResponse := messages.NewQPGSnResponse(response)
		jsonQPGSResponse, err := json.Marshal(QPGSResponse)
		if err != nil {
			log.Fatalf("failed to parse response to json with :%v", err)
		}
		err = mqtt.Send("phocus/stats/qpgs1", 0, false, string(jsonQPGSResponse), 10)
		if err != nil {
			log.Fatalf("mqtt send of QPGS1 failed with: %v\ntype of thing sent was: %T", err, jsonQPGSResponse)
		}
		log.Printf("%v sent to mqtt with err: %v\njson object looked like: %s\n", QPGSResponse, err, jsonQPGSResponse)

		log.Println("QPGS2")
		serial.Write("QPGS2")
		response, err = serial.Read()
		if err != nil {
			log.Fatalf("failed to read from serial with :%v", err)
		}
		QPGSResponse = messages.NewQPGSnResponse(response)
		jsonQPGSResponse, err = json.Marshal(QPGSResponse)
		if err != nil {
			log.Fatalf("failed to parse response to json with :%v", err)
		}
		err = mqtt.Send("phocus/stats/qpgs2", 0, false, string(jsonQPGSResponse), 10)
		if err != nil {
			log.Fatalf("mqtt send of QPGS2 failed with: %v\ntype of thing sent was: %T", err, jsonQPGSResponse)
		}

		log.Printf("%v sent to mqtt with err: %v\njson object looked like: %s\n", QPGSResponse, err, jsonQPGSResponse)
		return err
	case "QID":
		log.Println("send QID")
	}
	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	log.Println("Starting up phocus")

	// serial
	serial.Setup()

	// router setup for async rest api for queueing
	router := gin.Default()
	router.GET("/queue", GetQueue)
	router.GET("/queue/:id", GetMessage)
	router.POST("/queue", PostMessage)
	router.DELETE("/queue", DeleteQueue)

	// spawns a go-routine which handles web requests
	go router.Run("localhost:8080")

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
		log.Printf("re-checking queue of length: %d", len(queue))
		// if there is an entry at [0] then run that command
		if len(queue) > 0 {
			Interpret(queue[0])
			queue = queue[1:len(queue)]
		}
		// min sleep between comms with inverter
		time.Sleep(5 * time.Second)
	}
}
