package main

import (
	"strings"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

func TestParseConfigFromFile(t *testing.T) {
	configuration, err := ParseConfigFromFile("config.json.example")
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", configuration.MQTT.Host)
	assert.Equal(t, 1883, configuration.MQTT.Port)
	assert.Equal(t, "go_phocus_client", configuration.MQTT.Client.Name)
	assert.Equal(t, 2400, configuration.Connection.Serial.Baud)
	assert.Equal(t, "/dev/ttyUSB0", configuration.Connection.Serial.Port)
	assert.Equal(t, 5, configuration.Connection.Serial.Retries)
	assert.Equal(t, 5, configuration.MQTT.Retries)
	assert.Equal(t, 2, configuration.Messages.Read.TimeoutSeconds)
	assert.Equal(t, 2*time.Second, time.Duration(configuration.Messages.Read.TimeoutSeconds)*time.Second)
	assert.Equal(t, 15, configuration.DelaySeconds)
	assert.Equal(t, 5, configuration.RandDelaySeconds)
	assert.Equal(t, 5, configuration.MinDelaySeconds)
}

func TestParseConfigSerial(t *testing.T) {
	configJSON := `{
        "Connection": {
            "Serial": {
                "Port": "/dev/ttyUSB0",
                "Baud": 2400,
                "Retries": 5
            }
        },
        "MQTT": {
            "Host": "192.168.1.1",
            "Port": 1883,
            "Client": {
                "Name": "go_phocus_client"
            },
            "Retries": 5
        },
        "Messages": {
            "Read": {
                "TimeoutSeconds": 2
            }
        },
        "DelaySeconds": 15,
        "RandDelaySeconds": 5,
        "MinDelaySeconds": 5,
        "Profiling": false
    }`

	reader := strings.NewReader(configJSON)
	configuration, err := ParseConfig(reader)
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", configuration.MQTT.Host)
	assert.Equal(t, 1883, configuration.MQTT.Port)
	assert.Equal(t, "go_phocus_client", configuration.MQTT.Client.Name)
	assert.Equal(t, 2400, configuration.Connection.Serial.Baud)
	assert.Equal(t, "/dev/ttyUSB0", configuration.Connection.Serial.Port)
	assert.Equal(t, 5, configuration.Connection.Serial.Retries)
	assert.Equal(t, 5, configuration.MQTT.Retries)
	assert.Equal(t, 2, configuration.Messages.Read.TimeoutSeconds)
	assert.Equal(t, 2*time.Second, time.Duration(configuration.Messages.Read.TimeoutSeconds)*time.Second)
	assert.Equal(t, 15, configuration.DelaySeconds)
	assert.Equal(t, 5, configuration.RandDelaySeconds)
	assert.Equal(t, 5, configuration.MinDelaySeconds)
}

func TestParseConfigIP(t *testing.T) {
	configJSON := `{
        "Connection": {
            "IP": {
                "Host": "example.com",
                "Port": 1080
            }
        },
        "MQTT": {
            "Host": "192.168.1.1",
            "Port": 1883,
            "Client": {
                "Name": "go_phocus_client"
            },
            "Retries": 5
        },
        "Messages": {
            "Read": {
                "TimeoutSeconds": 2
            }
        },
        "DelaySeconds": 15,
        "RandDelaySeconds": 5,
        "MinDelaySeconds": 5,
        "Profiling": false
    }`

	reader := strings.NewReader(configJSON)
	configuration, err := ParseConfig(reader)
	assert.NoError(t, err)
	assert.Equal(t, "example.com", configuration.Connection.IP.Host)
	assert.Equal(t, 1080, configuration.Connection.IP.Port)
	assert.Equal(t, "192.168.1.1", configuration.MQTT.Host)
	assert.Equal(t, 1883, configuration.MQTT.Port)
	assert.Equal(t, "go_phocus_client", configuration.MQTT.Client.Name)
	assert.Equal(t, 5, configuration.MQTT.Retries)
	assert.Equal(t, 2, configuration.Messages.Read.TimeoutSeconds)
	assert.Equal(t, 2*time.Second, time.Duration(configuration.Messages.Read.TimeoutSeconds)*time.Second)
	assert.Equal(t, 15, configuration.DelaySeconds)
	assert.Equal(t, 5, configuration.RandDelaySeconds)
	assert.Equal(t, 5, configuration.MinDelaySeconds)
}

func TestRouter(t *testing.T) {
	// Create a channel to communicate the server's start or error status
	startCh := make(chan error)
	var client mqtt.Client

	// Start the server in a goroutine
	go func() {
		startCh <- Router(client, true)
	}()

	time.Sleep(51 * time.Millisecond)

	select {
	case err := <-startCh:
		if err != nil {
			t.Errorf("Failed to start the server: %v", err)
		}
	default:
		// If no error received within the duration, assume the server has started successfully
	}

	// Perform additional test logic or assertions related to the running server here
	close(startCh)
}
