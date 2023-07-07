//go:build linux || darwin

package phocus_serial

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	cmd := StartCmd("socat", "PTY,link=./com1,raw,echo=1,crnl", "PTY,link=./com2,raw,echo=1,crnl")
	defer TerminateCmd(cmd)
	time.Sleep(200 * time.Millisecond)

	port1, err := Setup("./com1", 2400)
	defer port1.Port.Close()

	assert.NoError(t, err)
	assert.NotEqual(t, nil, port1)

	port2, err := Setup("./com2", 2400)
	defer port2.Port.Close()

	assert.NoError(t, err)
	assert.NotEqual(t, nil, port2)

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
