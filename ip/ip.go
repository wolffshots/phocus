// Package phocus_ip facilitates communication
// over IP with an inverter
package phocus_ip

import (
	"errors" // creating custom err messages
	"time"   // timeouts

	comms "github.com/wolffshots/phocus/v2/comms" // common types for comms
	crc "github.com/wolffshots/phocus/v2/crc"     // checksum generation
)

type Port struct {
	// Socket
	Host    string
	Port    int
	Retries int
}

// Open opens a connection to the inverter.
//
// Returns the port or an error if the port fails to open.
func (ip *Port) Open() (comms.Port, error) {
	var err error
	for i := 0; i < ip.Retries; i++ {
		// TODO
	}
	return ip, err
}

// Close just closes the port
//
// Returns the error if there is one
func (ip *Port) Close() error {

	return nil // TODO
}

// Read from the open socket until reaching a carriage return, nil or nothing.
// Takes a duration as an input and times out the read after that long.
//
// Returns the read string and the error
func (ip *Port) Read(timeout time.Duration) (string, error) {

	return "", nil // TODO
}

// Write a string to the socket
// The input should just be the "payload" string as
// the CRC is calculated and added to that in Write
func (ip *Port) Write(input string) (int, error) {
	_ = crc.Encode(input) // TODO use message
	// if ip == nil {
	return 0, errors.New("ip port is nil on write")
	// }
	// TODO
	// return 0, nil
}
