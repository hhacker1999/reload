package main

import (
	"fmt"
	"reload/pkg/stack"
)

func main() {
	st := stack.NewStack()
	for i := 0; i < 20; i++ {
		st.Push(i)
	}
	st.Print()
	val, _ := st.Pop()
	fmt.Println(val)
	st.Print()
}
