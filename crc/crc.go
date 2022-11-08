package crc

import (
	"github.com/sigurn/crc16"
)

func Checksum(input string) (uint16, error) {
	table := crc16.MakeTable(crc16.CRC16_XMODEM)
	result := crc16.Checksum([]byte(input), table)
	return result, nil
}
