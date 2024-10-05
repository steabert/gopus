package worker

import (
	"context"
	"fmt"

	"github.com/steabert/gopus/opus"
	"github.com/steabert/gopus/rds"
)

// InsertSongFromPath adds a song from an .opus file to the database.
func InsertSongFromPath(path string) error {
	ctx := context.Background()
	tags, err := opus.ReadTags(path)
	if err != nil {
		return fmt.Errorf("failed to read Opus tags, %v", err)
	}

	err = rds.Database.AddSong(ctx, rds.AddSongParams{
		Title: tags.Title,
		Path:  &path,
	})

	if err != nil {
		return fmt.Errorf("failed to add song to database, %v", err)
	}

	return nil
}
