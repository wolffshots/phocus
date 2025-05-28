package phocus_ip

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPort_Open(t *testing.T) {
	tests := []struct {
		name    string
		port    Port
		wantErr bool
	}{
		{
			name: "Fail to open",
			port: Port{
				Host:    "localhost",
				Port:    8080,
				Retries: 3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.port.Open()
			if (err != nil) != tt.wantErr {
				t.Errorf("Port.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("Expected non-nil comms.Port, got nil")
			}
		})
	}
}

func TestPort_Close(t *testing.T) {
	port := Port{}
	if err := port.Close(); err != nil {
		t.Errorf("Port.Close() error = %v, want nil", err)
	}
}

func TestPort_Read(t *testing.T) {
	port := Port{}
	timeout := time.Second

	result, err := port.Read(timeout)
	assert.Equal(t, errors.New("ip port is not open"), err)
	if result != "" {
		t.Errorf("Port.Read() got = %v, want empty string", result)
	}
}

func TestPort_Write(t *testing.T) {
	port := Port{}
	input := "test_payload"

	n, err := port.Write(input)
	if err == nil || err.Error() != "ip port is nil on write" {
		t.Errorf("Port.Write() error = %v, want 'ip port is nil on write'", err)
	}
	if n != 0 {
		t.Errorf("Port.Write() bytes written = %v, want 0", n)
	}
}
