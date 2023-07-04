package phocus_messages

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	phocus_crc "github.com/wolffshots/phocus/v2/crc"
)

func TestInterpretQPGSn(t *testing.T) {
	// test grabbed input
	input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\x06\x6e\r"
	want := &QPGSnResponse{5, true, "92932004102443", "Off-grid", "", "237.0", "50.01", "000.0", "00.00", "0483", "0387", "009", "51.1", "000", "069", "020.4", "000", "00942", "00792", "007", InverterStatus{"off", "off", "off", "Battery voltage normal", "connected", "on", "0"}, "Parallel output", "Solar first", "060", "080", "10", "00.0", "006", fmt.Sprintf("0x%02x%02x", 0x06, 0x6e)}
	actual, err := InterpretQPGSn(input, 5)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, actual)

	assert.Equal(t, false, phocus_crc.Verify(input[1:]))
	assert.Equal(t, true, phocus_crc.Verify(input))
	assert.Equal(t, uint16(0xf22d), phocus_crc.Checksum(input[1:len(input)-3]))

	// test empty input
	input = ""
	want = (*QPGSnResponse)(nil)
	actual, err = InterpretQPGSn(input, 0)
	assert.Equal(t, errors.New("can't create a response from an empty string"), err)
	assert.Equal(t, want, actual)
}
