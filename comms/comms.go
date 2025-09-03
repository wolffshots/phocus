package phocus_comms

import (
	"time"
)

type Serial struct {
	Port    string
	Baud    int
	Retries int
}
type IP struct {
	Host    string
	Port    int
	Retries int
}
type ConnectionType string

type Connection struct {
	Type   ConnectionType
	Serial *Serial `json:"Serial,omitempty"`
	IP     *IP     `json:"IP,omitempty"`
}

const (
	ConnectionTypeSerial ConnectionType = "Serial"
	ConnectionTypeIP     ConnectionType = "IP"
)

type Port interface {
	Open() (Port, error)
	Close() error
	Read(timeout time.Duration) (string, error)
	Write(string) (int, error)
}
