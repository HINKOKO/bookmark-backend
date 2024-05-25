package main

import (
	"bookmarks/internal/repository"
	"bookmarks/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
}

func main() {
	var app application

	// Cmd line reading
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=12345 dbname=bookmarkers sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()

	// Connect to DB
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// populate releavant field of application struct
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
