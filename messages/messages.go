// Package messages contains the
// various queries and commands that
// that can be sent with phocus
package messages

import (
	"github.com/google/uuid"
	"log"
)

// Message is the shape of a message for phocus to interpret and handle queuing of
type Message struct {
	ID      uuid.UUID `json:"id"`
	Command string    `json:"command"`
	Payload string    `json:"payload"`
}

// Interpret converts the generic `phocus` message into a specific inverter message
func Interpret(input Message /* write func, read func, publish func */) error {
	switch input.Command {
	case "xxxx":
		// write func
		// read func
		// private package scoped struct instantiation
		// publish func
	default:
		log.Println("Unexpected message on queue")
	}
	return nil
}

type Command interface {
	New()
	Print()
}
