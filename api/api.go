// Package phocus_api is the wrapper for the http and queueing systems
package phocus_api

import (
	"errors"
	"log" // formatted logging
	"math/rand"
	"net/http"
	"sync"
	"time" // for sleeping
	"unicode/utf8"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	messages "github.com/wolffshots/phocus/v2/messages"
	diagnostics "github.com/wolffshots/phocus/v2/diagnostics"
)

const MAX_QUEUE_LENGTH = 50

// Queue of messages seeded with QID to run at startup
var Queue = []messages.Message{
	{ID: uuid.New(), Command: "QID", Payload: ""},
}

// QueueMutex controls access to the Queue
var QueueMutex sync.Mutex

// ValueMutex controls access to the values
var ValueMutex sync.Mutex

var LastQPGSResponse *messages.QPGSnResponse

var upgrader = websocket.Upgrader{
	// ReadBufferSize:  4096,
	// WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return (origin == "http://127.0.0.1:8080" ||
			origin == "http://localhost:8080" ||
			true)
	},
}

// AddQPGSnMessages is the meat of the QueueQPGSn functionality
// TODO refactor to take a number in for how many inverters
func AddQPGSnMessages(timeBetween time.Duration) error {
	QueueMutex.Lock()
	if len(Queue) < 2 {
		Queue = append(
			Queue,
			messages.Message{ID: uuid.New(), Command: "QPGS1", Payload: ""},
		)
		QueueMutex.Unlock()
		time.Sleep(timeBetween)
		QueueMutex.Lock()
		Queue = append(
			Queue,
			messages.Message{ID: uuid.New(), Command: "QPGS2", Payload: ""},
		)
		QueueMutex.Unlock()
		time.Sleep(timeBetween)
	} else {
		QueueMutex.Unlock()
		return errors.New("queue too long")
	}
	return nil
}

// QueueQPGSn is a simple loop to add QPGSn to the Queue as long as it isn't too long
func QueueQPGSn(delaySeconds int, randDelaySeconds int) {
	for {
		AddQPGSnMessages(time.Duration(delaySeconds+rand.Intn(randDelaySeconds)) * time.Second)
	}
}

// PostMessage enqueues a new message manually (requires knowledge of commands and a generated uuid on the request)
// as long as there is space in the queue for it
func PostMessage(c *gin.Context) {
	var newMessage messages.Message
	// Call BindJSON to bind the received JSON to
	// newMessage - will throw an error if it can't cast ID to UUID
	if err := c.BindJSON(&newMessage); err != nil || newMessage.Command == "" {
		log.Printf("Error binding to JSON: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Coudln't bind JSON to message"})
	} else {
		QueueMutex.Lock()
		// append new message to the Queue if there is space
		if len(Queue) < MAX_QUEUE_LENGTH {
			Queue = append(Queue, newMessage)
			QueueMutex.Unlock()
			c.IndentedJSON(http.StatusCreated, newMessage)
		} else {
			QueueMutex.Unlock()
			c.IndentedJSON(http.StatusInsufficientStorage, gin.H{"message": "Message Queue already full!"})
		}
	}
}

// GetQueue is called to view the current Queue as JSON
func GetQueue(c *gin.Context) {
	QueueMutex.Lock()
	tempQueue := Queue
	QueueMutex.Unlock()
	c.IndentedJSON(http.StatusOK, tempQueue)
}

func SetLast(newResponse *messages.QPGSnResponse) {
	ValueMutex.Lock()
	LastQPGSResponse = newResponse
	ValueMutex.Unlock()
}

// GetLast is called to view the current Last Response as JSON
func GetLast(c *gin.Context) {
	ValueMutex.Lock()
	c.JSON(http.StatusOK, LastQPGSResponse)
	ValueMutex.Unlock()
}

