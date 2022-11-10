package crc

import (
	"testing"
)

// TODO condense tests with for loops

func TestChecksumQPGS0(t *testing.T) {
	input := "QPGS0"
	want := uint16(0x3FDA)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Checksum("QPGS0") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestEncodeQPGS0(t *testing.T) {
	input := "QPGS0"
	want := "QPGS0\x3F\xDA\r"
	result, err := Encode(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Encode("QPGS0") = %s, %v, want match for %s, nil`, result, err, want)
	}
}

func TestChecksumQPGS1(t *testing.T) {
	input := "QPGS1"
	want := uint16(0x2FFB)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Checksum("QPGS1") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestEncodeQPGS1(t *testing.T) {
	input := "QPGS1"
	want := "QPGS1\x2F\xFB\r"
	result, err := Encode(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Encode("QPGS1") = %s, %v, want match for %s, nil`, result, err, want)
	}
}

func TestChecksumQPGS2(t *testing.T) {
	input := "QPGS2"
	want := uint16(0x1F98)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Checksum("QPGS2") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestEncodeQPGS2(t *testing.T) {
	input := "QPGS2"
	want := "QPGS2\x1F\x98\r"
	result, err := Encode(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Encode("QPGS2") = %s, %v, want match for %s, nil`, result, err, want)
	}
}

func TestChecksumQPGS3(t *testing.T) {
	input := "QPGS3"
	want := uint16(0x0FB9)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Checksum("QPGS4") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestEncodeQPGS3(t *testing.T) {
	input := "QPGS3"
	want := "QPGS3\x0F\xB9\r"
	result, err := Encode(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Encode("QPGS3") = %s, %v, want match for %s, nil`, result, err, want)
	}
}

func TestChecksumQPGS4(t *testing.T) {
	input := "QPGS4"
	want := uint16(0x7F5E)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Checksum("QPGS4") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestEncodeQPGS4(t *testing.T) {
	input := "QPGS4"
	want := "QPGS4\x7F\x5E\r"
	result, err := Encode(input)
	if !(want == result) || err != nil {
		t.Fatalf(`Encode("QPGS4") = %s, %v, want match for %s, nil`, result, err, want)
	}
}
