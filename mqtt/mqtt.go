package mqtt

import (
	"fmt"                                 // string formatting
	"github.com/eclipse/paho.mqtt.golang" // mqtt client
	"log"                                 // logging to stdout
	"os"                                  // verbose logging
	"time"                                // current time and timeouts
)

var client mqtt.Client

func Setup() {
	// start mqtt setup
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	var mqttBroker = "192.168.88.124" // TODO these should be config vars
	var mqttPort = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBroker, mqttPort))
	opts.SetClientID("go_phocus_client")
	opts.SetDefaultPublishHandler(messagePublishedHandler)
	opts.OnConnect = connectionHandler
	opts.OnConnectionLost = connectionLostHandler
	opts.SetPingTimeout(5 * time.Second)
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

    // TODO extract to a update function that uses Send
	err := Send("phocus/stats/start_time", 0, false, time.Now().Format(time.RFC822), 10)
	if err != nil {
		log.Fatalf("Failed to send initial setup stats to MQTT with err: %v", err)
	}
}

func Send(topic string, qos byte, retained bool, payload interface{}, timeout time.Duration) error {
	token := client.Publish(topic, qos, retained, payload)
	err := token.Error()
	if err != nil {
		token.WaitTimeout(timeout)
		err = token.Error()
	}
	return err
}

// mqtt handlers
var messagePublishedHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v\n", err)
}
