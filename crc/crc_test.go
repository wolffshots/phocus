package crc

import (
	"testing"
)

func TestChecksum(t *testing.T) {
	inputs := []string{"QPGS0", "QPGS1", "QPGS2", "QPGS3", "QPGS4"}
	wants := []uint16{0x3FDA, 0x2FFB, 0x1F98, 0x0FB9, 0x7F5E}
	for index, input := range inputs {
		result, err := Checksum(input)
		if !(wants[index] == result) || err != nil {
			t.Fatalf(`{Checksum("%s") = 0x%x, %v, wants 0x%x, nil}`, input, result, err, wants[index])
		}
	}
}

func TestEncode(t *testing.T) {
	inputs := []string{"QPGS0", "QPGS1", "QPGS2", "QPGS3", "QPGS4"}
	wants := []string{"QPGS0\x3F\xDA\r", "QPGS1\x2F\xFB\r", "QPGS2\x1F\x98\r", "QPGS3\x0F\xB9\r", "QPGS4\x7F\x5E\r"}
	for index, input := range inputs {
		result, err := Encode(input)
		if !(wants[index] == result) || err != nil {
			t.Fatalf(`{Encode("%s") = %s, %v, wants %s, nil}`, input, result, err, wants[index])
		}
	}
}

func TestVerify(t *testing.T) {
	inputs := []string{"QPGS0\x3F\xDA\r", "QPGS1\x2F\xFB\r", "QPGS4\x3F\xDA\r", "QPGS2\x2F\xFB\r"}
	wants := []bool{true, true, false, false}
	for index, input := range inputs {
		result, err := Verify(input)
		if !(wants[index] == result) || err != nil {
			t.Fatalf(`{Verify(%s) = %t , %v, wants %t, nil}`, input, result, err, wants[index])
		}
	}
}
