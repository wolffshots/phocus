//go:build linux || darwin

package phocus_serial

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestSerial(t *testing.T) {
	cmd := StartCmd("socat", "PTY,link=./serial1,raw,echo=1,crnl", "PTY,link=./serial2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(10 * time.Millisecond)

	t.Run("TestSetup", func(t *testing.T) {
		var buf bytes.Buffer

		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()
		badPort, err := Setup("./bad_port", 2400, 5)
		assert.Equal(t, syscall.Errno(0x2), err)
		assert.Equal(t, "./bad_port", badPort.Path)
		for i, message := range strings.Split(buf.String(), "\n") {
			if len(message) > 20 {
				assert.Equal(t, fmt.Sprintf("Failed to set up serial %d times with err: no such file or directory", i+1), message[20:])
			} else {
				assert.Equal(t, "", message)
			}
		}

		buf.Reset()

		time.Sleep(51 * time.Millisecond)

		port1, err := Setup("./serial1", 2400, 5)
		assert.NoError(t, err)
		defer port1.Port.Close()
		assert.Equal(t, "./serial1", port1.Path)

		for i, message := range strings.Split(buf.String(), "\n") {
			if len(message) > 20 {
				assert.Equal(t, fmt.Sprintf("Succeeded to set up serial after %d times", i+1), message[20:])
			} else {
				assert.Equal(t, "", message)
			}
		}

		buf.Reset()

		port2, err := Setup("./serial2", 2400, 5)
		assert.NoError(t, err)
		defer port2.Port.Close()
		assert.Equal(t, "./serial2", port2.Path)

		for i, message := range strings.Split(buf.String(), "\n") {
			if len(message) > 20 {
				assert.Equal(t, fmt.Sprintf("Succeeded to set up serial after %d times", i+1), message[20:])
			} else {
				assert.Equal(t, "", message)
			}
		}
	})

	t.Run("TestWrite", func(t *testing.T) {
		port1, err := Setup("./serial1", 2400, 5)
		assert.NoError(t, err)
		written, err := port1.Write(port1.Port, "test")
		assert.Equal(t, 7, written)
		assert.NoError(t, err)

		port1.Port.Close()
		written, err = port1.Write(port1.Port, "test")
		assert.Equal(t, -1, written)
		assert.Equal(t, syscall.Errno(0x9), err)

		port1.Port = nil
		written, err = port1.Write(port1.Port, "test")
		assert.Equal(t, 0, written)
		assert.Equal(t, errors.New("port is nil on write"), err)
	})

	t.Run("TestRead", func(t *testing.T) {
		port1, err := Setup("./serial1", 2400, 5)
		assert.NoError(t, err)
		read, err := port1.Read(port1.Port, 1*time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("read returned nothing"), err)

		err = port1.Port.Close()
		assert.NoError(t, err)
		read, err = port1.Read(port1.Port, 1*time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, serial.PortClosed, err.(*serial.PortError).Code())

		port1.Port = nil
		read, err = port1.Read(port1.Port, 1*time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("port is nil on read"), err)
	})
}
