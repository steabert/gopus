package main

import (
	"fmt"
	"os"

	"github.com/steabert/gopus/rds"
	"github.com/steabert/gopus/worker"
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
	// Expecting at least 1 subcommand (add or find)
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	cmdArgs := os.Args[2:]

	switch cmd {
	case "add":
		if len(cmdArgs) != 1 {
			usage()
			os.Exit(1)
		}

		err := rds.Open("rw")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open database, %v", err)
			os.Exit(1)
		}

		dir := cmdArgs[0]
		err = worker.ScanDirectory(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan directory, %v", err)
			os.Exit(1)
		}
	case "find":
		// Open database in "ro" mode and search on artist/album/song
		// based on flag(s) given and return result as a list of paths.
	default:
		usage()
		os.Exit(1)
	}
}
