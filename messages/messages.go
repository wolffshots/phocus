// Package phocus_messages contains the
// various queries and commands that
// that can be sent with phocus
package phocus_messages

import (
	"time"

	"github.com/google/uuid"
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
func Interpret(port phocus_serial.Port, input Message, readTimeout time.Duration) error {
	switch input.Command {
	case "QPGS1":
		// send
		SendQPGSn(port, 1)
		// receive
		response, err := ReceiveQPGSn(port, readTimeout, 1)
		if err != nil {
			return err
		} else {
			// interpret/handle
			QPGSnResponse, err := InterpretQPGSn(response, 1)
			if err != nil {
				return err
			}
			// publish stuff here
			PublishQPGSn(QPGSnResponse, 1)
			return nil
		}
	case "QPGS2":
		// send
		SendQPGSn(port, 2)
		// receive
		response, err := ReceiveQPGSn(port, readTimeout, 2)
		if err != nil {
			return err
		} else {
			// interpret/handle
			QPGSnResponse, err := InterpretQPGSn(response, 2)
			if err != nil {
				return err
			}
			// publish stuff here
			PublishQPGSn(QPGSnResponse, 2)
			return nil
		}
	case "QID":
		// send
		SendQID(port, nil)
		// receive
		response, err := ReceiveQID(port, readTimeout)
		if err != nil {
			return err
		} else {
			// interpret/handle
			QIDResponse, err := InterpretQID(response)
			if err != nil {
				return err
			}
			// publish stuff here
			PublishQID(QIDResponse)
		}
	default:
		// generic handling (not suitable for complicated queries)
		// send
		SendGeneric(port, input.Command, input.Payload)
		// receive
		response, err := ReceiveGeneric(port, input.Command, readTimeout)
		if err != nil {
			return err
		} else {
			// interpret/handle
			GenericResponse, err := InterpretGeneric(response)
			if err != nil {
				return err
			}
			// publish stuff here
			PublishGeneric(GenericResponse, input.Command)
		}
	}
	return nil
}
