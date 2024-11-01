// Package phocus_serial facilitates communication
// over RS232 with an inverter
package phocus_serial

import (
	"errors" // creating custom err messages
	"fmt"    // string formatting
	"log"    // logging
	"time"   // timeouts

	comms "github.com/wolffshots/phocus/v2/comms" // common types for comms
	crc "github.com/wolffshots/phocus/v2/crc"     // checksum generation
	serial "go.bug.st/serial"                     // rs232 serial
)

// type Writer func(port serial.Port, input string) (int, error)
// type Reader func(port serial.Port, timeout time.Duration) (string, error)

// port is the object representing the serial device/connection
// var port serial.Port

type Port struct {
	Port    *serial.Port
	Path    string
	Baud    int
	Retries int
}

// Open opens a connection to the inverter.
//
// Returns the port or an error if the port fails to open.
func (sp *Port) Open() (comms.Port, error) {
	fmt.Printf("Opening serial port: %s\n", sp.Path)
	var err error
	var port serial.Port
	for i := 0; i < sp.Retries; i++ {
		mode := &serial.Mode{
			BaudRate: sp.Baud,
		}
		port, err = serial.Open(sp.Path, mode)
		if err != nil {
			log.Printf("Failed to set up serial %d times with err: %v", i+1, err)
			time.Sleep(50 * time.Millisecond)
		} else {
			sp.Port = &port
			log.Printf("Succeeded to set up serial after %d times", i+1)
			break
		}
	}
	return sp, err
}

// Close just closes the port
//
// Returns the error if there is one
func (sp *Port) Close() error {
	fmt.Printf("Closing serial port: %s\n", sp.Path)
	if sp.Port == nil || sp == nil {
		return errors.New("serial port was nil on close call")
	}
	err := (*sp.Port).Close()
	if err == nil {
		sp.Port = nil
	}
	return err
}

// Read from the open serial port until reaching a carriage return, nil or nothing.
// Takes a duration as an input and times out the read after that long.
//
// Returns the read string and the error
func (sp *Port) Read(timeout time.Duration) (string, error) {
	log.Printf("Starting read\n")
	if sp.Port == nil || sp == nil {
		log.Printf("Port nil on read\n")
		return "", errors.New("port is nil on read")
	}
	buff := make([]byte, 140)
	err := (*sp.Port).SetReadTimeout(timeout)
	if err != nil {
		return "", errors.New("failed to set timeout")
	}
	var response string
	for {
		n, readErr := (*sp.Port).Read(buff)
		if readErr != nil {
			log.Printf("Err reading from port: %v", readErr)
			err = readErr
			break
		} else if n == 0 {
			log.Println("EOF")
			if response == "" {
				err = errors.New("read returned nothing")
			}
			break
		} else if string(buff[:n]) == "\r" {
			log.Printf("Encountered carriage return\n")
			response = fmt.Sprintf("%v%v", response, string(buff[:n]))
			break
		}
		log.Printf("Current response: %v\n", string(buff[:n]))
		response = fmt.Sprintf("%v%v", response, string(buff[:n]))
	}
	return response, err
}

// Write a string to the open serial port
// The input should just be the "payload" string as
// the CRC is calculated and added to that in Write
func (sp *Port) Write(input string) (int, error) {
	log.Printf("Starting write\n")
	if sp.Port == nil || sp == nil {
		return -1, errors.New("port is nil on write")
	}
	message := crc.Encode(input)
	n, err := (*sp.Port).Write([]byte(message))
	if err != nil {
		return -1, err
	}
	return n, err
}
