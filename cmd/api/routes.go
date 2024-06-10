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

	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Get("/bookmarks/{category}", app.GetProjectsByCategory)
	mux.Get("/dashboard", app.Dashboard)
	mux.Get("/auth/{provider}", app.HandleAuth)
	mux.Get("/auth/{provider}/callback", app.HandleCallback)
	mux.Post("/register", app.RegisterNewUser)
	mux.Post("/login", app.ClassicLogin)
	mux.Get("/confirm-email", app.ConfirmEmail)

	// USer information - Feed Dashboard && related screen with user data
	mux.Get("/user-info", app.GetUserInfo)

	mux.Get("/contributors", app.GetContributors)

	// protected route section - now we are not kidding anymore
	mux.Route("/contributor", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Get("/dashboard", app.Dashboard)
		// mux.Post("/{category}/{project}/new-resource", app.InsertNewBookmark)
	})

	return mux
}
