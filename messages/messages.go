// Package phocus_messages contains the
// various queries and commands that
// that can be sent with phocus
package phocus_messages

import (
	"time"

	"github.com/google/uuid"
	phocus_mqtt "github.com/wolffshots/phocus/v2/mqtt"
	phocus_serial "github.com/wolffshots/phocus/v2/serial" // comms with inverter
)

// Message is the shape of a message for phocus to interpret and handle queuing of
type Message struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Payload string    `json:"payload"`
}

// Interpret converts the generic `phocus` message into a specific inverter message
// TODO add even more generalisation and separated implementation details here
func Interpret(
	client phocus_mqtt.Client,
	port phocus_serial.Port,
	input Message,
	readTimeout time.Duration,
) (*QPGSnResponse, error) {
	switch input.Command {
	case "QPGS1":
		// send
		_, err := SendQPGSn(port, 1)
		if err != nil {
			return nil, err
		}
		// receive
		response, err := ReceiveQPGSn(port, readTimeout, 1)
		if err != nil {
			return nil, err
		} else {
			// interpret/handle
			QPGSnResponse, err := InterpretQPGSn(response, 1)
			if err != nil {
				return nil, err
			}
			// publish stuff here
			return PublishQPGSn(client, QPGSnResponse, 1)
		}
	case "QPGS2":
		// send
		_, err := SendQPGSn(port, 2)
		if err != nil {
			return nil, err
		}
		// receive
		response, err := ReceiveQPGSn(port, readTimeout, 2)
		if err != nil {
			return nil, err
		} else {
			// interpret/handle
			QPGSnResponse, err := InterpretQPGSn(response, 2)
			if err != nil {
				return nil, err
			}
			// publish stuff here
			return PublishQPGSn(client, QPGSnResponse, 2)
		}
	case "QID":
		// send
		_, err := SendQID(port, nil)
		if err != nil {
			return nil, err
		}
		// receive
		response, err := ReceiveQID(port, readTimeout)
		if err != nil {
			return nil, err
		} else {
			// interpret/handle
			QIDResponse, err := InterpretQID(response)
			if err != nil {
				return nil, err
			}
			// publish stuff here
			return nil, PublishQID(client, QIDResponse)
		}
	default:
		// generic handling (not suitable for complicated queries)
		// send
		_, err := SendGeneric(port, input.Command, input.Payload)
		if err != nil {
			return nil, err
		}
		// receive
		response, err := ReceiveGeneric(port, input.Command, readTimeout)
		if err != nil {
			return nil, err
		} else {
			// interpret/handle
			GenericResponse, err := InterpretGeneric(response)
			if err != nil {
				return nil, err
			}
			// publish stuff here
			return nil, PublishGeneric(client, GenericResponse, input.Command)
		}
	}
}
