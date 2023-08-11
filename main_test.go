package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	configuration, err := ParseConfig("config.json.example")
	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", configuration.MQTT.Host)
	assert.Equal(t, 1883, configuration.MQTT.Port)
	assert.Equal(t, "go_phocus_client", configuration.MQTT.Client.Name)
	assert.Equal(t, 2400, configuration.Serial.Baud)
	assert.Equal(t, "/dev/ttyUSB0", configuration.Serial.Port)
	assert.Equal(t, 5, configuration.Serial.Retries)
	assert.Equal(t, 5, configuration.MQTT.Retries)
	assert.Equal(t, 2, configuration.Messages.Read.TimeoutSeconds)
	assert.Equal(t, 2*time.Second, time.Duration(configuration.Messages.Read.TimeoutSeconds)*time.Second)
	assert.Equal(t, 15, configuration.DelaySeconds)
	assert.Equal(t, 5, configuration.RandDelaySeconds)
}

func TestRouter(t *testing.T) {
	// Create a channel to communicate the server's start or error status
	startCh := make(chan error)

	// Start the server in a goroutine
	go func() {
		startCh <- Router()
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
