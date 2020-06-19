package core

import (
	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/core/handlers/auth"
)

// Router assembles and returns a *chi.Mux router with all the API routes.
func (app *App) Router() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthcheck", healthcheckHandler)
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.Login(app.authenticationService))
		r.Post("/register", auth.Register(app.authenticationService))
	})

	return router
}
