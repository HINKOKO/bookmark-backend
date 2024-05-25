package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to connect to db: %v\n", err)
		os.Exit(1)
	}

	// Ping the database
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	conn, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres")
	return conn, nil
}
