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
	cmd := StartCmd("socat", "-d", "-d", "PTY,link=./serial1,raw,echo=0,crnl", "PTY,link=./serial2,raw,echo=0,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(51 * time.Millisecond)

	t.Run("TestSetup", func(t *testing.T) {
		var buf bytes.Buffer

		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		// sanity check with bug serial
		_, err := serial.Open("./bad_port", &serial.Mode{BaudRate: 2400})
		assert.Equal(t, syscall.Errno(0x2), err)

		badSerialPort := Port{
			Path:    "./bad_port",
			Baud:    2400,
			Retries: 5,
		}
		_, err = badSerialPort.Open()
		assert.Equal(t, "./bad_port", badSerialPort.Path)
		assert.NoFileExists(t, "./bad_port")
		for i, message := range strings.Split(buf.String(), "\n") {
			if len(message) > 20 {
				assert.Equal(t, fmt.Sprintf("Failed to set up serial %d times with err: no such file or directory", i+1), message[20:])
			} else {
				assert.Equal(t, "", message)
			}
		}
		assert.Equal(t, syscall.Errno(0x2), err)

		buf.Reset()

		time.Sleep(51 * time.Millisecond)

		// setup virtual ports
		serialPort1 := Port{
			Path:    "./serial1",
			Baud:    2400,
			Retries: 5,
		}
		commonPort1, err := serialPort1.Open()
		assert.NoError(t, err)
		defer commonPort1.Close()
		assert.Equal(t, "./serial1", serialPort1.Path)

		for i, message := range strings.Split(buf.String(), "\n") {
			if len(message) > 20 {
				assert.Equal(t, fmt.Sprintf("Succeeded to set up serial after %d times", i+1), message[20:])
			} else {
				assert.Equal(t, "", message)
			}
		}

		buf.Reset()

		serialPort2 := Port{
			Path:    "./serial2",
			Baud:    2400,
			Retries: 5,
		}
		commonPort2, err := serialPort2.Open()
		assert.NoError(t, err)
		defer commonPort2.Close()
		assert.Equal(t, "./serial2", serialPort2.Path)

		for i, message := range strings.Split(buf.String(), "\n") {
			if len(message) > 20 {
				assert.Equal(t, fmt.Sprintf("Succeeded to set up serial after %d times", i+1), message[20:])
			} else {
				assert.Equal(t, "", message)
			}
		}
	})

	t.Run("TestWrite", func(t *testing.T) {
		serialPort1 := Port{
			Path:    "./serial1",
			Baud:    2400,
			Retries: 5,
		}
		commonPort1, err := serialPort1.Open()
		assert.NoError(t, err)
		assert.Equal(t, "./serial1", serialPort1.Path)

		written, err := commonPort1.Write("test")
		assert.Equal(t, 7, written)
		assert.NoError(t, err)

		commonPort1.Close()
		written, err = commonPort1.Write("test")
		assert.Equal(t, -1, written)
		assert.Equal(t, errors.New("port is nil on write"), err)
	})

	t.Run("TestRead", func(t *testing.T) {
		serialPort1 := Port{
			Path:    "./serial1",
			Baud:    2400,
			Retries: 5,
		}
		commonPort1, err := serialPort1.Open()
		assert.NoError(t, err)
		assert.Equal(t, "./serial1", serialPort1.Path)

		read, err := commonPort1.Read(1 * time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("read returned nothing"), err)

		err = commonPort1.Close()
		assert.NoError(t, err)
		read, err = commonPort1.Read(1 * time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("port is nil on read"), err)
	})

	t.Run("TestReadWrite", func(t *testing.T) {
		serialPort1 := Port{
			Path:    "./serial1",
			Baud:    2400,
			Retries: 5,
		}
		commonPort1, err := serialPort1.Open()
		assert.NoError(t, err)
		defer commonPort1.Close()
		assert.Equal(t, "./serial1", serialPort1.Path)

		serialPort2 := Port{
			Path:    "./serial2",
			Baud:    2400,
			Retries: 5,
		}
		commonPort2, err := serialPort2.Open()
		assert.NoError(t, err)
		defer commonPort2.Close()
		assert.Equal(t, "./serial2", serialPort2.Path)

		// clear past read
		_, _ = commonPort2.Read(50 * time.Millisecond)

		written, err := commonPort1.Write("test")
		assert.Equal(t, 7, written)
		assert.NoError(t, err)

		time.Sleep(51 * time.Millisecond)

		read, err := commonPort2.Read(50 * time.Millisecond)
		assert.Equal(t, "test\x9b\x06\r", read)
		assert.NoError(t, err)
	})
}
