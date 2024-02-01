package reload

import (
	"errors"
	"fmt"
	"reload/internal/runner"
)

func (r *Reload) createBinary() error {
	run := runner.NewRunner(
		make(chan string),
		[]string{"go", "build", "-o", "main", r.root + "/" + r.mainPath},
	)
	oChan, err := run.Start()
	if err != nil {
		return err
	}

	for {
		output := <-oChan
		fmt.Println(output)
		if output == runner.EXIT_ERROR {
			return errors.New("Cannot build binary")
		} else if output == runner.EXIT_SUCCESS {
			break
		}
	}
	return nil
}

func (r *Reload) moveBinary() error {
	run := runner.NewRunner(
		make(chan string),
		[]string{"mv", "./main", r.root + "/" + ".reload/main"},
	)
	oChan, err := run.Start()
	if err != nil {
		return err
	}

	for {
		output := <-oChan
		fmt.Println(output)
		if output == runner.EXIT_ERROR {
			return errors.New("Unable to move binary")
		} else if output == runner.EXIT_SUCCESS {
			break
		}
	}
	return nil
}

func (r *Reload) startBinary() {
	run := runner.NewRunner(make(chan string), []string{".reload/main"})
	r.currentRunner = &run
	oChan, err := r.currentRunner.Start()
	if err != nil {
		r.currentRunner.Cleanup()
		fmt.Println("Error starting binary ", err)
		return
	}
	r.currentProcessOutput = oChan
	go r.printBinaryOutput()
}

func (r *Reload) printBinaryOutput() {
	for {
		output, ok := <-r.currentProcessOutput
		if !ok {
			return
		}
		fmt.Print(output)
	}
}
