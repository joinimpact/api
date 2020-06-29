package core

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joinimpact/api/internal/authentication"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/internal/tags"
	"github.com/joinimpact/api/internal/users"
	"github.com/rs/zerolog"
)

// APIRevision defines the whole-number version of the API, used for the /api/v{version} route.
const APIRevision = 1

// App represents a servable app.
type App struct {
	config                *config.Config
	logger                *zerolog.Logger
	authenticationService authentication.Service
	usersService          users.Service
	organizationsService  organizations.Service
	tagsService           tags.Service
}

// NewApp creates and returns a new *App with the provided Config.
func NewApp(config *config.Config, logger *zerolog.Logger, authenticationService authentication.Service, usersService users.Service, organizationsService organizations.Service, tagsService tags.Service) *App {
	return &App{
		config,
		logger,
		authenticationService,
		usersService,
		organizationsService,
		tagsService,
	}
}

// Serve serves the App on the port specified in the config.
func (app *App) Serve() error {
	// Create a new router.
	router := chi.NewRouter()

	// JSON middleware
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	// Apply the Logger middleware if dev mode is enabled.
	if app.config.DevMode {
		router.Use(middleware.Logger)
	}

	// Add the healthcheck.
	router.Get("/healthcheck", healthcheckHandler)

	// Mount the API router at /api/v1
	router.Mount(fmt.Sprintf("/api/v%d", APIRevision), app.Router())

	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), router)
}
