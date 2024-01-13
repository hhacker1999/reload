package stack

import (
	"errors"
	"fmt"
	"os"
)

type Stack struct {
	arr []os.DirEntry
}

func NewStack() Stack {
	return Stack{}
}

func (s *Stack) Push(value os.DirEntry) {
	s.arr = append(s.arr, value)
}

func (s *Stack) Pop() (os.DirEntry, error) {
	if len(s.arr) == 0 {
		return nil, errors.New("Stack is empty")
	}

	val := s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return val, nil
}

func (s *Stack) Print() {
	fmt.Println(s.arr)
}

func (s *Stack) IsEmpty() bool {
	return len(s.arr) == 0
}

type PathStack struct {
	arr []string
}

func NewPathStack() PathStack {
	return PathStack{}
}

func (s *PathStack) Push(value string) {
	s.arr = append(s.arr, value)
}

func (s *PathStack) Pop() (string, error) {
	if len(s.arr) == 0 {
		return "", errors.New("PathStack is empty")
	}

	val := s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return val, nil
}

func (s *PathStack) Print() {
	fmt.Println(s.arr)
}

func (s *PathStack) IsEmpty() bool {
	return len(s.arr) == 0
}
