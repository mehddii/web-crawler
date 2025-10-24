package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Persister interface {
	Save(uri string, text string, metadata string) error
}

type DBConnection struct {
	file string
	db   *sql.DB
}

func NewConnection(file string) (*DBConnection, error) {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	return &DBConnection{
		file: file,
		db:   db,
	}, nil
}

type SqlitePersister struct {
	dbc   *DBConnection
	table string
}

func NewSqlitePersister(file string, table string) (*SqlitePersister, error) {
	dbc, err := NewConnection(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Database connection was established successfully")

	return &SqlitePersister{
		dbc:   dbc,
		table: table,
	}, nil
}

func (persister *SqlitePersister) Save(uri string, text string, metadata string) error {
	err := initializeDB(persister.dbc, persister.table)
	if err != nil {
		log.Println(err)
		return err
	}

	query := fmt.Sprintf(
		"INSERT INTO %q(uri, text, metadata) VALUES(%q, %q, %q)",
		persister.table,
		uri,
		text,
		metadata,
	)

	_, err = persister.dbc.db.Exec(query)
	if err != nil {
		log.Println("")
		return err
	}
	log.Println(uri, "was inserted successfully.")

	return nil
}

func initializeDB(dbc *DBConnection, table string) error {
	if dbc == nil {
		return errors.New("No Database connection detected. Operation failed.")
	}

	err := dbc.db.Ping()
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %q (
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			uri TEXT NOT NULL,
			text TEXT,
			metadata TEXT
		);
	`, table)

	_, err = dbc.db.Exec(query)
	if err != nil {
		return err
	}
	log.Printf("Table %s was created successfully", table)

	return nil
}
