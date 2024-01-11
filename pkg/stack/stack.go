package stack

import (
	"errors"
	"fmt"
)

type Stack struct {
	arr []int
}

func NewStack() Stack {
	return Stack{}
}

func (s *Stack) Push(value int) {
	s.arr = append(s.arr, value)
}

func (s *Stack) Pop() (int, error) {
	if len(s.arr) == 0 {
		return 0, errors.New("Stack is empty")
	}

	val := s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return val, nil
}

func (s *Stack) Print() {
	fmt.Println(s.arr)
}
