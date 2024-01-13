package main

import (
	"fmt"
	"os"
	"reload/internal/reload"
	"runtime"
)

func main() {
	// Check for runtime
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		fmt.Println("Application only works in unix like systems")
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error reading cwd ", err)
		os.Exit(1)
	}
	exitChan := make(chan struct{})
	r := reload.NewReload(cwd, exitChan)
	r.ReadFiles()
	//
	// run := runner.NewRunner(make(chan string), []string{"go", "run", "test.go"})
	// out, err := run.Start()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer run.Cleanup()
	// for {
	// 	val := <-out
	// 	fmt.Println(val)
	// }
}
