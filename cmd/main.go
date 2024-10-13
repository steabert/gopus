package main

import (
	"errors"
	"fmt"
	"os"
)

func usage() {
	fmt.Println(`gopus is an .opus song library manager.

Usage:
    gopus add <path>

  where <path> is a directory containing .opus files to
  be searched and added to the database.

    gopus find [-t title] [-a album] [-c creator] [-p performer]

  where results are filtered based on the provided flags.`)
}

func main() {
	var err error

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	cmdArgs := os.Args[2:]

	switch cmd {
	case "scan":
		err = scan(cmdArgs)
	case "list":
		err = list(cmdArgs)
	default:
		err = errors.New("no command given")
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "program encountered an error: %v\n", err)
		os.Exit(1)
	}
}
