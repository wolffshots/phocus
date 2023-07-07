// Package phocus_crc contains all the dependencies and functions
// for calculating and verifying checksums.
package phocus_crc

import (
	"fmt"
	"strings"

	"github.com/sigurn/crc16" // 16 bit checksum generation
)

// Checksum takes an input string and returns the crc for that
// string.
func Checksum(input string) uint16 {
	table := crc16.MakeTable(crc16.CRC16_XMODEM)
	result := crc16.Checksum([]byte(input), table)
	return result
}

// Encode takes an input string and returns that same string
// with it's crc attached.
func Encode(input string) string {
	checksum := Checksum(input)
	result := input + string([]byte{byte((checksum >> 8) & 0xff), byte(checksum & 0xff)}) + "\r"
	return result
}

// Verify requires an input string with the crc attached and
// will return whether the crc matches the content.
func Verify(input string) bool {
	input = strings.TrimRight(input, "\r")
	if len(input) < 2 {
		return false
	}
	crc := input[len(input)-2:]
	remainder := input[:len(input)-2]
	if remainder == "" { // we take the stance that empty inputs aren't valid
		return false
	}
	calculatedCrc := Checksum(remainder)
	calculatedCrcString := string([]byte{byte((calculatedCrc >> 8) & 0xff), byte(calculatedCrc & 0xff)})
	encodedRemainder := Encode(remainder)

	result := calculatedCrcString == crc && (input+"\r") == encodedRemainder
	if !result {
		fmt.Printf("%x\n", strings.TrimRight(encodedRemainder[len(input)-2:], "\r"))
		fmt.Printf("crc matches: %t\ninput matches: %t\ncrc: %x\nshould be: %x\n\n", calculatedCrcString == crc, (input+"\r") == encodedRemainder, crc, calculatedCrcString)
	}
	return result
}
