package opportunities

import (
	"context"
	"fmt"
	"io"
	"time"

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
	// DeleteOpportunity deletes a single opportunity by ID.
	DeleteOpportunity(ctx context.Context, id int64) error
	// UpdateOpportunity updates changed fields on an opportunity entity.
	UpdateOpportunity(ctx context.Context, opportunity OpportunityView) error
	// GetOpportunityTags gets all of a user's tags.
	GetOpportunityTags(ctx context.Context, opportunityID int64) ([]models.Tag, error)
	// AddOpportunityTags adds tags to a user by tag name.
	AddOpportunityTags(ctx context.Context, opportunityID int64, tags []string) (int, error)
	// RemoveOpportunityTag removes a tag from an opportunity by id.
	RemoveOpportunityTag(ctx context.Context, opportunityID, tagID int64) error
	// UploadProfilePicture uploads a profile picture to the CDN and adds it to the opportunity.
	UploadProfilePicture(opportunityID int64, fileReader io.Reader) (string, error)
	// RequestOpportunityMembership creates a membership request (as a volunteer) to join an opportunity.
	RequestOpportunityMembership(ctx context.Context, opportunityID int64, volunteerID int64) (int64, error)
	// GetOpportunityVolunteers returns an array of OpportunityMembership volunteer objects for a specified opportunity by ID.
	GetOpportunityVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembership, error)
	// GetOpportunityPendingVolunteers returns an array of OpportunityMembershipRequest objects for a specified opportunity by ID.
	GetOpportunityPendingVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembershipRequest, error)
	// PublishOpportunity attempts to publish an opportunity and returns an error if the opportunity is unpublishable.
	PublishOpportunity(ctx context.Context, opportunityID int64) error
	// UnpublishOpportunity unpublishes an opportunity.
	UnpublishOpportunity(ctx context.Context, opportunityID int64) error
}

// service represents the intenral implementation of the opportunities Service.
type service struct {
	opportunityRepository                  models.OpportunityRepository
	opportunityRequirementsRepository      models.OpportunityRequirementsRepository
	opportunityLimitsRepository            models.OpportunityLimitsRepository
	opportunityTagRepository               models.OpportunityTagRepository
	opportunityMembershipRepository        models.OpportunityMembershipRepository
	opportunityMembershipRequestRepository models.OpportunityMembershipRequestRepository
	opportunityMembershipInviteRepository  models.OpportunityMembershipInviteRepository
	tagRepository                          models.TagRepository
	config                                 *config.Config
	logger                                 *zerolog.Logger
	snowflakeService                       snowflakes.SnowflakeService
	emailService                           email.Service
	cdnClient                              *cdn.Client
}

