package phocus_mqtt

import (
	"fmt"  // string formatting
	"log"  // logging to stdout
	"os"   // verbose logging
	"time" // current time and timeouts

	mqtt "github.com/eclipse/paho.mqtt.golang" // mqtt client
)

var client mqtt.Client

// Setup sets the logging and opens a connection to the broker
func Setup(hostname string, clientId string) error {
	// start mqtt setup
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	var mqttBroker = hostname // TODO these should be config vars
	var mqttPort = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBroker, mqttPort))
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(messagePublishedHandler)
	opts.OnConnect = connectionHandler
	opts.OnConnectionLost = connectionLostHandler
	opts.AutoReconnect = true
	opts.SetPingTimeout(5 * time.Second)
	client = mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Printf("Failed to connect to mqtt with err: %v", err)
		return err
	}

	// time needs to be formatted as iso8601 and rfc3339 is the closest to that
	err = Send("phocus/stats/start_time", 0, false, time.Now().Format(time.RFC3339), 10)
	if err != nil {
		log.Printf("Failed to send initial setup stats to mqtt with err: %v", err)
	}

	return err
}

// Send uses the mqtt client to publish some data to a topic with a timeout
func Send(topic string, qos byte, retained bool, payload interface{}, timeout time.Duration) error {
	token := client.Publish(topic, qos, retained, payload)
	err := token.Error()
	if err != nil {
		token.WaitTimeout(timeout)
		err = token.Error()
	}
	return err
}

// Error publishes a caught error to the error stat
func Error(qos byte, retained bool, payload error, timeout time.Duration) error {
	err := Send("phocus/stats/error", qos, retained, fmt.Sprint(payload), timeout)
	return err
}

// MQTT handlers

// messagePublishedHandler is called on every message publish
var messagePublishedHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

// connectionHandler is called on connection to the broker
var connectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

// connectionLostHandler is called on disconnection from the broker
var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v\n", err)
}
