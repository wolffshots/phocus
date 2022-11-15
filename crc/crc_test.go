package crc

import (
	"github.com/stretchr/testify/assert"
	"testing" // testing framework
)

func TestChecksum(t *testing.T) {
	inputs := []string{"QPGS0", "QPGS1", "QPGS2", "QPGS3", "QPGS4", ""}
	wants := []uint16{0x3FDA, 0x2FFB, 0x1F98, 0x0FB9, 0x7F5E, 0x0000}
	for index, input := range inputs {
		result, err := Checksum(input)
		assert.Equal(t, nil, err)
		assert.Equal(t, wants[index], result)
	}
}

func TestEncode(t *testing.T) {
	inputs := []string{"QPGS0", "QPGS1", "QPGS2", "QPGS3", "QPGS4", ""}
	wants := []string{"QPGS0\x3F\xDA\r", "QPGS1\x2F\xFB\r", "QPGS2\x1F\x98\r", "QPGS3\x0F\xB9\r", "QPGS4\x7F\x5E\r", "\x00\x00\r"}
	for index, input := range inputs {
		result, err := Encode(input)
		assert.Equal(t, nil, err)
		assert.Equal(t, wants[index], result)
	}
}

func TestVerify(t *testing.T) {
	inputs := []string{"QPGS0\x3F\xDA\r", "QPGS1\x2F\xFB\r", "QPGS4\x3F\xDA\r", "QPGS2\x2F\xFB\r", "\x00\x00\r"}
	wants := []bool{true, true, false, false, false}
	for index, input := range inputs {
		result, err := Verify(input)
		assert.Equal(t, nil, err)
		assert.Equal(t, wants[index], result)
	}
}
