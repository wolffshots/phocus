package main

import (
    "fmt"                                       // for printing
    "time"                                      // for sleeping
    "net/http"                                  // for statuses primarily
    "os"
    "log"

    "github.com/gin-gonic/gin"                  // for web server
    "github.com/google/uuid"                    // for generating UUIDs for commands
    "github.com/eclipse/paho.mqtt.golang"       // mqtt client
)

// shape of a message for phocus to interpret and handle queuing of
type message struct {
    ID      uuid.UUID   `json:"id"`
    Command string      `json:"command"`
    Payload string      `json:"payload"`
}

// queue of messages seeded with QID to run at startup
var messages = []message{
    {ID: uuid.New(), Command:"QID", Payload:""},
}

// loop and add QPGSi x n to the queue as long as it isn't too long
func queueQPGSn() {
    for {
        if(len(messages)<20){
            messages = append(
                messages, 
                message{ID:uuid.New(), Command: "QPGSn", Payload:""}, 
            )
        }
        time.Sleep(30*time.Second)
    }
}

// enqueue new message manually (requires knowledge of commands and a generated uuid on the request)
func postMessages(c *gin.Context) {
    var newMessage message
    
    // Call BindJSON to bind the received JSON to
    // newMessage.
    if err := c.BindJSON(&newMessage); err != nil {
        return
    }

    // Add the new album to the slice.
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
    
    if (id == "next" && len(messages)>0){
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

// function to check CRC values
// TODO

// mqtt handlers
var messagePublishedHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("Connect lost: %v", err)
}

func main() {
    router := gin.Default()
    router.GET("/messages", getMessages)
    router.GET("/messages/:id", getMessageByID)
    router.POST("/messages", postMessages)
    router.DELETE("/messages", deleteMessages)

    // spawns a go-routine which handles web requests
    go router.Run("localhost:8080")

    // start mqtt setup
  	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
    var broker = "192.168.88.49" // TODO these should be config vars
    var port = 1883
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID("go_phocus_client")
    opts.SetDefaultPublishHandler(messagePublishedHandler)
    opts.OnConnect = connectionHandler
    opts.OnConnectionLost = connectionLostHandler
    opts.SetPingTimeout(5 * time.Second)
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

	token := client.Publish("go_phocus_client/boot_time", 0, false, time.Now().Format(time.RFC822))
    token.Wait()

    // sleep to make sure web server comes on before polling starts
    time.Sleep(5*time.Second)

    // spawn go-routine to repeatedly enqueue QPGSn commands
    go queueQPGSn()

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
        time.Sleep(10*time.Second)
    }
}
