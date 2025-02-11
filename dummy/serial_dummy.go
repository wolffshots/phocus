package main

import (
	"log"

	"go.bug.st/serial"
)

func serialDummy() {
	port, err := serial.Open("./testport2", &serial.Mode{BaudRate: 2400})
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	buffer := make([]byte, 128)

	for {
		n, err := port.Read(buffer)
		if err != nil {
			log.Println("Error reading from serial port:", err)
			break
		}
		if n == 0 {
			continue
		}

		// Log the received message
		received := string(buffer[:n])
		log.Println("Received:", received)

		// Determine the response based on the received message
		var response string
		switch received {
		case "ping\n":
			response = "pong\r"
		case "hello\n":
			response = "world\r"
		default:
			response = "unknown command\r"
		}

		// Send the response back through the serial port
		_, err = port.Write([]byte(response))
		if err != nil {
			log.Println("Error writing to serial port:", err)
		}
	}
}
