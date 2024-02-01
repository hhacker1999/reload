package reload

import (
	"context"
	"fmt"
	"os"
	"reload/pkg/arrs"
	"reload/pkg/stack"
)

func (r *Reload) readFileNonRec() error {
	r.files = []string{}
	dirStack := stack.NewStack()
	pathStack := stack.NewPathStack()

	c, err := os.ReadDir(r.root)
	if err != nil {
		return err
	}

	for _, v := range c {
		if v.IsDir() {
			if arrs.Contains(v.Name(), r.blackList) {
				continue
			}
			dirStack.Push(v)
			pathStack.Push(r.root)
		} else {
			r.files = append(r.files, r.root+"/"+v.Name())
		}
	}

	if err != nil {
		return err
	}

	for !dirStack.IsEmpty() && !pathStack.IsEmpty() {
		currentDir, _ := dirStack.Pop()
		currentPath, _ := pathStack.Pop()
		cPath := currentPath + "/" + currentDir.Name()
		arr, _ := os.ReadDir(cPath)
		for _, v := range arr {
			if v.IsDir() {
				if v.Name()[0] == '.' {
					continue
				}
				dirStack.Push(v)
				pathStack.Push(cPath)
			} else {
				r.files = append(r.files, cPath+"/"+v.Name())
			}
		}
	}

	return nil
}

func (r *Reload) indexFiles() error {
	// NOTE: This setup is done because read files can be called multiple times when we decide
	// to do indexing again
	if r.isIndexing {
		fmt.Println("Indexing is already in process")
		return nil
	}

	r.indexingMutex.Lock()
	r.isIndexing = true
	defer r.indexingMutex.Unlock()

	for _, v := range r.reloadFiles {
		v.cancel()
	}
	r.reloadFiles = make(map[string]ReloadFileConfig)

	r.initialise()

	err := r.readFileNonRec()
	if err != nil {
		fmt.Println("Error occured ", err)
		r.isIndexing = false
		return err
	}
	for _, v := range r.files {
		temp, err := NewReloadFile(v, r.changeChan)
		if err != nil {
			fmt.Println("Error occured ", err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		temp.StartListening(ctx)
		reloadFileConfig := ReloadFileConfig{
			file:   &temp,
			ctx:    ctx,
			cancel: cancel,
		}
		r.reloadFiles[v] = reloadFileConfig
	}

	r.isIndexing = false
	return nil
}
