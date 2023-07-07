package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	configuration, err := ParseConfig("config.json.example")
	assert.Equal(t, nil, err)
	assert.Equal(t, "192.168.1.1", configuration.MQTT.Host)
	assert.Equal(t, 1883, configuration.MQTT.Port)
	assert.Equal(t, "go_phocus_client", configuration.MQTT.Client.Name)
	assert.Equal(t, 2400, configuration.Serial.Baud)
	assert.Equal(t, "/dev/ttyUSB0", configuration.Serial.Port)
}
