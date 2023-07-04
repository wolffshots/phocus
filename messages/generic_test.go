package phocus_messages

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	phocus_crc "github.com/wolffshots/phocus/v2/crc"
)

func TestInterpretGeneric(t *testing.T) {
	// test grabbed input
	input := "(92932004102443\x2e\x2a\r"
	want := &GenericResponse{"92932004102443"}
	actual, err := InterpretGeneric(input)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, actual)

	assert.Equal(t, false, phocus_crc.Verify(input[1:]))
	assert.Equal(t, true, phocus_crc.Verify(input))
	assert.Equal(t, uint16(0x2e2a), phocus_crc.Checksum(input[:len(input)-3]))

	// test grabbed input
	input = "(ACK\x94\x7b\r"
	want = &GenericResponse{"ACK"}
	actual, err = InterpretGeneric(input)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, actual)

	// test grabbed input
	input = "(NAK\x94\x7b\r"
	want = &GenericResponse{"NAK"}
	actual, err = InterpretGeneric(input)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, actual)

	// test empty input
	input = ""
	want = (*GenericResponse)(nil)
	actual, err = InterpretGeneric(input)
	assert.Equal(t, errors.New("can't create a response from an empty string"), err)
	assert.Equal(t, want, actual)
}
