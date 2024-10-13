package worker

import (
	"context"
	"fmt"
	"strconv"

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

	err = rds.Database.AddSong(ctx, info.Comments["TITLE"])
	err = rds.Database.AddAlbum(ctx, rds.AddAlbumParams{
		Title:  info.Comments["ALBUM"],
		Artist: info.Comments["ALBUMARTIST"],
	})
	err = rds.Database.AddArtist(ctx, info.Comments["ARTIST"])

	track, err := strconv.Atoi(info.Comments["TRACKNUMBER"])
	if err != nil {
		return fmt.Errorf("invalid track number, %v", err)
	}
	err = rds.Database.AddRecording(ctx, rds.AddRecordingParams{
		Path:   path,
		Song:   info.Comments["TITLE"],
		Artist: info.Comments["ARTIST"],
		Album:  info.Comments["ALBUM"],
		Cddb:   info.Comments["CDDB"],
		Track:  int64(track),
	})

	if err != nil {
		return fmt.Errorf("failed to add song to database, %v", err)
	}

	return nil
}
