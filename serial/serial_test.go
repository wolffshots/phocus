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

	cmd := StartCmd("socat", "PTY,link=./com1,raw,echo=1,crnl", "PTY,link=./com2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(100 * time.Millisecond)

	port1, err := Setup("./com1", 2400, 5)
	defer port1.Port.Close()

	assert.NoError(t, err)
	assert.Equal(t, "./com1", port1.Path)

	for i, message := range strings.Split(buf.String(), "\n") {
		if len(message) > 20 {
			assert.Equal(t, fmt.Sprintf("Succeeded to set up serial after %d times", i+1), message[20:])
		} else {
			assert.Equal(t, "", message)
		}
	}

	buf.Reset()

	port2, err := Setup("./com2", 2400, 5)
	defer port2.Port.Close()

	assert.NoError(t, err)
	assert.Equal(t, "./com2", port2.Path)

	for i, message := range strings.Split(buf.String(), "\n") {
		if len(message) > 20 {
			assert.Equal(t, fmt.Sprintf("Succeeded to set up serial after %d times", i+1), message[20:])
		} else {
			assert.Equal(t, "", message)
		}
	}

	t.Run("TestWrite", func(t *testing.T) {
		written, err := port1.Write("test")
		assert.Equal(t, 7, written)
		assert.NoError(t, err)
	})

	t.Run("TestReadTimeout", func(t *testing.T) {
		read, err := port1.Read(1 * time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("read timed out"), err)
	})

	t.Run("TestRead", func(t *testing.T) {
		// Read operation 2 (asynchronous)
		// readChannel := make(chan string)
		// errChannel := make(chan error)

		// go func() {
		// 	read, _ := port2.Read(1000 * time.Millisecond)
		// 	readChannel <- read
		// }()

		// time.Sleep(100 * time.Millisecond)

		// written, err := port1.Write("test")
		// assert.Equal(t, 7, written)
		// assert.NoError(t, err)

		// select {
		// case err := <-errChannel:
		// 	assert.NoError(t, err)
		// case read := <-readChannel:
		// 	assert.Equal(t, "test", read)
		// }
	})
}
