package main

import (
	"github.com/gin-gonic/gin" // for web server
	"github.com/google/uuid"   // for generating UUIDs for commands
	"net/http"                 // for statuses primarily
	"time"                     // for sleeping
	"wolffshots/phocus/mqtt"
	"wolffshots/phocus/sensors"
	"wolffshots/phocus/serial"
    "log"
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
    log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
    log.Println("Starting up phocus")
	
    // serial
    serial.Setup()
	serial.Write("QPGS0")
    log.Println(serial.Read())

	// router setup for async rest api for queueing
	router := gin.Default()
	router.GET("/messages", getMessages)
	router.GET("/messages/:id", getMessageByID)
	router.POST("/messages", postMessages)
	router.DELETE("/messages", deleteMessages)

	// spawns a go-routine which handles web requests
	go router.Run("localhost:8080")

	// mqtt
    mqtt.Setup()

    // sensors
	sensors.Register()

	// sleep to make sure web server comes on before polling starts
	time.Sleep(5 * time.Second)

	// spawn go-routine to repeatedly enqueue QPGSn commands
	go queueQPGSn()

	// loop to check queue and dequeue index 0, run it process result and wait 30 seconds
	for {
		log.Println("re-checking")
		log.Println(len(messages))
		log.Println("re-running")
		// if there is an entry at [0] then run that command
		if len(messages) > 0 {
			log.Println(messages[0])
			messages = messages[1:len(messages)]
		}
		// min sleep between comms with inverter
		time.Sleep(10 * time.Second)
	}
}