// NewService creates and returns a new Opportunities service with the provifded dependencies.
func NewService(opportunityRepository models.OpportunityRepository, opportunityRequirementsRepository models.OpportunityRequirementsRepository, opportunityLimitsRepository models.OpportunityLimitsRepository, opportunityTagRepository models.OpportunityTagRepository, opportunityMembershipRepository models.OpportunityMembershipRepository, opportunityMembershipRequestRepository models.OpportunityMembershipRequestRepository, opportunityMembershipInviteRepository models.OpportunityMembershipInviteRepository, tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service) Service {
	return &service{
		opportunityRepository,
		opportunityRequirementsRepository,
		opportunityLimitsRepository,
		opportunityTagRepository,
		opportunityMembershipRepository,
		opportunityMembershipRequestRepository,
		opportunityMembershipInviteRepository,
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

		if !shouldAppear(view) {
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

	view.Tags, _ = s.GetOpportunityTags(ctx, opportunity.ID)

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

		err = s.opportunityLimitsRepository.Save(*limits)
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

		err = s.opportunityRequirementsRepository.Save(*requirements)
		if err != nil {
			return NewErrServerError()
		}
	}

	return nil
}

// DeleteOpportunity deletes a single opportunity by ID.
func (s *service) DeleteOpportunity(ctx context.Context, id int64) error {
	err := s.opportunityRepository.DeleteByID(id)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	return nil
}

// GetOpportunityTags gets all of an opportunity's tags.
func (s *service) GetOpportunityTags(ctx context.Context, opportunityID int64) ([]models.Tag, error) {
	_, err := s.opportunityRepository.FindByID(opportunityID)
	if err != nil {
		return nil, NewErrOrganizationNotFound()
	}

	// Find all OpportunityTag objects by organization ID.
	opportunityTags, err := s.opportunityTagRepository.FindByOpportunityID(opportunityID)
	if err != nil {
		return nil, NewErrServerError()
	}

	tags := []models.Tag{}
	for _, opportunityTag := range opportunityTags {
		// Get the tag by ID.
		tag, err := s.tagRepository.FindByID(opportunityTag.TagID)
		if err != nil {
			// Tag not found, skip.
			s.logger.Error().Err(err).Msg("Error in GetOpportunityTags: OpportunityTag object missing valid Tag")
			continue
		}

		// Append the tag to the tags array.
		tags = append(tags, *tag)
	}

	return tags, nil
}

// AddOpportunityTags adds tags to a user by tag name.
func (s *service) AddOpportunityTags(ctx context.Context, opportunityID int64, tags []string) (int, error) {
	// successfulTags counts how many tags were inserted correctly.
	successfulTags := 0

	_, err := s.opportunityRepository.FindByID(opportunityID)
	if err != nil {
		return successfulTags, NewErrOrganizationNotFound()
	}

	for _, tag := range tags {
		tag, err := s.tagRepository.FindByName(tag)
		if err != nil {
			// Log the error and skip the tag.
			s.logger.Error().Err(err).Msg("Error in AddOpportunityTags finding a tag")
			continue
		}

		// Increment the successful tags value as the tag was found.
		successfulTags++

		// Check to see if the organization already has this tag.
		_, err = s.opportunityTagRepository.FindOpportunityTagByID(opportunityID, tag.ID)
		if err == nil {
			// The organization already has this tag, skip.
			continue
		}

		// Create a new UserTag entity.
		id := s.snowflakeService.GenerateID()
		err = s.opportunityTagRepository.Create(models.OpportunityTag{
			Model: models.Model{
				ID: id,
			},
			OpportunityID: opportunityID,
			TagID:         tag.ID,
		})
		if err != nil {
			s.logger.Error().Err(err).Msg("Error in AddOpportunityTags creating a OpportunityTag")
			return successfulTags - 1, NewErrServerError()
		}
	}

	return successfulTags, nil
}

// RemoveOpportunityTag removes a tag from an organization by id.
func (s *service) RemoveOpportunityTag(ctx context.Context, opportunityID, tagID int64) error {
	opportunityTag, err := s.opportunityTagRepository.FindOpportunityTagByID(opportunityID, tagID)
	if err != nil {
		return NewErrTagNotFound()
	}

	return s.opportunityTagRepository.DeleteByID(opportunityTag.ID)
}

// UploadProfilePicture uploads a profile picture to the CDN and adds it to the opportunity.
func (s *service) UploadProfilePicture(opportunityID int64, fileReader io.Reader) (string, error) {
	url, err := s.cdnClient.UploadImage(fmt.Sprintf("opportunity-picture-%d-%d.png", opportunityID, time.Now().UTC().Unix()), fileReader)
	if err != nil {
		return "", err
	}

	return url, s.opportunityRepository.Update(models.Opportunity{
		Model: models.Model{
			ID: opportunityID,
		},
		ProfilePicture: url,
	})
}

// RequestOpportunityMembership creates a membership request (as a volunteer) to join an opportunity.
func (s *service) RequestOpportunityMembership(ctx context.Context, opportunityID int64, volunteerID int64) (int64, error) {
	_, err := s.opportunityMembershipRepository.FindUserInOpportunity(opportunityID, volunteerID)
	if err == nil {
		return 0, NewErrMembershipAlreadyRequested()
	}

	_, err = s.opportunityMembershipRequestRepository.FindInOpportunityByVolunteerID(opportunityID, volunteerID)
	if err == nil {
		return 0, NewErrMembershipAlreadyRequested()
	}

	// Create an ID for the request.
	id := s.snowflakeService.GenerateID()

	// Cretae the membership request entity.
	opportunityMembershipRequest := models.OpportunityMembershipRequest{
		Model: models.Model{
			ID: id,
		},
		Accepted:      false,
		VolunteerID:   volunteerID,
		OpportunityID: opportunityID,
	}

	// Attempt to create the entity.
	err = s.opportunityMembershipRequestRepository.Create(opportunityMembershipRequest)
	if err != nil {
		return 0, NewErrServerError()
	}

	return id, nil
}

// GetOpportunityVolunteers returns an array of OpportunityMembership volunteer objects for a specified opportunity by ID.
func (s *service) GetOpportunityVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembership, error) {
	// Get all memberships.
	memberships, err := s.opportunityMembershipRepository.FindByOpportunityID(opportunityID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return memberships, nil
}

// GetOpportunityPendingVolunteers returns an array of OpportunityMembershipRequest objects for a specified opportunity by ID.
func (s *service) GetOpportunityPendingVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembershipRequest, error) {
	// Get all membership requests by opportunity ID.
	requests, err := s.opportunityMembershipRequestRepository.FindByOpportunityID(opportunityID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return requests, nil
}

// PublishOpportunity attempts to publish an opportunity and returns an error if the opportunity is unpublishable.
func (s *service) PublishOpportunity(ctx context.Context, opportunityID int64) error {
	opportunity, err := s.opportunityRepository.FindByID(opportunityID)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	// Validate that the opportunity is publishable.
	if invalidFields, ok := isPublishable(*opportunity); !ok {
		return NewErrOpportunityNotPublishable(invalidFields)
	}

	opportunity.Public = true

	err = s.opportunityRepository.Save(*opportunity)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}

// UnpublishOpportunity unpublishes an opportunity.
func (s *service) UnpublishOpportunity(ctx context.Context, opportunityID int64) error {
	opportunity, err := s.opportunityRepository.FindByID(opportunityID)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	opportunity.Public = false

	err = s.opportunityRepository.Save(*opportunity)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}
