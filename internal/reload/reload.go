package reload

import (
	"context"
	"fmt"
	"os"
	"time"
)

type ReloadFileConfig struct {
	file   *ReloadFile
	ctx    context.Context
	cancel context.CancelFunc
}

type Reload struct {
	changeChan          chan string
	blackList           []string
	root                string
	files               []string
	reloadFiles         []ReloadFileConfig
	exitChan            chan struct{}
	debounceStartTime   *time.Time
	currentlyInDebounce bool
}

func NewReload(root string, exitChan chan struct{}) Reload {
	return Reload{
		root:       root,
		changeChan: make(chan string),
		blackList:  []string{".git"},
		exitChan:   exitChan,
	}
}

func (r *Reload) ReadFiles() {
	err := r.readFileNonRec()
	if err != nil {
		fmt.Println("Error occured ", err)
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
	r.listenChangeEvents()
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

	fmt.Println("Doing some things here bruh")
	r.debounceStartTime = nil
	r.currentlyInDebounce = false
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
