package rds

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

//go:embed schema.sql
var ddl string

func Load() (*Queries, error) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "file:gopus.db")
	if err != nil {
		return nil, fmt.Errorf("while opening database: %v", err)
	}
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, fmt.Errorf("while creating db tables: %v", err)
	}
	queries := New(db)
	return queries, nil
}
