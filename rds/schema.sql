CREATE TABLE IF NOT EXISTS song (
    title  TEXT NOT NULL,

    CONSTRAINT PK 
        PRIMARY KEY ( title )
);

CREATE TABLE IF NOT EXISTS artist (
    name  TEXT NOT NULL,

    CONSTRAINT PK 
        PRIMARY KEY ( name )
);

CREATE TABLE IF NOT EXISTS album (
    title  TEXT NOT NULL,
    artist TEXT NOT NULL,

    CONSTRAINT PK 
        PRIMARY KEY ( title )
);

CREATE TABLE IF NOT EXISTS recording (
    path   TEXT NOT NULL,
    song   TEXT NOT NULL,
    artist TEXT NOT NULL,
    album  TEXT NOT NULL,
    cddb   TEXT NOT NULL,
    track  INTEGER NOT NULL,

    CONSTRAINT PK 
        PRIMARY KEY ( path )

    CONSTRAINT artist_recorded_song_on_album_fk
        FOREIGN KEY       ( artist )
        REFERENCES artist ( name )

    CONSTRAINT album_contains_artist_song_fk
        FOREIGN KEY      ( album )
        REFERENCES album ( title )

    CONSTRAINT song_recorded_by_artist_on_album_fk
        FOREIGN KEY     ( song )
        REFERENCES song ( title )
);
