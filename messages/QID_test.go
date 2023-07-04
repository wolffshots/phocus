package phocus_messages

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	phocus_crc "github.com/wolffshots/phocus/v2/crc"
)

func TestSendQID(t *testing.T) {
	assert.Equal(t, "QID\xd6\xea\r", phocus_crc.Encode("QID"))
}

func TestInterpretQID(t *testing.T) {
	// test grabbed input
	input := "(92932004102443\x94\x7b\r"
	want := &QIDResponse{"92932004102443"}
	actual, err := InterpretQID(input)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, actual)

	assert.Equal(t, true, phocus_crc.Verify(input[1:]))
	assert.Equal(t, true, phocus_crc.Verify(input))
	assert.Equal(t, uint16(0x947b), phocus_crc.Checksum(input[1:len(input)-3]))

	// test empty input
	input = ""
	want = (*QIDResponse)(nil)
	actual, err = InterpretQID(input)
	assert.Equal(t, errors.New("can't create a response from an empty string"), err)
	assert.Equal(t, want, actual)

	// test short input
	input = "(9\x94\x7b\r"
	want = (*QIDResponse)(nil)
	actual, err = InterpretQID(input)
	assert.Equal(t, errors.New("response is malformed or shorter than expected"), err)
	assert.Equal(t, want, actual)
}
