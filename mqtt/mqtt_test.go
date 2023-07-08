package phocus_mqtt

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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
		"bad host",
		1883,
		5,
		"test_client_name",
	)
	assert.Equal(t, errors.New("no servers defined to connect to"), err)
	assert.Equal(t, nil, client)

	for i, message := range strings.Split(buf.String(), "\n") {
		if len(message) > 20 {
			assert.Equal(t, fmt.Sprintf("Failed to set up mqtt %d times with err: no servers defined to connect to", i+1), message[20:])
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
}

func TestSend(t *testing.T) {
	client = nil
	err := Send("test/topic", 0, false, "payload", 10*time.Millisecond)
	assert.Equal(t, errors.New("client not defined in send"), err)

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)
	err = Send("test/topic", 0, false, "payload", 10*time.Millisecond)
	assert.Equal(t, errors.New("client not connected in send"), err)
}

func TestError(t *testing.T) {
	client = nil
	err := Error(0, false, errors.New("example error"), 10*time.Millisecond)
	assert.Equal(t, errors.New("client not defined in send"), err)

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)
	err = Error(0, false, errors.New("example error"), 10*time.Millisecond)
	assert.Equal(t, errors.New("client not connected in send"), err)
}

func TestMessagePublishedHandler(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	client = nil
	messagePublishedHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Client is nil in messagePublishedHandler\n", buf.String()[20:])

	buf.Reset()

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)

	messagePublishedHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Message is nil in messagePublishedHandler\n", buf.String()[20:])
}

func TestConnectionHandler(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	client = nil
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
	client = nil
	connectionLostHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Client is nil in connectionLostHandler\n", buf.String()[20:])

	buf.Reset()

	opts := mqtt.NewClientOptions()
	client = mqtt.NewClient(opts)

	connectionLostHandler(client, nil)
	assert.True(t, len(buf.String()) > 20)
	assert.Equal(t, "Error is nil in connectionLostHandler\n", buf.String()[20:])
}
