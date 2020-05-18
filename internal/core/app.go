package core

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

// App represents a servable app.
type App struct {
	config *Config
	logger *zerolog.Logger
}

// NewApp creates and returns a new *App with the provided Config.
func NewApp(config *Config, logger *zerolog.Logger) *App {
	return &App{config, logger}
}

// Serve serves the App on the port specified in the config.
func (app *App) Serve() error {
	// Create a new router.
	router := chi.NewRouter()
	// Apply the Logger middleware if dev mode is enabled.
	if app.config.DevMode {
		router.Use(middleware.Logger)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), router)
}
