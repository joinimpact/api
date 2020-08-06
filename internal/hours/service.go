package hours

import (
	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/pubsub"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/pkg/location"
	"github.com/rs/zerolog"
)

// Service represents a service for tracking volunteer hours.
type Service interface {
}

// service is the internal implementation of the hours.Service interface.
type service struct {
	volunteeringHourLogRepository        models.VolunteeringHourLogRepository
	volunteeringHourLogRequestRepository models.VolunteeringHourLogRequestRepository
	opportunityRepository                models.OpportunityRepository
	organizationRepository               models.OrganizationRepository
	userRepository                       models.UserRepository
	config                               *config.Config
	logger                               *zerolog.Logger
	snowflakeService                     snowflakes.SnowflakeService
	emailService                         email.Service
	broker                               pubsub.Broker
	locationService                      location.Service
	cdnClient                            *cdn.Client
}

// NewService creates and returns a new hours.Service.
func NewService(volunteeringHourLogRepository models.VolunteeringHourLogRepository, volunteeringHourLogRequestRepository models.VolunteeringHourLogRequestRepository, opportunityRepository models.OpportunityRepository, organizationRepository models.OrganizationRepository, userRepository models.UserRepository, userProfileFieldRepository models.UserProfileFieldRepository, userTagRepository models.UserTagRepository, tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service, broker pubsub.Broker, locationService location.Service) Service {
	return &service{
		volunteeringHourLogRepository,
		volunteeringHourLogRequestRepository,
		opportunityRepository,
		organizationRepository,
		userRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		broker,
		locationService,
		cdn.NewCDNClient(config),
	}
}
