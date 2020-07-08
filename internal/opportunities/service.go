package opportunities

import (
	"context"
	"fmt"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
)

// Service represents a service for interacting with opportunities.
type Service interface {
	// GetOrganizationOpportunities gets all opportunities by organization ID.
	GetOrganizationOpportunities(ctx context.Context, organizationID int64) ([]OpportunityView, error)
	// GetOpportunity returns an opportunity by ID.
	GetOpportunity(ctx context.Context, id int64) (*OpportunityView, error)
	// CreateOpportunity creates a new opportunity and returns the ID on success.
	CreateOpportunity(ctx context.Context, opportunity OpportunityView) (int64, error)
	// UpdateOpportunity updates changed fields on an opportunity entity.
	UpdateOpportunity(ctx context.Context, opportunity OpportunityView) error
}

// service represents the intenral implementation of the opportunities Service.
type service struct {
	opportunityRepository             models.OpportunityRepository
	opportunityRequirementsRepository models.OpportunityRequirementsRepository
	opportunityLimitsRepository       models.OpportunityLimitsRepository
	tagRepository                     models.TagRepository
	config                            *config.Config
	logger                            *zerolog.Logger
	snowflakeService                  snowflakes.SnowflakeService
	emailService                      email.Service
	cdnClient                         *cdn.Client
}

// NewService creates and returns a new Opportunities service with the provifded dependencies.
func NewService(opportunityRepository models.OpportunityRepository, opportunityRequirementsRepository models.OpportunityRequirementsRepository, opportunityLimitsRepository models.OpportunityLimitsRepository, tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service) Service {
	return &service{
		opportunityRepository,
		opportunityRequirementsRepository,
		opportunityLimitsRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		cdn.NewCDNClient(config),
	}
}

// GetOrganizationOpportunities gets all opportunities by organization ID.
func (s *service) GetOrganizationOpportunities(ctx context.Context, organizationID int64) ([]OpportunityView, error) {
	views := []OpportunityView{}

	opportunities, err := s.opportunityRepository.FindByOrganizationID(organizationID)
	if err != nil {
		return views, NewErrServerError()
	}

	for _, opportunity := range opportunities {
		view, err := s.GetOpportunity(ctx, opportunity.ID)
		if err != nil {
			continue
		}

		views = append(views, *view)
	}

	return views, nil
}

// GetOpportunity returns an opportunity by ID.
func (s *service) GetOpportunity(ctx context.Context, id int64) (*OpportunityView, error) {
	view := &OpportunityView{}
	view.Requirements = &Requirements{}
	view.Limits = &Limits{}

	opportunity, err := s.opportunityRepository.FindByID(id)
	if err != nil {
		return nil, NewErrOpportunityNotFound()
	}

	view.ID = opportunity.ID
	view.OrganizationID = opportunity.OrganizationID
	view.CreatorID = opportunity.CreatorID
	view.ProfilePicture = opportunity.ProfilePicture
	view.Title = opportunity.Title
	view.Description = opportunity.Description
	view.Public = opportunity.Public

	opportunityRequirements, err := s.opportunityRequirementsRepository.FindByOpportunityID(opportunity.ID)
	if err == nil {
		if opportunityRequirements.AgeLimitActive {
			view.Requirements.AgeLimit = AgeLimit{
				Active: true,
				From:   opportunityRequirements.AgeLimitFrom,
				To:     opportunityRequirements.AgeLimitTo,
			}
		}
		if opportunityRequirements.ExpectedHoursActive {
			view.Requirements.ExpectedHours = ExpectedHours{
				Active: true,
				Hours:  opportunityRequirements.ExpectedHours,
			}
		}
	}

	opportunityLimits, err := s.opportunityLimitsRepository.FindByOpportunityID(opportunity.ID)
	if err == nil {
		if opportunityLimits.VolunteersCapActive {
			view.Limits.VolunteersCap = VolunteersCap{
				Active: true,
				Cap:    opportunityLimits.VolunteersCap,
			}
		}
	}

	return view, nil
}

