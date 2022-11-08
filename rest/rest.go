package rest

import (
	"github.com/gin-gonic/gin" // for web server
	"net/http"                 // for statuses primarily
)

func Setup() { // TODO err handling
	// router setup for async rest api for queueing
	router := gin.Default()
	router.GET("/messages", getMessages)
	router.GET("/messages/:id", getMessageByID)
	router.POST("/messages", postMessages)
	router.DELETE("/messages", deleteMessages)

	// spawns a go-routine which handles web requests
	router.Run("localhost:8080")
}
