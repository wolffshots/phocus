package serial

import (
	"fmt"              // formatting
	"go.bug.st/serial" // rs232 serial
	"log"              // logging
	"strings"          // string operations
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

func Read() (string, error) {
	buff := make([]byte, 140)
	var response = ""
	// doesn't need to be this big but the biggest response we expect is 135 chars so might as well be slightly bigger than that
	// even though it reads one at a time in the current setup
	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			log.Println("\nEOF")
			break
		}
		response = fmt.Sprintf("%v%v", response, string(buff[:n]))
		if string(buff[:n]) == "\r" {
			log.Print("read a \\r - response was: ")
			// this is what needs to be parsed for values based on the type of query it was
			log.Printf("other units:  \t%v\n", strings.Split(response, " ")[0])
			log.Printf("serial number:\t%v\n", strings.Split(response, " ")[1])
			// TODO seperate out the deserialisation of the commands to a generic function call with the input type as a parameter
			// we can handle updating mqtt values from that parser
			// TODO capture and make sense of the CRC in the response
			// crc.CalculateCRC("some input string")
			// crc.Verify
			break
		}
	}

	return response, err
}