// GetLast is called to view the current Last Response as JSON on a websocket
func GetLastWS(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade err:", err)
		return
	}
	defer c.Close()
	for {
		ValueMutex.Lock()
		bytesResponse := []byte(messages.EncodeQPGSn(LastQPGSResponse))
		ValueMutex.Unlock()
		if utf8.Valid(bytesResponse) {
			_ = c.WriteMessage(websocket.TextMessage, bytesResponse)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

type LastStateOfCharge struct {
	BatteryStateOfCharge string
}

// GetLastStateOfCharge is called to view the current Last State of Charge as JSON
func GetLastStateOfCharge(c *gin.Context) {
	ValueMutex.Lock()
	if LastQPGSResponse != nil {
		c.JSON(http.StatusOK, LastStateOfCharge{BatteryStateOfCharge: LastQPGSResponse.BatteryStateOfCharge})
	} else {
		c.JSON(http.StatusOK, LastStateOfCharge{BatteryStateOfCharge: "null"})
	}
	ValueMutex.Unlock()
}

// GetHealth is a simple endpoint to return a 200
func GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "UP")
}

// GetDiagnostics returns current system diagnostics as JSON
func GetDiagnostics(c *gin.Context) {
	diagData := diagnostics.GetDiagnostics()
	c.JSON(http.StatusOK, diagData)
}

// GetMessage attempts to select a specified message from the Queue and returns it or fails
//
// Attempts to get and return the Message with the supplied `id` from the Queue otherwise it returns a 404.
//
// Handles `next` as a reserved word for the next Message in the Queue
func GetMessage(c *gin.Context) {
	id := c.Param("id")
	QueueMutex.Lock()
	if id == "next" && len(Queue) > 0 {
		message := Queue[0]
		QueueMutex.Unlock()
		c.IndentedJSON(http.StatusOK, message)
	} else {
		for _, message := range Queue {
			if message.ID.String() == id {
				QueueMutex.Unlock()
				c.IndentedJSON(http.StatusOK, message)
				return
			}
		}
		QueueMutex.Unlock()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "message not found"})
	}
}

// DeleteQueue clears the current Queue
func DeleteQueue(c *gin.Context) {
	QueueMutex.Lock()
	Queue = []messages.Message{}
	QueueMutex.Unlock()
	c.Status(http.StatusNoContent)
}

// DeleteMessage attempts to delete a specified message from the Queue
//
// If the Queue is empty it or if the ID is not found it returns a 404 otherwise it returns an empty 204
func DeleteMessage(c *gin.Context) {
	id := c.Param("id")
	QueueMutex.Lock()
	if len(Queue) > 0 {
		for index, a := range Queue {
			if a.ID.String() == id {
				Queue = append(Queue[:index], Queue[index+1:]...)
				QueueMutex.Unlock()
				c.Status(http.StatusNoContent)
				return
			}
		}
		QueueMutex.Unlock()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "message not found"})
	} else {
		QueueMutex.Unlock()
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "message not found"})
	}
}

// SetupRouter adds the endpoints on the router for Queue management
//
// returns router object
func SetupRouter(mode string, profiling bool) *gin.Engine {
	gin.SetMode(mode)
	// router setup for async rest api for Queueing
	router := gin.Default()
	if profiling {
		log.Println("Starting profiling")
		pprof.Register(router)
	} else {
		log.Println("Running without profiling")
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://192.168.88.*:5173", "http://192.168.88.*"},
		MaxAge:       12 * time.Hour,
	}))
	router.GET("/health", GetHealth)
	router.GET("/diagnostics", GetDiagnostics)
	router.GET("/queue", GetQueue)
	router.GET("/queue/:id", GetMessage)
	router.GET("/last", GetLast)
	router.GET("/last-ws", GetLastWS)
	router.GET("/last/soc", GetLastStateOfCharge)
	router.POST("/queue", PostMessage)
	router.DELETE("/queue", DeleteQueue)
	router.DELETE("/queue/:id", DeleteMessage)
	return router
}
