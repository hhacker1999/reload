package capture

import (
	"fmt"
	"os"
	"os/exec"
)

type Capture struct {
	outChan  chan []byte
	exitChan chan struct{}
}

func NewCapture(oChan chan []byte) Capture {
	return Capture{
		outChan: oChan,
	}
}

func (c *Capture) StartCapturing() error {
	exec := exec.Command("stty", "-f", "/dev/tty", "raw")
	err := exec.Run()
	if err != nil {
		fmt.Println("Error starting capture ", err)
		return err
	}
	for {
		select {
		case <-c.exitChan:
			return nil
		default:
			os.Stdin.Read(<-c.outChan)
		}
	}
}

func (c *Capture) Cleanup() {
	defer exec.Command("stty", "-f", "/dev/tty", "-raw").Run()
	c.exitChan <- struct{}{}
}
