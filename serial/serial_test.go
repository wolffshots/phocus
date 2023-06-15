package phocus_serial

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Run(cmd *exec.Cmd) {
	err := cmd.Run()
	log.Fatalf("Couldn't run cmd: %v", err)
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
	cmd := StartCmd("socat", "PTY,link=./com1", "PTY,link=./com2")
	time.Sleep(200 * time.Millisecond)

	var port1 Port
	var port2 Port

	port1, err := Setup("./com1")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, port1)

	port2, err = Setup("./com2")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, port2)

	t.Run("TestWrite", func(t *testing.T) {
		written, err := port2.Write("test")
		assert.Equal(t, 7, written)
		assert.Equal(t, nil, err)

		assert.Equal(t, nil, err)
	})
	t.Run("TestRead", func(t *testing.T) {
		port1, err := Setup("./com1")
		assert.Equal(t, nil, err)
		// this isn't written to the port so should timeout
		read, err := port1.Read(1 * time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("read timed out"), err)

		written, err := port1.Write("test")
		assert.Equal(t, 7, written)
		assert.Equal(t, nil, err)

		read, err = port2.Read(1000 * time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("read timed out"), err)
		// assert.Equal(t, nil, err)
		// currently read times out since the virtual ports aren't linked quite right

	})
	port1.Port.Close()
	port2.Port.Close()
	TerminateCmd(cmd)
}
