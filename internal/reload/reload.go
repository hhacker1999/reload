package reload

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reload/internal/runner"
	"time"
)

type ReloadFileConfig struct {
	file   *ReloadFile
	ctx    context.Context
	cancel context.CancelFunc
}

type Reload struct {
	changeChan           chan string
	blackList            []string
	root                 string
	files                []string
	reloadFiles          []ReloadFileConfig
	exitChan             chan struct{}
	debounceStartTime    *time.Time
	currentlyInDebounce  bool
	currentRunner        *runner.Runner
	currentProcessOutput chan string
	mainPath             string
}

func NewReload(root string, mainPath string, exitChan chan struct{}) Reload {
	fmt.Println(BANNER)
	return Reload{
		root:       root,
		changeChan: make(chan string),
		blackList:  []string{".git", ".reload"},
		exitChan:   exitChan,
		mainPath:   mainPath,
	}
}

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
	fmt.Println("Inside start binary")
	run := runner.NewRunner(make(chan string), []string{".reload/main"})
	fmt.Println("New runner created")
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
			fmt.Println("Channel is closed")
			return
		}
		fmt.Print(output)
	}
}

func (r *Reload) ReadFiles() error {
	r.initialise()

	err := r.readFileNonRec()
	if err != nil {
		fmt.Println("Error occured ", err)
		return err
	}
	for _, v := range r.files {
		temp, err := NewReloadFile(v, r.changeChan)
		if err != nil {
			fmt.Println("Error occured ", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		temp.StartListening(ctx)
		r.reloadFiles = append(r.reloadFiles, ReloadFileConfig{
			file:   &temp,
			ctx:    ctx,
			cancel: cancel,
		})
	}

	err = r.createBinary()
	if err != nil {
		fmt.Println("Error creating binary ", err)
		return err
	}

	err = r.moveBinary()
	if err != nil {
		return err
	}

	r.startBinary()

	r.listenChangeEvents()
	return nil
}

func (r *Reload) listenChangeEvents() {
	for {
		select {
		case changedFile := <-r.changeChan:
			fmt.Println("Changes detected in file ", changedFile)

			// Start debouncing
			cTime := time.Now()
			r.debounceStartTime = &cTime
			if !r.currentlyInDebounce {
				r.currentlyInDebounce = true
				go r.runnerRebooter()
			}
		case <-r.exitChan:
			fmt.Println("Exiting Reload")
		}
	}
}

func (r *Reload) runnerRebooter() {
	if r.debounceStartTime == nil {
		r.currentlyInDebounce = false
		return
	}

	// Waiting for debounce to be over
	for time.Now().Sub(*r.debounceStartTime) <= time.Second*2 {
	}

	r.debounceStartTime = nil
	r.currentlyInDebounce = false
	err := r.createBinary()
	if err != nil {
		return
	}

	err = r.moveBinary()
	if err != nil {
		return
	}

	if r.currentRunner != nil {
		r.currentRunner.Cleanup()
		r.currentRunner = nil
	}

	r.startBinary()
}

type ReloadFile struct {
	modTime   time.Time
	path      string
	eventChan chan string
}

func NewReloadFile(path string, eventChant chan string) (ReloadFile, error) {
	stat, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return ReloadFile{}, err
	}

	temp := ReloadFile{
		modTime:   stat.ModTime(),
		eventChan: eventChant,
		path:      path,
	}

	return temp, nil
}

func (r *ReloadFile) StartListening(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				cStat, err := os.Stat(r.path)
				if err != nil {
					fmt.Println("File moved ", r.path)
					return
				}
				cTime := cStat.ModTime()
				if r.modTime.Before(cTime) {
					r.modTime = cTime
					r.eventChan <- r.path
				}
				time.Sleep(time.Second * 2)
			}
		}
	}()
}
