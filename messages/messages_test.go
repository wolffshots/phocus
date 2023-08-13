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
	"go.bug.st/serial"
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
func TestMesages(t *testing.T) {
	cmd := StartCmd("socat", "PTY,link=./messages1,raw,echo=1,crnl", "PTY,link=./messages2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(10 * time.Millisecond)

	t.Run("TestInterpretWriteErrors", func(t *testing.T) {
		port1, err := phocus_serial.Setup("./messages1", 2400, 5)
		assert.NoError(t, err)
		assert.NoError(t, port1.Port.Close())
		port1.Port = nil

		err = Interpret(port1, Message{uuid.New(), "QPGS1", ""}, 0*time.Second)
		assert.EqualError(t, err, "port is nil on write")

		err = Interpret(port1, Message{uuid.New(), "QPGS2", ""}, 0*time.Second)
		assert.EqualError(t, err, "port is nil on write")

		err = Interpret(port1, Message{uuid.New(), "QID", ""}, 0*time.Second)
		assert.EqualError(t, err, "port is nil on write")

		err = Interpret(port1, Message{uuid.New(), "SOMETHING_ELSE", ""}, 0*time.Second)
		assert.EqualError(t, err, "port is nil on write")
	})

	t.Run("TestInterpretReadErrors", func(t *testing.T) {
		port2, err := phocus_serial.Setup("./messages1", 2400, 5)
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
	})

	t.Run("TestInterpret", func(t *testing.T) {
		port1, err := phocus_serial.Setup("./messages1", 2400, 5)
		assert.NoError(t, err)
		defer port1.Port.Close()

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\x2d\r", nil
		}
		err = Interpret(port1, Message{uuid.New(), "QPGS1", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "1 92932004102453 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\x9f\x50\r", nil
		}
		err = Interpret(port1, Message{uuid.New(), "QPGS2", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "92932004102453\xa7\x4a\r", nil
		}
		err = Interpret(port1, Message{uuid.New(), "QID", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")

		port1.Read = func(port serial.Port, timeout time.Duration) (string, error) {
			return "SOME_RESPONSE\xb2\xb2\r", nil
		}
		err = Interpret(port1, Message{uuid.New(), "SOME_MESSAGE", ""}, 0*time.Second)
		assert.EqualError(t, err, "client not defined in send")
	})
}
