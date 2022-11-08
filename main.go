package main

import (
	"fmt"                    // for printing
	"github.com/google/uuid" // for generating UUIDs for commands
	"strings"                // for string parsing
	"time"                   // for sleeping
	"wolffshots/phocus/crc"
	"wolffshots/phocus/mqtt"
	"wolffshots/phocus/rest"
	"wolffshots/phocus/sensors"
	"wolffshots/phocus/serial"
)

// shape of a message for phocus to interpret and handle queuing of
type message struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Payload string    `json:"payload"`
}

// queue of messages seeded with QID to run at startup
var messages = []message{
	{ID: uuid.New(), Command: "QID", Payload: ""},
}

// loop and add QPGSi x n to the queue as long as it isn't too long
func queueQPGSn() {
	for {
		if len(messages) < 20 {
			messages = append(
				messages,
				message{ID: uuid.New(), Command: "QPGSn", Payload: ""},
			)
		}
		time.Sleep(30 * time.Second)
	}
}

// enqueue new message manually (requires knowledge of commands and a generated uuid on the request)
func postMessages(c *gin.Context) {
	var newMessage message
	// Call BindJSON to bind the received JSON to
	// newMessage - will throw an error if it can't cast ID to UUID
	if err := c.BindJSON(&newMessage); err != nil {
		return
	}
	messages = append(messages, newMessage)
	c.IndentedJSON(http.StatusCreated, newMessage)
}

// view current queue
func getMessages(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, messages)
}

// get specific message
func getMessageByID(c *gin.Context) {
	id := c.Param("id")

	if id == "next" && len(messages) > 0 {
		c.IndentedJSON(http.StatusOK, messages[0])
	} else {
		for _, a := range messages {
			if a.ID.String() == id {
				c.IndentedJSON(http.StatusOK, a)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}
}

// clear current queue
func deleteMessages(c *gin.Context) {
	messages = []message{}
}

// function to interpret message and run relevant action (command or query)
// TODO

// function to decode response
// TODO

func main() {
	serial.Setup()

	go rest.Setup()

	mqtt.Setup()

	sensors.Register()

	// sleep to make sure web server comes on before polling starts
	time.Sleep(5 * time.Second)

	// spawn go-routine to repeatedly enqueue QPGSn commands
	go queueQPGSn()

	// loop to check queue and dequeue index 0, run it process result and wait 30 seconds
	for {
		fmt.Println("re-checking")
		fmt.Println(len(messages))
		fmt.Println("re-running")
		// if there is an entry at [0] then run that command
		if len(messages) > 0 {
			fmt.Println(messages[0])
			messages = messages[1:len(messages)]
		}
		// min sleep between comms with inverter
		time.Sleep(10 * time.Second)
	}
}
