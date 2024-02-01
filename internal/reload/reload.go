package reload

import (
	"fmt"
	"reload/internal/runner"
	"sync"
	"time"
)

type Reload struct {
	changeChan           chan string
	blackList            []string
	root                 string
	files                []string
	reloadFiles          map[string]ReloadFileConfig
	exitChan             chan struct{}
	debounceStartTime    *time.Time
	currentlyInDebounce  bool
	currentRunner        *runner.Runner
	currentProcessOutput chan string
	mainPath             string
	indexingMutex        *sync.Mutex
	rebootingMutex       *sync.Mutex
	isIndexing           bool
	isRebooting          bool
}

func NewReload(root string, mainPath string, exitChan chan struct{}) Reload {
	fmt.Println(BANNER)
	return Reload{
		reloadFiles:    make(map[string]ReloadFileConfig),
		root:           root,
		changeChan:     make(chan string),
		blackList:      []string{".git", ".reload", "migrate.log"},
		exitChan:       exitChan,
		mainPath:       mainPath,
		indexingMutex:  &sync.Mutex{},
		rebootingMutex: &sync.Mutex{},
	}
}

func (r *Reload) ReadFiles() error {
	err := r.indexFiles()
	if err != nil {
		return err
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

	if r.isRebooting {
		fmt.Println("Rebooting is already in process")
		return
	}
	r.rebootingMutex.Lock()
	r.isRebooting = true
	defer r.rebootingMutex.Unlock()

	// Waiting for debounce to be over
	for time.Now().Sub(*r.debounceStartTime) <= time.Millisecond*800 {
	}

	r.debounceStartTime = nil
	r.currentlyInDebounce = false
	err := r.createBinary()
	if err != nil {
		r.isRebooting = false
		return
	}

	err = r.moveBinary()
	if err != nil {
		r.isRebooting = false
		return
	}

	if r.currentRunner != nil {
		r.currentRunner.Cleanup()
		r.currentRunner = nil
	}

	r.startBinary()

	go r.indexFiles()

	r.isRebooting = false
}
