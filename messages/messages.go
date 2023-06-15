// Package phocus_messages contains the
// various queries and commands that
// that can be sent with phocus
package phocus_messages

import (
	"log"

	"github.com/google/uuid"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
	// comms with inverter
)

// Message is the shape of a message for phocus to interpret and handle queuing of
type Message struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Payload string    `json:"payload"`
}

// Response is the generic response from a Message
type Response struct {
	QPGSnResponse
}

var Port phocus_serial.Port

// Interpret converts the generic `phocus` message into a specific inverter message
// TODO add even more generalisation and separated implementation details here
func Interpret(port phocus_serial.Port, input Message) (*Response, error) {
	switch input.Command {
	case "QPGS1":
		response, err := HandleQPGS(port, 1)
		if err != nil {
			log.Printf("Failed to handle %s :%v\n", input.Command, err)
		}
		if response == nil {
			return nil, err
		} else {
			return &Response{QPGSnResponse: *response}, err
		}
	case "QPGS2":
		response, err := HandleQPGS(port, 2)
		if err != nil {
			log.Printf("Failed to handle %s :%v\n", input.Command, err)
		}
		if response == nil {
			return nil, err
		} else {
			return &Response{QPGSnResponse: *response}, err
		}
	case "QID":
		log.Println("TODO send QID")
	default:
		log.Println("Unexpected message on queue")
	}
	return nil, nil
}
