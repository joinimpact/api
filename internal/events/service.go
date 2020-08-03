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
	// UpdateEvent updates an event.
	UpdateEvent(ctx context.Context, request ModifyEventRequest) error
	// GetEvent gets a single event by ID.
	GetEvent(ctx context.Context, eventID int64) (*EventView, error)
	// GetMinimalEvent gets a single event with only the event base fields.
	GetMinimalEvent(ctx context.Context, eventID int64) (*EventView, error)
	// GetEventResponses gets all responses to an event.
	GetEventResponses(ctx context.Context, eventID int64) ([]models.EventResponse, error)
	// GetUserEventResponse gets a user's response to a single event.
	GetUserEventResponse(ctx context.Context, userID, eventID int64) (*models.EventResponse, error)
	// SetEventResponseCanAttend creates or updates an EventResponse with the "can attend" status.
	SetEventResponseCanAttend(ctx context.Context, eventID, userID int64) error
	// SetEventResponseCanNotAttend creates or updates an EventResponse with the "can not attend" status.
	SetEventResponseCanNotAttend(ctx context.Context, eventID, userID int64) error
	// GetOpportunityEvents gets all events by opportunity ID.
	GetOpportunityEvents(ctx context.Context, opportunityID int64) ([]EventView, error)
	// GetUserEvents gets events from all of a user's enrolled opportunities.
	GetUserEvents(ctx context.Context, userID int64) ([]EventView, error)
	// DeleteEvent deletes a single event by ID.
	DeleteEvent(ctx context.Context, eventID int64) error
}

// service represents the internal implementation of the Service.
type service struct {
	eventRepository                 models.EventRepository
	eventResponseRepository         models.EventResponseRepository
	opportunityMembershipRepository models.OpportunityMembershipRepository
	tagRepository                   models.TagRepository
	config                          *config.Config
	logger                          *zerolog.Logger
	snowflakeService                snowflakes.SnowflakeService
	emailService                    email.Service
	cdnClient                       *cdn.Client
	locationService                 location.Service
}

// NewService creates and returns a new events.Service with the provided
// dependencies.
func NewService(eventRepository models.EventRepository, eventResponseRepository models.EventResponseRepository, opportunityMembershipRepository models.OpportunityMembershipRepository, tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service, locationService location.Service) Service {
	return &service{
		eventRepository,
		eventResponseRepository,
		opportunityMembershipRepository,
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

// UpdateEvent updates an event.
func (s *service) UpdateEvent(ctx context.Context, request ModifyEventRequest) error {
	event := s.requestToEvent(request)

	// Validate the event to the minimum requirements.
	if !validateEvent(&event) {
		return NewErrServerError()
	}

	err := s.eventRepository.Update(ctx, event)
	if err != nil {
		return NewErrEventNotFound()
	}

	return nil
}

// GetEvent gets a single event by ID.
func (s *service) GetEvent(ctx context.Context, eventID int64) (*EventView, error) {
	// Find the event by ID.
	event, err := s.eventRepository.FindByID(ctx, eventID)
	if err != nil {
		return nil, NewErrEventNotFound()
	}

	// Convert the event to view.
	view, err := s.eventToView(*event)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error converting event to view")
		return nil, NewErrServerError()
	}

	responsesSummary, err := s.getEventResponsesSummary(ctx, eventID)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error getting EventResponsesSummary")
		return nil, NewErrServerError()
	}
	view.EventResponsesSummary = responsesSummary

	return view, nil
}

// GetMinimalEvent gets a single event by ID.
func (s *service) GetMinimalEvent(ctx context.Context, eventID int64) (*EventView, error) {
	// Find the event by ID.
	event, err := s.eventRepository.FindByID(ctx, eventID)
	if err != nil {
		return nil, NewErrEventNotFound()
	}

	// Convert the event to view.
	view, err := s.eventToMinimalView(*event)
	if err != nil {
		return nil, NewErrServerError()
	}

	return view, nil
}

