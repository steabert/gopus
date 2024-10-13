-- name: AddSong :exec
INSERT INTO song ( title ) VALUES ( ? );

-- name: AddArtist :exec
INSERT INTO artist ( name ) VALUES ( ? );

-- name: AddAlbum :exec
INSERT INTO album ( title, artist ) VALUES ( ?, ? );

-- name: AddRecording :exec
INSERT INTO recording
  ( path, song, artist, album, cddb, track ) 
VALUES
  ( ?, ?, ?, ?, ?, ?);

-- name: ListRecordingsMatchingSong :many
SELECT path, song, album, track FROM recording WHERE song LIKE ? ORDER BY album, track;

-- name: ListRecordingsMatchingAlbum :many
SELECT path, song, album, track FROM recording WHERE album LIKE ? ORDER BY album, track;

-- name: ListRecordingsMatchingArtist :many
SELECT path, song, album, track FROM recording WHERE artist LIKE ? ORDER BY album, track;
