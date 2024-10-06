package rds

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

var Database *Queries

//go:embed schema.sql
var ddl string

func init() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "file:gopus.db")
	if err != nil {
		panic(fmt.Errorf("unable to open database, %v", err))
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		// panic(fmt.Errorf("failed to access database, %v", err))
	}

	Database = New(db)
}
