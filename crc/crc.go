// Package containing all the dependencies and functions
// for calculating and verifying checksums.
package crc

import (
	"github.com/sigurn/crc16"
	"log"
)

// Takes an input string and returns the crc for that
// string.
func Checksum(input string) (uint16, error) {
	table := crc16.MakeTable(crc16.CRC16_XMODEM)
	result := crc16.Checksum([]byte(input), table)
	return result, nil
}

// Takes an input string and returns that same string
// with it's crc attached.
func Encode(input string) (string, error) {
	checksum, err := Checksum(input)
	log.Printf("crc for encode: %x\n", checksum)
	if err != nil {
		log.Fatal(err)
	}
	result := input + string([]byte{byte((checksum >> 8) & 0xff), byte(checksum & 0xff)}) + "\r"

	log.Printf("result to be written: %s\n", result)

	return result, err
}

// Verify requires an input string with the crc attached and
// will return whether the crc matches the content.
//
// Returns an error if the crc can't be detached or
// if another component returns an error.
func Verify(input string) (bool, error) {
	crc := input[len(input)-3 : len(input)-1]
	remainder := input[:len(input)-3]
	calculatedCrc, err := Checksum(remainder)
	calculatedCrcString := string([]byte{byte((calculatedCrc >> 8) & 0xff), byte(calculatedCrc & 0xff)})
	encodedRemainder, err := Encode(remainder)

	return calculatedCrcString == crc && input == encodedRemainder, err
}
