package serial

import (
	"errors"                // creating custom err messages
	"fmt"                   // formatting
	"go.bug.st/serial"      // rs232 serial
	"log"                   // logging
	"time"                  // timeouts
	"wolffshots/phocus/crc" // checksum generation
)

var port serial.Port
var err error

func Setup() error { // TODO add error handling
	// specify serial port
	mode := &serial.Mode{
		BaudRate: 2400,
	}
	port, err = serial.Open("/dev/ttyUSB0", mode) // TODO move to environment variable
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Write a string to the open serial port
// The input should be the "payload" string as
// the CRC is calculated and added to that in Write
func Write(input string) (int, error) {
	message, err := crc.Encode(input)
	if err != nil {
		log.Fatal(err)
	}
	n, err := port.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	return n, err
}

// Read from the open serial port until reaching a carriage return, nil or nothing
// Takes a duration as an input and times out the read after that long
// Returns the read string and the error
func Read(timeout time.Duration) (string, error) {
	log.Printf("Starting read\n")
	buff := make([]byte, 140)
	dataChannel := make(chan string, 1)
	var response = ""
	// doesn't need to be this big but the biggest response we expect is 135 chars so might as well be slightly bigger than that
	// even though it reads one at a time in the current setup
	go func() {
		for {
			n, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
			} else if n == 0 {
				log.Println("\nEOF")
				break
			} else if string(buff[:n]) == "\r" {
				response = fmt.Sprintf("%v%v", response, string(buff[:n]))
				break
			} else {
				// log.Printf("%v",string(buff[:n]))
				// default case, no need to do anything special
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
