package main

import (
	"os"

	"github.com/joinimpact/api/internal/authentication"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/conversations"
	"github.com/joinimpact/api/internal/database/postgres"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/events"
	"github.com/joinimpact/api/internal/hours"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/internal/pubsub"
	"github.com/joinimpact/api/internal/search"
	opportunitiesSearch "github.com/joinimpact/api/internal/search/stores/opportunities"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/internal/tags"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/internal/websocket/hub"
	"github.com/joinimpact/api/internal/websocket/hubmanager"
	"github.com/joinimpact/api/internal/websocket/socketserver"
	"github.com/joinimpact/api/pkg/location"

	"github.com/joinimpact/api/internal/migrations"

	"github.com/joinimpact/api/internal/core"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// APIVersion defines the current version of the API.
const APIVersion = "1.0.0"

func main() {
	// Create a new default configuration based on environment variables.
	config := config.NewConfig()
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
	// Close the database on program exit.
	defer db.Close()

	if config.DevMode {
		db.LogMode(true)
	}

	// Create a new MigrationService to handle automatic migrations.
	migrationService := migrations.NewMigrationService(db)
	// Auto migrate models into the database.
	err = migrationService.Migrate(
		&models.User{},
		&models.UserProfileField{},
		&models.Organization{},
		&models.OrganizationProfileField{},
		&models.OrganizationMembership{},
		&models.OrganizationMembershipInvite{},
		&models.OrganizationTag{},
		&models.Opportunity{},
		&models.OpportunityLimits{},
		&models.OpportunityRequirements{},
		&models.OpportunityTag{},
		&models.OpportunityMembership{},
		&models.OpportunityMembershipRequest{},
		&models.OpportunityMembershipInvite{},
		&models.Event{},
		&models.EventResponse{},
		&models.Conversation{},
		&models.ConversationOpportunityMembershipRequest{},
		&models.ConversationMembership{},
		&models.Message{},
		&models.PasswordResetKey{},
		&models.UserTag{},
		&models.Tag{},
		&models.ThirdPartyIdentity{},
		&models.VolunteeringHourLog{},
		&models.VolunteeringHourLogRequest{},
	)
	if err != nil {
		// Error migrating the database, panic.
		log.Fatal().Err(err).Msg("Error migrating the database")
	}

	// Dependencies/external services
	snowflakeService, err := snowflakes.NewSnowflakeService()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing snowflake service")
	}
	emailService := email.NewService(config, email.NewSender(
		"Impact",
		"no-reply@joinimpact.org",
	))
	locationService, err := location.NewService(&location.Options{
		APIKey: config.GoogleMapsAPIKey,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing Google Maps API")
	}

	// Repositories
	userRepository := postgres.NewUserRepository(db, &log.Logger)
	passwordResetRepository := postgres.NewPasswordResetRepository(db, &log.Logger)
	userProfileFieldRepository := postgres.NewUserProfileFieldRepository(db, &log.Logger)
	thirdPartyIdentityRepository := postgres.NewThirdPartyIdentityRepository(db, &log.Logger)
	userTagRepository := postgres.NewUserTagRepository(db, &log.Logger)
	tagRepository := postgres.NewTagRepository(db, &log.Logger)
	organizationRepository := postgres.NewOrganizationRepository(db, &log.Logger)
	organizationMembershipRepository := postgres.NewOrganizationMembershipRepository(db, &log.Logger)
	organizationMembershipInviteRepository := postgres.NewOrganizationMembershipInviteRepository(db, &log.Logger)
	organizationProfileFieldRepository := postgres.NewOrganizationProfileFieldRepository(db, &log.Logger)
	organizationTagRepository := postgres.NewOrganizationTagRepository(db, &log.Logger)
	opportunityRepository := postgres.NewOpportunityRepository(db, &log.Logger)
	opportunityRequirementsRepository := postgres.NewOpportunityRequirementsRepository(db, &log.Logger)
	opportunityLimitsRepository := postgres.NewOpportunityLimitsRepository(db, &log.Logger)
	opportunityTagRepository := postgres.NewOpportunityTagRepository(db, &log.Logger)
	opportunityMembershipRepository := postgres.NewOpportunityMembershipRepository(db, &log.Logger)
	opportunityMembershipRequestRepository := postgres.NewOpportunityMembershipRequestRepository(db, &log.Logger)
	opportunityMembershipInviteRepository := postgres.NewOpportunityMembershipInviteRepository(db, &log.Logger)
	conversationRepository := postgres.NewConversationRepository(db, &log.Logger)
	conversationMembershipRepository := postgres.NewConversationMembershipRepository(db, &log.Logger)
	conversationOpportunityMembershipRequestRepository := postgres.NewConversationOpportunityMembershipRequestRepository(db, &log.Logger)
	conversationOrganizationMembershipRepository := postgres.NewConversationOrganizationMembershipRepository(db, &log.Logger)
	messageRepository := postgres.NewMessageRepository(db, &log.Logger)
	eventRepository := postgres.NewEventRepository(db, &log.Logger)
	eventResponseRepository := postgres.NewEventResponseRepository(db, &log.Logger)
	volunteeringHourLogRepository := postgres.NewVolunteeringHourLogRepository(db, &log.Logger)
	volunteeringHourLogRequestRepository := postgres.NewVolunteeringHourLogRequestRepository(db, &log.Logger)

	// Elastic client
	elasticClient, err := search.NewElasticsearch(config.ElasticHost, config.ElasticPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing Elasticsearch client")
	}
	// Elastic services
	opportunitiesSearchService := opportunitiesSearch.NewStore(elasticClient, opportunityRepository, opportunityRequirementsRepository, opportunityLimitsRepository, opportunityTagRepository, tagRepository, organizationRepository, &log.Logger, config)
	err = opportunitiesSearchService.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting Opportunities search service")
	}

	// Pub/sub service
	broker := pubsub.NewBroker()

	// Internal services
	usersService := users.NewService(userRepository, userProfileFieldRepository, userTagRepository, tagRepository, config, &log.Logger, snowflakeService, locationService)
	authenticationService := authentication.NewService(userRepository, passwordResetRepository, thirdPartyIdentityRepository, config, &log.Logger, snowflakeService, emailService)
	organizationsService := organizations.NewService(organizationRepository, organizationMembershipRepository, organizationMembershipInviteRepository, organizationProfileFieldRepository, organizationTagRepository, userRepository, tagRepository, config, &log.Logger, snowflakeService, emailService, locationService)
	opportunitiesService := opportunities.NewService(opportunityRepository, opportunityRequirementsRepository, opportunityLimitsRepository, opportunityTagRepository, opportunityMembershipRepository, opportunityMembershipRequestRepository, opportunityMembershipInviteRepository, tagRepository, userRepository, userTagRepository, organizationRepository, config, &log.Logger, snowflakeService, emailService, opportunitiesSearchService, locationService)
	eventsService := events.NewService(eventRepository, eventResponseRepository, opportunityMembershipRepository, tagRepository, config, &log.Logger, snowflakeService, emailService, locationService)
	conversationsService := conversations.NewService(conversationRepository, conversationMembershipRepository, conversationOpportunityMembershipRequestRepository, conversationOrganizationMembershipRepository, messageRepository, opportunityRepository, userRepository, userProfileFieldRepository, userTagRepository, tagRepository, volunteeringHourLogRequestRepository, config, &log.Logger, snowflakeService, emailService, broker, locationService)
	tagsService := tags.NewService(tagRepository, config, &log.Logger, snowflakeService)
	hoursService := hours.NewService(volunteeringHourLogRepository, volunteeringHourLogRequestRepository, opportunityRepository, organizationRepository,
		userRepository, eventRepository, config, &log.Logger, snowflakeService, emailService, broker, locationService)

	// WebSocket services
	wsHub := hub.NewHub(hub.Options{})
	hubManager := hubmanager.NewHubManager(wsHub)
	wsManager := socketserver.NewWebSocketManager(wsHub, hubManager, broker, authenticationService, organizationsService, conversationsService)
	wsManager.SubscribeHub()
	websocketService := socketserver.NewService(wsManager)

	// Create a new app using the new config.
	app := core.NewApp(config, &log.Logger, websocketService, authenticationService, usersService, organizationsService, tagsService, opportunitiesService, eventsService, conversationsService, hoursService)

	// Print a message.
	log.Info().Int("port", int(config.Port)).Str("version", APIVersion).Msg("Listening")

	// Serve the app, panic if an error occurs.
	panic(app.Serve())
}
