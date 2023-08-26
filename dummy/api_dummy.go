package main

import (
	"log"
	"math/rand"
	"time"

	phocus_api "github.com/wolffshots/phocus/v2/api"
	phocus_messages "github.com/wolffshots/phocus/v2/messages"
)

func randomLast() {
	for {
		input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\x36\x29\r"
		actual, err := phocus_messages.InterpretQPGSn(input, rand.Intn(1000))
		if err == nil {
			phocus_api.SetLast(actual)
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	go randomLast()
	err := phocus_api.SetupRouter().Run("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("fatal err in router: %v", err)
	}
}
