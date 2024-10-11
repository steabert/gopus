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
	info, err := opus.ParseInfo(path)
	if err != nil {
		return fmt.Errorf("failed to read Opus info, %v", err)
	}

	// fmt.Printf("info: %+v\n", info)

	err = rds.Database.AddSong(ctx, info.Comments["TITLE"])
	err = rds.Database.AddAlbum(ctx, info.Comments["ALBUM"])
	err = rds.Database.AddArtist(ctx, info.Comments["ARTIST"])
	err = rds.Database.AddRecording(ctx, rds.AddRecordingParams{
		Path:   path,
		Song:   info.Comments["TITLE"],
		Artist: info.Comments["ARTIST"],
		Album:  info.Comments["ALBUM"],
	})

	if err != nil {
		return fmt.Errorf("failed to add song to database, %v", err)
	}

	return nil
}
