package main

import (
	"os"

	"github.com/joinimpact/api/internal/models"

	"github.com/joinimpact/api/internal/migrations"

	"github.com/joinimpact/api/internal/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// APIVersion defines the current version of the API.
const APIVersion = "1.0.0"

func main() {
	// Create a new default configuration based on environment variables.
	config := core.NewConfig()
	if config.DevMode {
		// If in dev mode, enable pretty printing.
		// This makes the logger less efficient, so it is only
		// enabled for dev mode.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Create a DatabaseService and connect to the database.
	databaseService := core.NewDatabaseService(config, &log.Logger)
	db, err := databaseService.DatabaseConnect()
	if err != nil {
		// Error connecting to the database, panic.
		log.Panic().Err(err).Msg("Error connecting to the database")
	}
	// Create a new MigrationService to handle automatic migrations.
	migrationService := migrations.NewMigrationService(db)
	// Auto migrate models into the database.
	migrationService.Migrate(
		&models.User{},
		&models.UserProfileField{},
		&models.Organization{},
		&models.OrganizationProfileField{},
		&models.OrganizationMembership{},
	)
	// Create a new app using the new config.
	app := core.NewApp(config, &log.Logger)

	// Print a message.
	log.Info().Int("port", int(config.Port)).Str("version", APIVersion).Msg("Listening")

	// Serve the app, panic if an error occurs.
	panic(app.Serve())
}
