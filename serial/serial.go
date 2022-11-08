package serial

import (
	"fmt"
	"go.bug.st/serial" // rs232 serial
	"log"
	"strings"
)

func Setup() { // TODO add error handling
	// specify serial port
	mode := &serial.Mode{
		BaudRate: 2400,
	}
	port, err := serial.Open("/dev/ttyUSB0", mode) // TODO move to environment variable
	if err != nil {
		log.Fatal(err)
	}
	n, err := port.Write([]byte("QPGS0\x3F\xDA\r"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Sent %v bytes\n", n)
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
			fmt.Println("\nEOF")
			break
		}
		response = fmt.Sprintf("%v%v", response, string(buff[:n]))
		if string(buff[:n]) == "\r" {
			fmt.Print("read a \\r - response was: ")
			// this is what needs to be parsed for values based on the type of query it was
			fmt.Printf("other units:  \t%v\n", strings.Split(response, " ")[0])
			fmt.Printf("serial number:\t%v\n", strings.Split(response, " ")[1])
			// TODO seperate out the deserialisation of the commands to a generic function call with the input type as a parameter
			// we can handle updating mqtt values from that parser
			// TODO capture and make sense of the CRC in the response
			// crc.CalculateCRC("some input string")
			break
		}
	}
}
