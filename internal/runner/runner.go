package runner

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

const EXIT_SUCCESS = "EXIT SUCCESS"
const EXIT_ERROR = "EXIT ERROR"

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
	wg            *sync.WaitGroup
}

func NewRunner(inputChan chan string, cmd []string) Runner {
	oChan := make(chan string)
	writer := NewRunnerWriter(oChan)
	return Runner{
		inputChan:     inputChan,
		exitChan:      make(chan struct{}, 1),
		cmd:           cmd,
		outputChan:    oChan,
		currentWriter: &writer,
		wg:            &sync.WaitGroup{},
	}
}

func (r *Runner) Start() (chan string, error) {
	fmt.Println("Executing ", r.cmd)
	exc := exec.CommandContext(context.Background(), r.cmd[0], r.cmd[1:]...)
	exc.Stdout = r.currentWriter
	exc.Stderr = r.currentWriter
	r.exc = exc
	iChan, err := r.exc.StdinPipe()
	if err != nil {
		fmt.Println("Error running the command ", r.cmd, err)
		return r.outputChan, err
	}
	r.inputPipe = iChan
	err = r.exc.Start()
	if err != nil {
		fmt.Println("Error starting command ", r.cmd, err)
		return r.outputChan, err
	}
	go r.listenExit()
	go r.listenInput()
	return r.outputChan, nil
}

func (r *Runner) listenExit() {
	r.wg.Add(1)

	defer r.wg.Done()
	err := r.exc.Wait()
	if err == nil {
		r.outputChan <- EXIT_SUCCESS
	} else {
		r.outputChan <- EXIT_ERROR
	}
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
	fmt.Println("Inside cleanup")
	r.exitChan <- struct{}{}
	fmt.Println("After exit chan")
	r.exc.Cancel()
	r.wg.Wait()
	close(r.inputChan)
	close(r.outputChan)
	close(r.exitChan)
}
