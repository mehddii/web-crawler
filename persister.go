package main

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type WebDocument struct {
	URI      string
	Text     string
	Metadata string // --todo: should probably have its own type
}

type Persister interface {
	Persist(ctx context.Context, doc WebDocument) error
}

type SQLPersister struct {
	DB *sql.DB
}

func (dbc *SQLPersister) Persist(ctx context.Context, doc WebDocument) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	_, err = dbc.DB.ExecContext(ctx, `
		INSERT INTO documents(id, uri, text, metadata)	
		VALUES (?, ?, ?, ?);
	`, id, doc.URI, doc.Text, doc.Metadata)

	return nil
}

func Init(ctx context.Context, source string) (Persister, error) {
	db, err := sql.Open("sqlite3", source)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS documents (
			id TEXT PRIMARY KEY,
			uri TEXT NOT NULL,
			text TEXT,
			metadata TEXT
		);
	`)
	if err != nil {
		return nil, err
	}

	return &SQLPersister{DB: db}, nil
}
