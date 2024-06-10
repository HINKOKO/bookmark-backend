package main

import (
	"bookmarks/internal/repository"
	"bookmarks/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

const port = 8080

type MailConfig struct {
	host     string
	port     int
	username string
	password string
	from     string
}

type OauthConfig struct {
	client_id       string
	client_secret   string
	client_callback string
}

type application struct {
	mailConfig   MailConfig
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	oauth        OauthConfig
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

	smtp_username := os.Getenv("SMTP_USERNAME")
	smtp_password := os.Getenv("SMTP_PASSWORD")
	smtp_from := os.Getenv("SMTP_FROM")

	// Cmd line reading
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=12345 dbname=bookmarkers sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecretstuff", "signing secret for jwt")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "jwt audience")
	// flag.StringVar(&app.CookieDomain, "domain", "localhost", "Cookie domain")
	// Adding smtp mail configuration
	flag.StringVar(&app.mailConfig.host, "smtp host", "sandbox.smtp.mailtrap.io", "smtp host")
	flag.IntVar(&app.mailConfig.port, "smtp port", 2525, "smtp port")
	flag.StringVar(&app.mailConfig.username, "smtp username", smtp_username, "smtp user")
	flag.StringVar(&app.mailConfig.password, "smtp password", smtp_password, "smtp password")
	flag.StringVar(&app.mailConfig.from, "smtp from", smtp_from, "smtp from")
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
		TokenExpiry:   time.Hour * 24,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refresh_token",
		CookieDomain:  "localhost",
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientCallback := os.Getenv("CLIENT_CALLBACK")

	if clientID == "" || clientSecret == "" || clientCallback == "" {
		log.Println("field is missing to param oauth")
		return
	}

	app.oauth = OauthConfig{
		client_id:       clientID,
		client_secret:   clientSecret,
		client_callback: clientCallback,
	}

	goth.UseProviders(
		google.New(clientID, clientSecret, clientCallback),
		github.New(os.Getenv("GITHUB_CLIENT"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_CALLBACK")),
	)

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