// GetEventResponses gets all responses to an event.
func (s *service) GetEventResponses(ctx context.Context, eventID int64) ([]models.EventResponse, error) {
	// Get all event memberships.
	ids, err := s.getEventMemberships(ctx, eventID)
	if err != nil {
		return nil, NewErrServerError()
	}

	responses := []models.EventResponse{}
	for _, id := range ids {
		code := models.EventResponseNull
		// Create a default response to append when one is not found in the database.
		response := &models.EventResponse{
			UserID:   id,
			EventID:  eventID,
			Response: &code,
		}

		// Find a response if one exists.
		eventResponse, err := s.eventResponseRepository.FindInEventByUserID(ctx, eventID, id)
		if err == nil {
			// If a response exists, overwrite the default response.
			response = eventResponse
		}

		responses = append(responses, *response)
	}

	return responses, nil
}

// GetUserEventResponse gets a user's response to a single event.
func (s *service) GetUserEventResponse(ctx context.Context, userID, eventID int64) (*models.EventResponse, error) {
	response, err := s.eventResponseRepository.FindInEventByUserID(ctx, eventID, userID)
	if err != nil {
		return nil, NewErrResponseNotFound()
	}

	return response, nil
}

// SetEventResponseCanAttend creates or updates an EventResponse with the "can attend" status.
func (s *service) SetEventResponseCanAttend(ctx context.Context, eventID, userID int64) error {
	code := models.EventResponseCanAttend
	return s.setEventResponse(ctx, eventID, userID, code)
}

// SetEventResponseCanNotAttend creates or updates an EventResponse with the "can not attend" status.
func (s *service) SetEventResponseCanNotAttend(ctx context.Context, eventID, userID int64) error {
	code := models.EventResponseCanNotAttend
	return s.setEventResponse(ctx, eventID, userID, code)
}

// setEventResponse creates or updates an event response.
func (s *service) setEventResponse(ctx context.Context, eventID, userID int64, code int) error {
	// Check if an event response already exists.
	response, err := s.eventResponseRepository.FindInEventByUserID(ctx, eventID, userID)
	if err == nil {
		// If found, update the existing response.
		response.Response = &code
		err := s.eventResponseRepository.Update(ctx, *response)
		if err != nil {
			return NewErrServerError()
		}

		return nil
	}

	// Existing response not found, create a new one.
	id := s.snowflakeService.GenerateID()
	response = &models.EventResponse{}
	response.ID = id
	response.EventID = eventID
	response.UserID = userID
	response.Response = &code

	err = s.eventResponseRepository.Create(ctx, *response)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}

// GetOpportunityEvents gets all events by opportunity ID.
func (s *service) GetOpportunityEvents(ctx context.Context, opportunityID int64) ([]EventView, error) {
	// Find the events by opportunity ID.
	events, err := s.eventRepository.FindByOpportunityID(ctx, opportunityID)
	if err != nil {
		return nil, NewErrServerError()
	}

	views := []EventView{}
	for _, event := range events {
		// Convert the event to view.
		view, err := s.eventToView(event)
		if err != nil {
			return nil, NewErrServerError()
		}

		views = append(views, *view)
	}

	return views, nil
}

// getEventMemberships gets the IDs of all users who are members of an event.
func (s *service) getEventMemberships(ctx context.Context, eventID int64) ([]int64, error) {
	// Get the event.
	event, err := s.eventRepository.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	// Get all memberships of the event's opportunity.
	memberships, err := s.opportunityMembershipRepository.FindByOpportunityID(ctx, event.OpportunityID)
	if err != nil {
		return nil, err
	}

	// Reduce all memberships to user IDs.
	userIDs := []int64{}
	for _, membership := range memberships {
		userIDs = append(userIDs, membership.UserID)
	}

	return userIDs, nil
}

// GetUserEvents gets events from all of a user's enrolled opportunities.
func (s *service) GetUserEvents(ctx context.Context, userID int64) ([]EventView, error) {
	memberships, err := s.opportunityMembershipRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	ids := []int64{}
	for _, membership := range memberships {
		ids = append(ids, membership.OpportunityID)
	}

	events, err := s.eventRepository.FindByOpportunityIDs(ctx, ids)
	if err != nil {
		return nil, NewErrServerError()
	}

	views := []EventView{}
	for _, event := range events {
		// Convert the event to view.
		view, err := s.eventToView(event)
		if err != nil {
			return nil, NewErrServerError()
		}

		views = append(views, *view)
	}

	return views, nil
}

// DeleteEvent deletes a single event by ID.
func (s *service) DeleteEvent(ctx context.Context, eventID int64) error {
	if err := s.eventRepository.DeleteByID(ctx, eventID); err != nil {
		return NewErrServerError()
	}

	return nil
}
