package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	// mux.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:8080"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300,
	// }))
	mux.Use(app.enableCORS)

	mux.Get("/bookmarks/{category}", app.GetProjectsByCategory)
	// mux.Post("/login", app.Login)
	mux.Post("/register", app.Register)

	return mux
}
