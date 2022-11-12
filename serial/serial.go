package serial

import (
	"fmt"              // formatting
	"go.bug.st/serial" // rs232 serial
	"log"              // logging
	"wolffshots/phocus/crc"
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

func Write(input string) (int, error) {
	message, err := crc.Encode(input)
	if err != nil {
		log.Fatal(err)
	}
	n, err := port.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Sent %v bytes\n", n)
	return n, err
}

// Reads from the open serial port until reaching a carriage return, nil or nothing
// TODO implement timeout
func Read() (string, error) {
    log.Printf("Starting read\n")
	buff := make([]byte, 140)
	var response = ""
	// doesn't need to be this big but the biggest response we expect is 135 chars so might as well be slightly bigger than that
	// even though it reads one at a time in the current setup
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
        } else if n == 0 {
			log.Println("\nEOF")
			break
        } else if string(buff[:n]) == "\r" {
		    response = fmt.Sprintf("%v%v", response, string(buff[:n]))
			break
        } else{
            // log.Printf("%v",string(buff[:n]))
            // default case, no need to do anything special
        }
		response = fmt.Sprintf("%v%v", response, string(buff[:n]))
	}

	return response, err
}
