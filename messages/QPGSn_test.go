package phocus_messages

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQPGSnResponse(t *testing.T) {
	// test grabbed input
	input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\x2d\r"
	want := &QPGSnResponse{5, true, "92932004102443", "Off-grid", "", "237.0", "50.01", "000.0", "00.00", "0483", "0387", "009", "51.1", "000", "069", "020.4", "000", "00942", "00792", "007", InverterStatus{"off", "off", "off", "Battery voltage normal", "connected", "on", "0"}, "Parallel output", "Solar first", "060", "080", "10", "00.0", "006", fmt.Sprintf("0x%x%x", 0xf2, 0x2d)}
	actual, err := NewQPGSnResponse(input, 5)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, actual)

	// test empty input
	input = ""
	want = (*QPGSnResponse)(nil)
	actual, err = NewQPGSnResponse(input, 0)
	assert.Equal(t, errors.New("can't create a response from an empty string"), err)
	assert.Equal(t, want, actual)
}
