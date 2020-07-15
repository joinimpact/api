package events

import (
	"context"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/pkg/location"
	"github.com/rs/zerolog"
)

// Service provides methods for interacting with events within opportunities.
type Service interface {
	// CreateEvent creates an event and returns the ID of the newly created event.
	CreateEvent(ctx context.Context, request ModifyEventRequest) (int64, error)
}

// service represents the internal implementation of the Service.
type service struct {
	eventRepository  models.EventRepository
	tagRepository    models.TagRepository
	config           *config.Config
	logger           *zerolog.Logger
	snowflakeService snowflakes.SnowflakeService
	emailService     email.Service
	cdnClient        *cdn.Client
	locationService  location.Service
}

// NewService creates and returns a new events.Service with the provided
// dependencies.
func NewService(eventRepository models.EventRepository, tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service, locationService location.Service) Service {
	return &service{
		eventRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		cdn.NewCDNClient(config),
		locationService,
	}
}

// CreateEvent creates an event and returns the ID of the newly created
// event.
func (s *service) CreateEvent(ctx context.Context, request ModifyEventRequest) (int64, error) {
	event := s.requestToEvent(request)

	// Validate the event to the minimum requirements.
	if !validateEvent(&event) {
		return 0, NewErrServerError()
	}

	// Generate an ID for the event.
	event.ID = s.snowflakeService.GenerateID()

	err := s.eventRepository.Create(ctx, event)
	if err != nil {
		return 0, NewErrEventNotFound()
	}

	return event.ID, nil
}

// GetEvent gets a single event by ID.
func (s *service) GetEvent(ctx context.Context, eventID int64) (*EventView, error) {
	return nil, nil
}

// GetOpportunityEvents gets all events by opportunity ID.
func (s *service) GetOpportunityEvents(ctx context.Context, opportunityID int64) ([]EventView, error) {
	return nil, nil
}
