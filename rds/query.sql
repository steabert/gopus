-- name: AddSong :exec
INSERT INTO song ( title ) VALUES ( ? );

-- name: AddArtist :exec
INSERT INTO artist ( name ) VALUES ( ? );

-- name: AddAlbum :exec
INSERT INTO album ( title ) VALUES ( ? );

-- name: AddRecording :exec
INSERT INTO recording ( path, song, artist, album ) VALUES ( ?, ?, ?, ?);

-- name: ListSongs :many
SELECT song, artist, album, path FROM recording WHERE song LIKE ? ORDER BY song;
