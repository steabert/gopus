package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
)

func main() {
	flag.Parse()

	var entryPoint string
	switch flag.NArg() {
	case 0:
		u, err := user.Current()
		if err != nil {
			panic(fmt.Errorf("while getting current user: %v", err))
		}
		entryPoint = u.HomeDir
	case 1:
		entryPoint = flag.Arg(0)
	default:
		flag.Usage()
		os.Exit(1)
	}

	if entryPoint == "" {
	}

	fmt.Printf("Searching for .opus files in %s\n", entryPoint)
}
