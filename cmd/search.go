package main

import (
	"errors"
	"fmt"

	"github.com/steabert/gopus/rds"
	"github.com/steabert/gopus/worker"
)

func list(args []string) error {
	var err error

	if len(args) != 1 {
		usage()
		return errors.New("no arguments, expected pattern to match")
	}

	err = rds.Open("ro")
	if err != nil {
		return fmt.Errorf("failed to open database, %v", err)
	}

	recordings, err := worker.MatchSong(args[0])
	if err != nil {
		return fmt.Errorf("failed to retrieve recordings, %v", err)
	}

	for _, recording := range recordings {
		fmt.Println(recording)
	}

	return nil
}
