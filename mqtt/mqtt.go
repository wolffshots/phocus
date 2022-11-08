package mqtt

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang" // mqtt client
	"log"
	"os"
	"time"
)

func Setup() {
	// start mqtt setup
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	var mqtt_broker = "192.168.88.49" // TODO these should be config vars
	var mqtt_port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqtt_broker, mqtt_port))
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
}

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
