package hours

import (
	"context"
	"time"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/pubsub"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/pkg/dbctx"
	"github.com/joinimpact/api/pkg/location"
	"github.com/rs/zerolog"
)

// Service represents a service for tracking volunteer hours.
type Service interface {
	// RequestHours requests hour validation from an organization as a volunteer.
	RequestHours(ctx context.Context, volunteerID, organizationID int64, requestedHours float32) (int64, error)
	// GetOrganizationRequests gets all volunteer hour requests per organization.
	GetOrganizationRequests(ctx context.Context, organizationID int64) (*VolunteeringHourLogRequestsResponse, error)
	// GetOpportunityRequests gets all volunteer hour requests per opportunity.
	GetOpportunityRequests(ctx context.Context, opportunityID int64) (*VolunteeringHourLogRequestsResponse, error)
	// AcceptRequest accepts a request by ID.
	AcceptRequest(ctx context.Context, granterID, requestID int64) error
	// DeclineRequest declines a request by ID.
	DeclineRequest(ctx context.Context, granterID, requestID int64) error
	// GetHoursByVolunteer gets a user's hours.
	GetHoursByVolunteer(ctx context.Context, volunteerID int64) (*VolunteeringHourLogsResponse, error)
}

// service is the internal implementation of the hours.Service interface.
type service struct {
	volunteeringHourLogRepository        models.VolunteeringHourLogRepository
	volunteeringHourLogRequestRepository models.VolunteeringHourLogRequestRepository
	opportunityRepository                models.OpportunityRepository
	organizationRepository               models.OrganizationRepository
	userRepository                       models.UserRepository
	eventRepository                      models.EventRepository
	config                               *config.Config
	logger                               *zerolog.Logger
	snowflakeService                     snowflakes.SnowflakeService
	emailService                         email.Service
	broker                               pubsub.Broker
	locationService                      location.Service
	cdnClient                            *cdn.Client
}

// NewService creates and returns a new hours.Service.
func NewService(volunteeringHourLogRepository models.VolunteeringHourLogRepository, volunteeringHourLogRequestRepository models.VolunteeringHourLogRequestRepository, opportunityRepository models.OpportunityRepository, organizationRepository models.OrganizationRepository, userRepository models.UserRepository, eventRepository models.EventRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service, broker pubsub.Broker, locationService location.Service) Service {
	return &service{
		volunteeringHourLogRepository,
		volunteeringHourLogRequestRepository,
		opportunityRepository,
		organizationRepository,
		userRepository,
		eventRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		broker,
		locationService,
		cdn.NewCDNClient(config),
	}
}

// RequestHours requests hour validation from an organization as a volunteer.
func (s *service) RequestHours(ctx context.Context, volunteerID, organizationID int64, requestedHours float32) (int64, error) {
	volunteeringHourLogRequest := models.VolunteeringHourLogRequest{}
	volunteeringHourLogRequest.Accepted = false
	volunteeringHourLogRequest.ID = s.snowflakeService.GenerateID()
	volunteeringHourLogRequest.VolunteerID = volunteerID
	volunteeringHourLogRequest.OrganizationID = organizationID
	volunteeringHourLogRequest.RequestedHours = requestedHours

	if err := s.volunteeringHourLogRequestRepository.Create(ctx, volunteeringHourLogRequest); err != nil {
		return 0, NewErrServerError()
	}

	return volunteeringHourLogRequest.ID, nil
}

// VolunteeringHourLogRequestsResponse contains volunteering hour log requests and a total number of pages.
type VolunteeringHourLogRequestsResponse struct {
	VolunteeringHourLogRequests []models.VolunteeringHourLogRequest `json:"hourLogRequests"`
	Pages                       uint                                `json:"pages"`
}

// VolunteeringHourLogsResponse contains volunteering hour logs and a total number of pages.
type VolunteeringHourLogsResponse struct {
	VolunteeringHourLogs []models.VolunteeringHourLog `json:"hourLogs"`
	Pages                uint                         `json:"pages"`
}

// GetOrganizationRequests gets all volunteer hour requests per organization.
func (s *service) GetOrganizationRequests(ctx context.Context, organizationID int64) (*VolunteeringHourLogRequestsResponse, error) {
	res, err := s.volunteeringHourLogRequestRepository.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return &VolunteeringHourLogRequestsResponse{
		VolunteeringHourLogRequests: res.VolunteeringHourLogRequests,
		Pages:                       uint(res.TotalResults/dbctx.Get(ctx).Limit) + 1,
	}, nil
}

// GetOpportunityRequests gets all volunteer hour requests per opportunity.
func (s *service) GetOpportunityRequests(ctx context.Context, opportunityID int64) (*VolunteeringHourLogRequestsResponse, error) {
	res, err := s.volunteeringHourLogRequestRepository.FindByOpportunityID(ctx, opportunityID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return &VolunteeringHourLogRequestsResponse{
		VolunteeringHourLogRequests: res.VolunteeringHourLogRequests,
		Pages:                       uint(res.TotalResults/dbctx.Get(ctx).Limit) + 1,
	}, nil
}

// AcceptRequest accepts a request by ID.
func (s *service) AcceptRequest(ctx context.Context, granterID, requestID int64) error {
	request, err := s.volunteeringHourLogRequestRepository.FindByID(ctx, requestID)
	if err != nil {
		return NewErrRequestNotFound()
	}

	hourLog := models.VolunteeringHourLog{}
	hourLog.ID = s.snowflakeService.GenerateID()
	hourLog.VolunteerID = request.VolunteerID
	hourLog.OrganizationID = request.OrganizationID
	hourLog.OpportunityID = request.OpportunityID
	hourLog.GrantedHours = request.RequestedHours
	hourLog.GranterID = granterID
	hourLog.GrantedOn = time.Now()

	err = s.volunteeringHourLogRepository.Create(ctx, hourLog)
	if err != nil {
		return NewErrServerError()
	}

	err = s.volunteeringHourLogRequestRepository.DeleteByID(ctx, request.ID)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}

// DeclineRequest declines a request by ID.
func (s *service) DeclineRequest(ctx context.Context, granterID, requestID int64) error {
	request, err := s.volunteeringHourLogRequestRepository.FindByID(ctx, requestID)
	if err != nil {
		return NewErrRequestNotFound()
	}

	err = s.volunteeringHourLogRequestRepository.DeleteByID(ctx, request.ID)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}

// GetHoursByVolunteer gets a user's hours.
func (s *service) GetHoursByVolunteer(ctx context.Context, volunteerID int64) (*VolunteeringHourLogsResponse, error) {
	res, err := s.volunteeringHourLogRepository.FindByVolunteerID(ctx, volunteerID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return &VolunteeringHourLogsResponse{
		VolunteeringHourLogs: res.VolunteeringHourLogs,
		Pages:                uint(res.TotalResults/dbctx.Get(ctx).Limit) + 1,
	}, nil
}
