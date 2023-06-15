// Package phocus_serial facilitates communication
// over RS232 with an inverter
package phocus_serial

import (
	"errors" // creating custom err messages
	"fmt"    // formatting
	"log"    // logging
	"time"   // timeouts

	crc "github.com/wolffshots/phocus/v2/crc" // checksum generation
	"go.bug.st/serial"                        // rs232 serial
)

// port is the object representing the serial device/connection
// var port serial.Port

type Port struct {
	Port serial.Port
	Path string
}

// err is the error placeholder for serial connections
var err error

// Setup opens a connection to the inverter.
//
// Returns the port or an error if the port fails to open.
func Setup(portPath string) (Port, error) {
	// specify serial port
	mode := &serial.Mode{
		BaudRate: 2400,
	}
	port, err := serial.Open(portPath, mode) // TODO move to environment variable
	return Port{port, portPath}, err
}

// Write a string to the open serial port
// The input should be the "payload" string as
// the CRC is calculated and added to that in Write
func (p *Port) Write(input string) (int, error) {
	message := crc.Encode(input)
	n, err := p.Port.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	return n, err
}

// Read from the open serial port until reaching a carriage return, nil or nothing.
// Takes a duration as an input and times out the read after that long.
//
// Returns the read string and the error
func (p *Port) Read(timeout time.Duration) (string, error) {
	log.Printf("Starting read\n")
	buff := make([]byte, 140)
	dataChannel := make(chan string, 1)
	var response = ""
	// doesn't need to be this big but the biggest response we expect is 135 chars so might as well be slightly bigger than that
	// even though it reads one at a time in the current setup
	go func() {
		for {
			n, err := p.Port.Read(buff)
			if err != nil {
				log.Printf("Err reading from port: %v", err)
			} else if n == 0 {
				log.Println("\nEOF")
				break
			} else if string(buff[:n]) == "\r" {
				response = fmt.Sprintf("%v%v", response, string(buff[:n]))
				break
			}
			response = fmt.Sprintf("%v%v", response, string(buff[:n]))
		}
		dataChannel <- response
	}()
	select {
	case results := <-dataChannel:
		return results, err
	case <-time.After(timeout):
		return "", errors.New("read timed out")
	}
}
