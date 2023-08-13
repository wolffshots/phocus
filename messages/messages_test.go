package phocus_messages

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	phocus_serial "github.com/wolffshots/phocus/v2/serial"
)

func Run(cmd *exec.Cmd) {
	err := cmd.Run()
	fmt.Printf("Couldn't run cmd: %v\n", err)
}

func StartCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	go Run(cmd)
	return cmd
}

func TerminateCmd(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		if runtime.GOOS == "windows" {
			cmd.Process.Signal(os.Kill)
		} else {
			cmd.Process.Signal(os.Interrupt)
		}
	} else {
		panic(fmt.Sprintf("command isn't running: %v", cmd))
	}
}

func TestInterpretWriteErrors(t *testing.T) {
	port1 := phocus_serial.Port{Port: nil, Path: "/nil"}

	err := Interpret(port1, Message{uuid.New(), "QPGS1", ""}, 0*time.Second)
	assert.EqualError(t, err, "port is nil on write")

	err = Interpret(port1, Message{uuid.New(), "QPGS2", ""}, 0*time.Second)
	assert.EqualError(t, err, "port is nil on write")

	err = Interpret(port1, Message{uuid.New(), "QID", ""}, 0*time.Second)
	assert.EqualError(t, err, "port is nil on write")

	err = Interpret(port1, Message{uuid.New(), "SOMETHING_ELSE", ""}, 0*time.Second)
	assert.EqualError(t, err, "port is nil on write")
}

func TestInterpretReadErrors(t *testing.T) {
	cmd := StartCmd("socat", "PTY,link=./com1,raw,echo=1,crnl", "PTY,link=./com2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)

	port2, err := phocus_serial.Setup("./com1", 2400, 5)
	assert.NoError(t, err)
	defer port2.Port.Close()

	err = Interpret(port2, Message{uuid.New(), "QPGS1", ""}, 0*time.Second)
	assert.EqualError(t, err, "read returned nothing")

	err = Interpret(port2, Message{uuid.New(), "QPGS2", ""}, 0*time.Second)
	assert.EqualError(t, err, "read returned nothing")

	err = Interpret(port2, Message{uuid.New(), "QID", ""}, 0*time.Second)
	assert.EqualError(t, err, "read returned nothing")

	err = Interpret(port2, Message{uuid.New(), "SOMETHING_ELSE", ""}, 0*time.Second)
	assert.EqualError(t, err, "read returned nothing")
}
