package phocus_mqtt

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	client, err := Setup(
		"badhost",
		1883,
		5,
		"test_client_name",
	)
	assert.Equal(t, errors.New("network Error : dial tcp: lookup badhost: no such host"), err)
	assert.Equal(t, nil, client)

	for i, message := range strings.Split(buf.String(), "\n") {
		if len(message) > 20 {
			assert.Equal(t, fmt.Sprintf("Failed to set up mqtt %d times with err: network Error : dial tcp: lookup badhost: no such host", i+1), message[20:])
		} else {
			assert.Equal(t, "", message)
		}
	}

	buf.Reset()

	client, err = Setup(
		"127.0.0.1",
		1883,
		5,
		"test_client_name",
	)
	assert.Equal(t, errors.New("network Error : dial tcp 127.0.0.1:1883: connect: connection refused"), err)
	assert.Equal(t, nil, client)

	for i, message := range strings.Split(buf.String(), "\n") {
		if len(message) > 20 {
			assert.Equal(t, fmt.Sprintf("Failed to set up mqtt %d times with err: network Error : dial tcp 127.0.0.1:1883: connect: connection refused", i+1), message[20:])
		} else {
			assert.Equal(t, "", message)
		}
	}

	buf.Reset()

	CreateClient = func(hostname string, port, retries int, clientId string) (mqtt.Client, error) {
		return nil, nil
	}

	client, err = Setup(
		"127.0.0.1",
		1883,
		5,
		"test_client_name",
	)
	assert.Equal(t, errors.New("client not defined in send"), err)
	assert.Equal(t, nil, client)

	message := buf.String()
	if len(message) > 20 {
		assert.Equal(t, "Failed to send initial setup stats to mqtt with err: client not defined in send", strings.Trim(message[20:], "\n"))
	} else {
		assert.Equal(t, "", message)
	}
}

func TestSend(t *testing.T) {
	var client mqtt.Client
	err := Send(client, "test/topic", 0, false, "payload", 10*time.Millisecond)
	assert.Equal(t, errors.New("client not defined in send"), err)

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)
	err = Send(client, "test/topic", 0, false, "payload", 10*time.Millisecond)
	assert.Equal(t, errors.New("client not connected in send"), err)
}

func TestError(t *testing.T) {
	var client mqtt.Client
	err := Error(client, 0, false, errors.New("example error"), 10*time.Millisecond)
	assert.Equal(t, errors.New("client not defined in send"), err)

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)
	err = Error(client, 0, false, errors.New("example error"), 10*time.Millisecond)
	assert.Equal(t, errors.New("client not connected in send"), err)
}

type mqttMessage struct {
	duplicate bool
	qos       byte
	retained  bool
	messageID uint16
	once      sync.Once
	ack       func()
}

func (m *mqttMessage) Duplicate() bool {
	return m.duplicate
}

func (m *mqttMessage) Qos() byte {
	return m.qos
}

func (m *mqttMessage) Retained() bool {
	return m.retained
}

func (m *mqttMessage) Topic() string {
	return "<some topic>"
}

func (m *mqttMessage) MessageID() uint16 {
	return m.messageID
}

func (m *mqttMessage) Payload() []byte {
	return bytes.NewBufferString("<some payload>").Bytes()
}

func (m *mqttMessage) Ack() {
	m.once.Do(m.ack)
}

func messageFromMQTTMessage(a *mqttMessage, ack func()) mqtt.Message {
	// Implement this function
	return mqtt.Message(a)
}

func TestMessagePublishedHandler(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	var client mqtt.Client
	messagePublishedHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Client is nil in messagePublishedHandler\n", buf.String()[20:])

	buf.Reset()

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)

	messagePublishedHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Message is nil in messagePublishedHandler\n", buf.String()[20:])

	buf.Reset()
	messagePublishedHandler(client, messageFromMQTTMessage(&mqttMessage{}, func() {}))
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Received message: <some payload> from topic: <some topic>\n", buf.String()[20:])
}

func TestConnectionHandler(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	var client mqtt.Client
	connectionHandler(client)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Client is nil in connectionHandler\n", buf.String()[20:])

	buf.Reset()

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)

	connectionHandler(client)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Connected\n", buf.String()[20:])
}

func TestConnectionLostHandler(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	var client mqtt.Client
	connectionLostHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Client is nil in connectionLostHandler\n", buf.String()[20:])

	buf.Reset()

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)

	connectionLostHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Error is nil in connectionLostHandler\n", buf.String()[20:])

	buf.Reset()
	connectionLostHandler(client, errors.New("some error"))
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Connection lost: some error\n", buf.String()[20:])
}
