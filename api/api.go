// Package phocus_api is the wrapper for the http and queueing api
package phocus_api

import (
	"log" // formatted logging
	"math/rand"
	"net/http"
	"sync"
	"time" // for sleeping

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	phocus_messages "github.com/wolffshots/phocus/messages" // message structures
)

// Queue of messages seeded with QID to run at startup
var Queue = []phocus_messages.Message{
	{ID: uuid.New(), Command: "QID", Payload: ""},
}
var QueueMutex sync.Mutex

// QueueQPGSn is a simple loop to add QPGSn to the Queue as long as it isn't too long
func QueueQPGSn() {
	for {
		QueueMutex.Lock()
		if len(Queue) < 2 {
			Queue = append(
				Queue,
				phocus_messages.Message{ID: uuid.New(), Command: "QPGS1", Payload: ""},
			)
			QueueMutex.Unlock()
			time.Sleep(time.Duration(15+rand.Intn(5)) * time.Second)
			QueueMutex.Lock()
			Queue = append(
				Queue,
				phocus_messages.Message{ID: uuid.New(), Command: "QPGS2", Payload: ""},
			)
			QueueMutex.Unlock()
			time.Sleep(time.Duration(15+rand.Intn(5)) * time.Second)
		} else {
			QueueMutex.Unlock()
		}

	}
}

// PostMessage enQueues a new message manually (requires knowledge of commands and a generated uuid on the request)
func PostMessage(c *gin.Context) {
	var newMessage phocus_messages.Message
	// Call BindJSON to bind the received JSON to
	// newMessage - will throw an error if it can't cast ID to UUID
	if err := c.BindJSON(&newMessage); err != nil || newMessage.Command == "" {
		log.Printf("Error binding to JSON: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Coudln't bind JSON to message"})
	} else {
		QueueMutex.Lock()
		// append new message to the Queue if there is space
		if len(Queue) < 50 {
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

// GetHealth is a simple endpoint to return a 200
func GetHealth(c *gin.Context) {
	c.String(http.StatusOK, "UP")
}

// GetMessage attempts to select a specified message from the Queue and returns it or fails
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
	Queue = []phocus_messages.Message{}
	QueueMutex.Unlock()
	c.Status(http.StatusNoContent)
}

// DeleteMessage attempts to delete a specified message from the Queue
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

// SetupRouter add the endpoints on the router for Queue management
func SetupRouter() *gin.Engine {
	// router setup for async rest api for Queueing
	router := gin.Default()
	router.GET("/health", GetHealth)
	router.GET("/queue", GetQueue)
	router.GET("/queue/:id", GetMessage)
	router.POST("/queue", PostMessage)
	router.DELETE("/queue", DeleteQueue)
	router.DELETE("/queue/:id", DeleteMessage)
	return router
}
