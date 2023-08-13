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

type Writer func(port serial.Port, input string) (int, error)
type Reader func(port serial.Port, timeout time.Duration) (string, error)

// port is the object representing the serial device/connection
// var port serial.Port

type Port struct {
	Port  serial.Port
	Path  string
	Write Writer
	Read  Reader
}

// Setup opens a connection to the inverter.
//
// Returns the port or an error if the port fails to open.
func Setup(portPath string, baud int, retries int) (Port, error) {
	var port serial.Port
	var err error
	for i := 0; i < retries; i++ {
		mode := &serial.Mode{
			BaudRate: baud,
		}
		port, err = serial.Open(portPath, mode)
		if err != nil {
			log.Printf("Failed to set up serial %d times with err: %v", i+1, err)
			time.Sleep(50 * time.Millisecond)
		} else {
			log.Printf("Succeeded to set up serial after %d times", i+1)
			break
		}
	}
	return Port{
		Port:  port,
		Path:  portPath,
		Write: Write,
		Read:  Read,
	}, err
}

// Write a string to the open serial port
// The input should just be the "payload" string as
// the CRC is calculated and added to that in Write
var Write = func(port serial.Port, input string) (int, error) {
	message := crc.Encode(input)
	if port == nil {
		return 0, errors.New("port is nil on write")
	}
	n, err := port.Write([]byte(message))
	if err != nil {
		return -1, err
	}
	return n, err
}

// Read from the open serial port until reaching a carriage return, nil or nothing.
// Takes a duration as an input and times out the read after that long.
//
// Returns the read string and the error
var Read = func(port serial.Port, timeout time.Duration) (string, error) {
	log.Printf("Starting read\n")
	buff := make([]byte, 140)
	if port == nil {
		return "", errors.New("port is nil on read")
	}
	port.SetReadTimeout(timeout)
	var err error
	var response string
	for {
		n, readErr := port.Read(buff)
		if readErr != nil {
			log.Printf("Err reading from port: %v", readErr)
			err = readErr
			break
		} else if n == 0 {
			log.Println("\nEOF")
			err = errors.New("read returned nothing")
			break
		} else if string(buff[:n]) == "\r" {
			response = fmt.Sprintf("%v%v", response, string(buff[:n]))
			break
		}
		response = fmt.Sprintf("%v%v", response, string(buff[:n]))
	}
	return response, err
}
