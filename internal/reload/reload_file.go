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
				time.Sleep(time.Millisecond * 500)
			}
		}
	}()
}
