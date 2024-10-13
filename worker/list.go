package worker

import (
	"context"
	"fmt"

	"github.com/steabert/gopus/rds"
)

type Recording struct {
	Path  string
	Song  string
	Album string
	Track int64
}

func MatchSong(pattern string) ([]rds.ListRecordingsMatchingSongRow, error) {
	ctx := context.Background()
	recordings, err := rds.Database.ListRecordingsMatchingSong(ctx, "%"+pattern+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recordings, %v", err)
	}

	return recordings, nil
}

func MatchAlbum(pattern string) ([]rds.ListRecordingsMatchingAlbumRow, error) {
	ctx := context.Background()
	recordings, err := rds.Database.ListRecordingsMatchingAlbum(ctx, "%"+pattern+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recordings, %v", err)
	}

	return recordings, nil
}

func MatchArtist(pattern string) ([]rds.ListRecordingsMatchingArtistRow, error) {
	ctx := context.Background()
	recordings, err := rds.Database.ListRecordingsMatchingArtist(ctx, "%"+pattern+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recordings, %v", err)
	}

	return recordings, nil
}
