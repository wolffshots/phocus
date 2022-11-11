package mqtt

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang" // mqtt client
	"log"
	"os"
	"time"
)

var client mqtt.Client

func Setup() {
	// start mqtt setup
	mqtt.ERROR = log.New(os.Stdout, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	// mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	var mqtt_broker = "192.168.88.124" // TODO these should be config vars
	var mqtt_port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqtt_broker, mqtt_port))
	opts.SetClientID("go_phocus_client")
	opts.SetDefaultPublishHandler(messagePublishedHandler)
	opts.OnConnect = connectionHandler
	opts.OnConnectionLost = connectionLostHandler
	opts.SetPingTimeout(5 * time.Second)
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	Send("homeassistant/sensor/phocus/start_time/config", 0, true, `{"unique_id":"phocus_start_time","name":"phocus - Start Time","state_topic":"phocus/stats/start_time","icon":"mdi:hammer-wrench","device":{"name":"phocus","identifiers":["phocus"],"model":"phocus","manufacturer":"phocus","sw_version":"1.1.0"},"force_update":false}`, 10)
	Send("phocus/stats/start_time", 0, false, time.Now().Format(time.RFC822), 10)
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
