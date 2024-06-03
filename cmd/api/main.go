package main

import (
	"bookmarks/internal/repository"
	"bookmarks/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
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

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error: no .env file found. Shutting down")
	}
}

func main() {
	var app application

	// Cmd line reading
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=12345 dbname=bookmarkers sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecretstuff", "signing secret for jwt")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "jwt audience")
	// flag.StringVar(&app.CookieDomain, "domain", "example.com", "Cookie domain")
	flag.Parse()

	// Connect to DB
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// populate releavant field of application struct
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Second * 45,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refresh_token",
		CookieDomain:  "localhost",
	}

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
