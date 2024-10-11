package worker

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func ScanDirectory(dir string) error {
	fmt.Printf("scanning %s for Opus files to add to the database...\n", dir)

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		err = InsertSongFromPath(path)
		if err != nil {
			fmt.Printf("[ERROR] failed to add %s, %v\n", path, err)
		} else {
			fmt.Printf("[OK] added %s\n", path)
		}
		return nil
	})

	return err
}
