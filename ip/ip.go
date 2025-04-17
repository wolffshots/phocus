// Package phocus_ip facilitates communication
// over IP with an inverter
package phocus_ip

import (
	"bytes"  // byte manipulation
	"errors" // creating custom err messages
	"fmt"    // string formatting
	"log"    // logging
	"net"    // networking
	"time"   // timeouts

	comms "github.com/wolffshots/phocus/v2/comms" // common types for comms
	crc "github.com/wolffshots/phocus/v2/crc"     // checksum generation
)

type Port struct {
	// Socket
	Host    string
	Port    int
	Retries int
	Conn    net.Conn // active connection
}

// Open opens a connection to the inverter.
//
// Returns the port or an error if the port fails to open.
func (ip *Port) Open() (comms.Port, error) {
	var err error
	address := fmt.Sprintf("%s:%d", ip.Host, ip.Port)
	log.Printf("Connecting to %s over IP stream\n", address)
	var conn net.Conn
	for i := 0; i < ip.Retries; i++ {
		log.Printf("Attempt %d/%d\n", i+1, ip.Retries)
		conn, err = net.DialTimeout("tcp", address, 5*time.Second)
		if err == nil {
			log.Printf("Connected to %s on attempt %d\n", address, i+1)
			ip.Conn = conn
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	return ip, err
}

// Close just closes the port
//
// Returns the error if there is one
func (ip *Port) Close() error {
	if ip.Conn != nil {
		return ip.Conn.Close()
	}
	return nil
}

// Read from the open socket until reaching a carriage return, nil or nothing.
// Takes a duration as an input and times out the read after that long.
//
// Returns the read string and the error
func (ip *Port) Read(timeout time.Duration) (string, error) {
	if ip.Conn == nil {
		return "", errors.New("ip port is not open")
	}
	ip.Conn.SetReadDeadline(time.Now().Add(timeout))
	var data []byte
	buf := make([]byte, 256)
	for {
		n, err := ip.Conn.Read(buf)
		if err != nil {
			return "", err
		}
		data = append(data, buf[:n]...)
		// Check if there's a carriage return in the data
		if i := bytes.IndexByte(data, '\r'); i != -1 {
			return string(data[:i+1]), nil
		}
	}
}

// Write a string to the socket
// The input should just be the "payload" string as
// the CRC is calculated and added to that in Write
func (ip *Port) Write(input string) (int, error) {
	encoded := crc.Encode(input)
	if ip.Conn == nil {
		return 0, errors.New("ip port is nil on write")
	}
	n, err := ip.Conn.Write([]byte(encoded))
	return n, err
}
