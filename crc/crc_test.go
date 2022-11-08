package crc

import (
	"testing"
)

func TestChecksumQPGS0(t *testing.T) {
	input := "QPGS0"
	want := uint16(0x3FDA)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`CalculateCRC("QPGS0") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestChecksumQPGS1(t *testing.T) {
	input := "QPGS1"
	want := uint16(0x2FFB)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`CalculateCRC("QPGS1") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestChecksumQPGS2(t *testing.T) {
	input := "QPGS2"
	want := uint16(0x1F98)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`CalculateCRC("QPGS2") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestChecksumQPGS3(t *testing.T) {
	input := "QPGS3"
	want := uint16(0x0FB9)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`CalculateCRC("QPGS4") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}

func TestChecksumQPGS4(t *testing.T) {
	input := "QPGS4"
	want := uint16(0x7F5E)
	result, err := Checksum(input)
	if !(want == result) || err != nil {
		t.Fatalf(`CalculateCRC("QPGS4") = 0x%x, %v, want match for 0x%x, nil`, result, err, want)
	}
}
