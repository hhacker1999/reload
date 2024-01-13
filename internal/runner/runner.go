package runner

import (
	"fmt"
	"io"
	"os/exec"
)

type RunnerWriter struct {
	oChan chan string
}

func NewRunnerWriter(oChan chan string) RunnerWriter {
	return RunnerWriter{oChan: oChan}
}

// Write(p []byte) (n int, err error)
func (r *RunnerWriter) Write(p []byte) (n int, err error) {
	r.oChan <- string(p)
	return len(p), nil
}

type Runner struct {
	cmd           []string
	exc           *exec.Cmd
	inputChan     chan string
	outputChan    chan string
	currentWriter *RunnerWriter
	inputPipe     io.WriteCloser
	exitChan      chan struct{}
}

func NewRunner(inputChan chan string, cmd []string) Runner {
	oChan := make(chan string)
	writer := NewRunnerWriter(oChan)
	return Runner{
		inputChan:     inputChan,
		cmd:           cmd,
		outputChan:    oChan,
		currentWriter: &writer,
	}
}

func (r *Runner) Start() (chan string, error) {
	fmt.Println("Executing ", r.cmd)
	exc := exec.Command(r.cmd[0], r.cmd[1:]...)
	exc.Stdout = r.currentWriter
	r.exc = exc
	iChan, err := exc.StdinPipe()
	if err != nil {
		fmt.Println("Error running the command ", r.cmd, err)
		return r.outputChan, err
	}
	r.inputPipe = iChan
	err = exc.Start()
	if err != nil {
		fmt.Println("Error starting command ", r.cmd, err)
		return r.outputChan, err
	}
	return r.outputChan, nil
}

func (r *Runner) listenInput() {
	for {
		select {
		case input := <-r.inputChan:
			r.inputPipe.Write([]byte(input))
		case <-r.exitChan:
			return
		}
	}

}

func (r *Runner) Cleanup() {
	r.exitChan <- struct{}{}
	r.exc.Cancel()
	close(r.inputChan)
	close(r.outputChan)
	close(r.exitChan)
}
