-- name: AddSong :exec
INSERT INTO songs ( title, path ) VALUES ( ?, ? );

-- name: ListTitles :many
SELECT title, path FROM songs WHERE title LIKE ? ORDER BY title;

-- name: GetTitle :many
SELECT title, path FROM songs WHERE id = ?;
