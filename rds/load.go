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

func Open(mode string) error {
	switch mode {
	case "ro", "rw":
	default:
		return fmt.Errorf("invalid mode: %s, use one of ro, rw\n", mode)
	}

	ctx := context.Background()

	dbConnection := fmt.Sprintf("file:gopus.db?mode=%s", mode)
	db, err := sql.Open("sqlite3", dbConnection)
	if err != nil {
		panic(fmt.Errorf("unable to open database, %v", err))
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return fmt.Errorf("failed to access database, %v", err)
	}

	Database = New(db)
	return nil
}
