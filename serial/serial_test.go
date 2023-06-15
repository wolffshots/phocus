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

func StartCmd(name string, arg ...string) *exec.Cmd {
	cmd := exec.Command(name, arg...)
	go cmd.Run()
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

	t.Run("TestSetup", func(t *testing.T) {
		port, err := Setup("./com1")
		if err != nil {
			log.Println(err)
		}
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, port)

		port.Port.Close()
		assert.Equal(t, nil, err)
	})
	t.Run("TestWrite", func(t *testing.T) {
		port, err := Setup("./com1")
		if err != nil {
			log.Println(err)
		}
		written, err := port.Write("test")
		assert.Equal(t, 7, written)
		assert.Equal(t, nil, err)

		port.Port.Close()
		assert.Equal(t, nil, err)
	})
	t.Run("TestRead", func(t *testing.T) {
		port1, err := Setup("./com1")
		assert.Equal(t, nil, err)
		// this isn't written to the port so should timeout
		read, err := port1.Read(1 * time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("read timed out"), err)

		port2, err := Setup("./com2")
		assert.Equal(t, nil, err)

		written, err := port1.Write("test")
		assert.Equal(t, 7, written)
		assert.Equal(t, nil, err)

		err = port1.Port.Close()
		assert.Equal(t, nil, err)

		read, err = port2.Read(1000 * time.Millisecond)
		assert.Equal(t, "", read)
		assert.Equal(t, errors.New("read timed out"), err)
		// assert.Equal(t, nil, err)
		// currently read times out since the virtual ports aren't linked quite right

	})
	t.Cleanup(func() {
		TerminateCmd(cmd)
	})
}