// CreateOpportunity creates a new opportunity and returns the ID on success.
func (s *service) CreateOpportunity(ctx context.Context, view OpportunityView) (int64, error) {
	opportunity := models.Opportunity{}

	opportunity.Active = true
	opportunity.OrganizationID = view.OrganizationID
	opportunity.CreatorID = view.CreatorID
	opportunity.Title = view.Title
	opportunity.Description = view.Description
	opportunity.Public = false

	// Generate a new ID for the opportunity.
	opportunity.ID = s.snowflakeService.GenerateID()
	err := s.opportunityRepository.Create(opportunity)
	if err != nil {
		return 0, NewErrServerError()
	}

	limits := models.OpportunityLimits{}
	limits.ID = s.snowflakeService.GenerateID()
	limits.OpportunityID = opportunity.ID
	if view.Limits != nil {
		limits.VolunteersCapActive = view.Limits.VolunteersCap.Active
		limits.VolunteersCap = view.Limits.VolunteersCap.Cap
	}
	fmt.Println(limits)
	err = s.opportunityLimitsRepository.Create(limits)
	if err != nil {
		return 0, NewErrServerError()
	}

	requirements := models.OpportunityRequirements{}
	requirements.ID = s.snowflakeService.GenerateID()
	requirements.OpportunityID = opportunity.ID
	if view.Requirements != nil {
		requirements.AgeLimitActive = view.Requirements.AgeLimit.Active
		requirements.AgeLimitFrom = view.Requirements.AgeLimit.From
		requirements.AgeLimitTo = view.Requirements.AgeLimit.To

		requirements.ExpectedHoursActive = view.Requirements.ExpectedHours.Active
		requirements.ExpectedHours = view.Requirements.ExpectedHours.Hours
	}
	fmt.Println(requirements)
	err = s.opportunityRequirementsRepository.Create(requirements)
	if err != nil {
		return 0, NewErrServerError()
	}

	return opportunity.ID, nil
}

// UpdateOpportunity updates changed fields on an opportunity entity.
func (s *service) UpdateOpportunity(ctx context.Context, view OpportunityView) error {
	existingOpportunity, err := s.opportunityRepository.FindByID(view.ID)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	opportunity := models.Opportunity{}

	opportunity.ID = existingOpportunity.ID
	opportunity.ProfilePicture = view.ProfilePicture
	opportunity.Title = view.Title
	opportunity.Description = view.Description
	opportunity.Public = view.Public

	err = s.opportunityRepository.Update(opportunity)
	if err != nil {
		return NewErrServerError()
	}

	if view.Limits != nil {
		limits, err := s.opportunityLimitsRepository.FindByOpportunityID(opportunity.ID)
		if err != nil {
			return NewErrServerError()
		}

		limits.VolunteersCapActive = view.Limits.VolunteersCap.Active
		limits.VolunteersCap = view.Limits.VolunteersCap.Cap

		err = s.opportunityLimitsRepository.Update(*limits)
		if err != nil {
			return NewErrServerError()
		}
	}

	if view.Requirements != nil {
		requirements, err := s.opportunityRequirementsRepository.FindByOpportunityID(opportunity.ID)
		if err != nil {
			return NewErrServerError()
		}

		requirements.AgeLimitActive = view.Requirements.AgeLimit.Active
		requirements.AgeLimitFrom = view.Requirements.AgeLimit.From
		requirements.AgeLimitTo = view.Requirements.AgeLimit.To

		requirements.ExpectedHoursActive = view.Requirements.ExpectedHours.Active
		requirements.ExpectedHours = view.Requirements.ExpectedHours.Hours

		err = s.opportunityRequirementsRepository.Update(*requirements)
		if err != nil {
			return NewErrServerError()
		}
	}

	return nil
}
