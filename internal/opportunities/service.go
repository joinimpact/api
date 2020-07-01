package opportunities

import (
	"context"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
)

// Service represents a service for interacting with opportunities.
type Service interface {
	// GetOpportunity returns an opportunity by ID.
	GetOpportunity(ctx context.Context, id int64) (*models.Opportunity, error)
	// CreateOpportunity creates a new opportunity and returns the ID on success.
	CreateOpportunity(ctx context.Context, opportunity models.Opportunity) (int64, error)
	// UpdateOpportunity updates changed fields on an opportunity entity.
	UpdateOpportunity(ctx context.Context, opportunity models.Opportunity) error
}

// service represents the intenral implementation of the opportunities Service.
type service struct {
	opportunityRepository models.OpportunityRepository
	tagRepository         models.TagRepository
	config                *config.Config
	logger                *zerolog.Logger
	snowflakeService      snowflakes.SnowflakeService
	emailService          email.Service
	cdnClient             *cdn.Client
}

// NewService creates and returns a new Opportunities service with the provifded dependencies.
func NewService(opportunityRepository models.OpportunityRepository, tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service) Service {
	return &service{
		opportunityRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		cdn.NewCDNClient(config),
	}
}

// GetOpportunity returns an opportunity by ID.
func (s *service) GetOpportunity(ctx context.Context, id int64) (*models.Opportunity, error) {
	opportunity, err := s.opportunityRepository.FindByID(id)
	if err != nil {
		return nil, NewErrOpportunityNotFound()
	}

	return opportunity, nil
}

// CreateOpportunity creates a new opportunity and returns the ID on success.
func (s *service) CreateOpportunity(ctx context.Context, opportunity models.Opportunity) (int64, error) {
	// Generate a new ID for the opportunity.
	opportunity.ID = s.snowflakeService.GenerateID()
	err := s.opportunityRepository.Create(opportunity)
	if err != nil {
		return 0, NewErrServerError()
	}

	return opportunity.ID, nil
}

// UpdateOpportunity updates changed fields on an opportunity entity.
func (s *service) UpdateOpportunity(ctx context.Context, opportunity models.Opportunity) error {
	err := s.opportunityRepository.Update(opportunity)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}
