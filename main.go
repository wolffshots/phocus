package main

import (
    "fmt"                       // for printing
    "time"                      // for sleeping
    "net/http"                  // for statuses primarily

    "github.com/gin-gonic/gin"  // for web server
)

// shape of a message for phocus to interpret and handle queuing of
type message struct {
    ID      string  `json:"id"`
    Command string  `json:"command"`
    Payload string  `json:"payload"`
}

// queue of messages seeded with QID to run at startup
var messages = []message{
    {ID: "0", Command:"QID", Payload:""},
}

// loop and add QPGSi x n to the queue as long as it isn't too long


// endpoints to enqueue new messages
func postMessages(c *gin.Context) {
    var newMessage message

    // Call BindJSON to bind the received JSON to
    // newAlbum.
    if err := c.BindJSON(&newMessage); err != nil {
        return
    }

    // Add the new album to the slice.
    messages = append(messages, newMessage)
    c.IndentedJSON(http.StatusCreated, newMessage)
}

// endpoint to view current queue
func getMessages(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, messages)
}

// get specific message
func getMessageByID(c *gin.Context) {
    id := c.Param("id")
    
    if (id == "next" && len(messages)>0){
        c.IndentedJSON(http.StatusOK, messages[0])
    } else {
        for _, a := range messages {
            if a.ID == id {
                c.IndentedJSON(http.StatusOK, a)
                return
            }
        }
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
    }
}
// endpoint to clear current queue
func deleteMessages(c *gin.Context) {
    messages = []message{}
}

// function to interpret message and run relevant action (command or query)

// function to decode response

// function to check CRC values

func main() {
    router := gin.Default()
    router.GET("/messages", getMessages)
    router.GET("/messages/:id", getMessageByID)
    router.POST("/messages", postMessages)
    router.DELETE("/messages", deleteMessages)

    // spawns a go-routine which handles web requests
    go router.Run("localhost:8080")
    
    // sleep to make sure web server comes on before polling starts
    time.Sleep(5*time.Second)

    // loop to check queue and dequeue index 0, run it process result and wait 30 seconds
    for {
        fmt.Println("re-checking")
        fmt.Println(len(messages))
        fmt.Println("re-running")
        // if there is 
        if(len(messages)>0){
            fmt.Println(messages[0])
            messages = messages[1:len(messages)]
        }
        // sleep between comms with inverter
        time.Sleep(30*time.Second)
    }
}
