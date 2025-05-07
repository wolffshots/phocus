// Package phocus_mqtt handles connecting to and
// communicating with an MQTT broker
package phocus_mqtt

import (
	"errors"
	"fmt"  // string formatting
	"log"  // logging to stdout
	"os"   // verbose logging
	"time" // current time and timeouts

	mqtt "github.com/eclipse/paho.mqtt.golang" // mqtt client
)

type Client mqtt.Client

var CreateClient = func(hostname string, port int, retries int, clientId string) (mqtt.Client, error) {
	var client mqtt.Client
	// start mqtt setup
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	var mqttBroker = hostname
	var mqttPort = port
	var err error
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBroker, mqttPort))
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(messagePublishedHandler)
	opts.OnConnect = connectionHandler
	opts.OnConnectionLost = connectionLostHandler
	opts.AutoReconnect = true
	opts.SetPingTimeout(5 * time.Second)
	for i := 0; i < retries; i++ {
		client = mqtt.NewClient(opts)
		token := client.Connect()
		token.WaitTimeout(30 * time.Second)
		err = token.Error()
		if err != nil {
			log.Printf("Failed to set up mqtt %d times with err: %v", i+1, err)
			time.Sleep(50 * time.Millisecond)
		} else {
			log.Printf("Succeeded to set up mqtt after %d times", i+1)
			break
		}
	}
	return client, err
}

// Setup sets the logging and opens a connection to the broker
func Setup(hostname string, port int, retries int, clientId string) (mqtt.Client, error) {
	client, err := CreateClient(hostname, port, retries, clientId)

	if err != nil {
		return nil, err // i explicitly make client nil
	}

	// time needs to be formatted as iso8601 and rfc3339 is the closest to that
	err = Send(client, "phocus/stats/start_time", 0, false, time.Now().Format(time.RFC3339), 10)
	if err != nil {
		log.Printf("Failed to send initial setup stats to mqtt with err: %v", err)
	}

	return client, err
}

// Send uses the mqtt client to publish some data to a topic with a timeout
func Send(client mqtt.Client, topic string, qos byte, retained bool, payload interface{}, timeout time.Duration) error {
	if client == nil {
		return errors.New("client not defined in send")
	} else if !client.IsConnected() {
		return errors.New("client not connected in send")
	}
	token := client.Publish(topic, qos, retained, payload)
	err := token.Error()
	if err != nil {
		token.WaitTimeout(timeout)
		err = token.Error()
	}
	return err
}

// Error publishes a caught error to the error stat
func Error(client mqtt.Client, qos byte, retained bool, payload error, timeout time.Duration) error {
	err := Send(client, "phocus/stats/error", qos, retained, fmt.Sprint(payload), timeout)
	return err
}

// MQTT handlers

// messagePublishedHandler is called on every message publish
var messagePublishedHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	if client == nil {
		log.Println("Client is nil in messagePublishedHandler")
	} else if msg == nil {
		log.Println("Message is nil in messagePublishedHandler")
	} else {
		log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}
}

// connectionHandler is called on connection to the broker
var connectionHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	if client == nil {
		log.Println("Client is nil in connectionHandler")
	} else {
		log.Println("Connected")
	}
}

// connectionLostHandler is called on disconnection from the broker
var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	if client == nil {
		log.Println("Client is nil in connectionLostHandler")
	} else if err == nil {
		log.Println("Error is nil in connectionLostHandler")
	} else {
		log.Printf("Connection lost: %v\n", err)
	}
}
