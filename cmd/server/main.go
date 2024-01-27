package main

import (
	"fmt"
	"os"
	"reload/internal/args"
	"reload/internal/reload"
	"runtime"
)

func main() {
	// Check for runtime
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		fmt.Println("Application only works in unix systems")
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error reading cwd ", err)
		os.Exit(1)
	}

	flags := args.NewAppFlags().
		WithFlag("-p", "")

	args, err := args.NewAppArgs(flags, args.HARD)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

  mainPath, _:=  args.GetFlagValue("-p")

	exitChan := make(chan struct{})
	r := reload.NewReload(cwd, mainPath,exitChan)
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
