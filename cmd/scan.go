package main

import (
	"errors"
	"fmt"

	"github.com/steabert/gopus/rds"
	"github.com/steabert/gopus/worker"
)

func scan(args []string) error {
	var err error

	if len(args) == 0 {
		usage()
		return errors.New("no arguments, expected at least 1 directory to search")
	}

	err = rds.Open("rwc")
	if err != nil {
		return fmt.Errorf("failed to open database, %v", err)
	}

	for _, dir := range args {
		err = worker.WalkDirInsert(dir)
		if err != nil {
			return fmt.Errorf("failed to scan directory, %v", err)
		}
	}

	return nil
}
