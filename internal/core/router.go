package core

import "github.com/go-chi/chi"

// Router assembles and returns a *chi.Mux router with all the API routes.
func (app *App) Router() *chi.Mux {
	router := chi.NewRouter()
	return router
}
