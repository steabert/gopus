package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/steabert/gopus/rds"
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
	queries, err := rds.Load()
	if err != nil {
		panic(fmt.Errorf("while opening database: %v", err))
	}

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
		path := cmdArgs[0]
		add(queries, path)
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

func add(queries *rds.Queries, path string) {
	fmt.Printf("scanning %s for .opus songs to add to the database", path)

	ctx := context.Background()
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		song, err := queries.AddSong(ctx, rds.AddSongParams{
			Title: "",
			Path:  sql.NullString{},
		})
		if err == nil {
			fmt.Printf("added %s\n", song.Title)
		}
		return err
	})
}

// func find(db *sql.DB, title string) {
// 	fmt.Println("find")
// }
