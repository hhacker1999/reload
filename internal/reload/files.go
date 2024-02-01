package reload

import (
	"os"
	"reload/pkg/arrs"
	"reload/pkg/stack"
)

func readFileRec(cwd string) ([]string, error) {
	var result []string
	entry, _ := os.ReadDir(cwd)
	for _, v := range entry {
		if v.IsDir() {
			if v.Name()[0] == '.' {
				continue
			}
			foo, _ := readFileRec(cwd + "/" + v.Name())
			result = append(result, foo...)
		} else {
			result = append(result, cwd+"/"+v.Name())
		}
	}
	return result, nil
}

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
