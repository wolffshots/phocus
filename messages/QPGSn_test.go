package messages

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewQPGSnResponse(t *testing.T) {
	// test grabbed input
	input := "1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006"
	want := &QPGSnResponse{true, "92932004102443", "Off-grid", "", "237.0", "50.01", "000.0", "00.00", "0483", "0387", "009", "51.1", "000", "069", "020.4", "000", "00942", "00792", "007", InverterStatus{"off", "off", "Battery voltage normal", "connected", "off", "1"}, "Parallel output", "Solar first", "060", "080", "10", "00.0", "", "006"}

	actual, err := NewQPGSnResponse(input)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, actual)

	// test empty input
	input = ""
	want = (*QPGSnResponse)(nil)
	actual, err = NewQPGSnResponse(input)
	assert.Equal(t, errors.New("can't create a response from an empty string"), err)
	assert.Equal(t, want, actual)
}
