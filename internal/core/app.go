package core

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joinimpact/api/internal/authentication"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/conversations"
	authm "github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/internal/tags"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/internal/websocket/socketserver"
	"github.com/rs/zerolog"
)

// APIRevision defines the whole-number version of the API, used for the /api/v{version} route.
const APIRevision = 1

// App represents a servable app.
type App struct {
	config                *config.Config
	logger                *zerolog.Logger
	websocketService      socketserver.Service
	authenticationService authentication.Service
	usersService          users.Service
	organizationsService  organizations.Service
	tagsService           tags.Service
	opportunitiesService  opportunities.Service
	eventsService         events.Service
	conversationsService  conversations.Service
}

// NewApp creates and returns a new *App with the provided Config.
func NewApp(config *config.Config, logger *zerolog.Logger, websocketService socketserver.Service, authenticationService authentication.Service, usersService users.Service, organizationsService organizations.Service, tagsService tags.Service, opportunitiesService opportunities.Service, eventsService events.Service, conversationsService conversations.Service) *App {
	return &App{
		config,
		logger,
		websocketService,
		authenticationService,
		usersService,
		organizationsService,
		tagsService,
		opportunitiesService,
		eventsService,
		conversationsService,
	}
}

// Serve serves the App on the port specified in the config.
func (app *App) Serve() error {
	// Create a new router.
	router := chi.NewRouter()

	// Apply the Logger middleware if dev mode is enabled.
	if app.config.DevMode {
		router.Use(middleware.Logger)
	}

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "docs", "swagger-ui"))
	FileServer(router, "/swagger-ui", filesDir)

	router.Group(func(router chi.Router) {
		// JSON middleware
		router.Use(middleware.SetHeader("Content-Type", "application/json"))

		// Add the healthcheck.
		router.Get("/healthcheck", healthcheckHandler)

		// Mount the API router at /api/v1.
		router.Mount(fmt.Sprintf("/api/v%d", APIRevision), app.Router())

		// Mount the WebSocket handler at /ws/v1.
		router.
			With(authm.CookieMiddleware(app.authenticationService)).
			Mount(fmt.Sprintf("/ws/v%d", APIRevision), app.websocketService.Handler())
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", app.config.Port), router)
}
