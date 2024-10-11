package main

import (
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

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
		dir := cmdArgs[0]
		scanDirectory(dir)
	case "find":
		// findCmd := flag.NewFlagSet("find", flag.ExitOnError)
		// err := findCmd.Parse(flag.Args())
		// if err != nil {
		// 	findCmd.Usage()
		// 	os.Exit(1)
		// }
		// find(db, filter)
	default:
		usage()
		os.Exit(1)
	}
}

func scanDirectory(dir string) {
	fmt.Printf("scanning %s for Opus files to add to the database...\n", dir)

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		err = worker.InsertSongFromPath(path)
		if err != nil {
			fmt.Printf("[ERROR] failed to add %s, %v\n", path, err)
		} else {
			fmt.Printf("[OK] added %s\n", path)
		}
		return err
	})
}
