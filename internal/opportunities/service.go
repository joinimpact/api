package opportunities

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/email/templates"
	"github.com/joinimpact/api/internal/models"
	opportunitiesSearch "github.com/joinimpact/api/internal/search/stores/opportunities"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/pkg/location"
	"github.com/rs/zerolog"
)

// Service represents a service for interacting with opportunities.
type Service interface {
	// GetOrganizationOpportunities gets all opportunities by organization ID.
	GetOrganizationOpportunities(ctx context.Context, organizationID int64) ([]OpportunityView, error)
	// GetVolunteerOpportunities gets all opportunities by user ID where
	// user is a registered volunteer.
	GetVolunteerOpportunities(ctx context.Context, userID int64) ([]OpportunityView, error)
	// GetOpportunity returns an opportunity by ID.
	GetOpportunity(ctx context.Context, id int64) (*OpportunityView, error)
	// GetMinimalOpportunity returns an opportunity without tags or a profile.
	GetMinimalOpportunity(ctx context.Context, id int64) (*OpportunityView, error)
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
	UploadProfilePicture(ctx context.Context, opportunityID int64, fileReader io.Reader) (string, error)
	// CanRequestOpportunityMembership checks if a user can request membership or not.
	CanRequestOpportunityMembership(ctx context.Context, opportunityID, volunteerID int64) error
	// AcceptOpportunityMembershipRequest accepts a membership request from a volunteer by user ID.
	AcceptOpportunityMembershipRequest(ctx context.Context, opportunityID, volunteerID, userID int64) error
	// RequestOpportunityMembership creates a membership request (as a volunteer) to join an opportunity.
	RequestOpportunityMembership(ctx context.Context, opportunityID int64, volunteerID int64) (int64, error)
	// GetOpportunityVolunteers returns an array of OpportunityMembership volunteer objects for a specified opportunity by ID.
	GetOpportunityVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembership, error)
	// GetOpportunityPendingVolunteers returns an array of OpportunityMembershipRequest objects for a specified opportunity by ID.
	GetOpportunityPendingVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembershipRequest, error)
	// GetOpportunityInvitedVolunteers returns an array of OpportunityMembershipInvite objects for a specified opportunity by ID.
	GetOpportunityInvitedVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembershipInvite, error)
	// PublishOpportunity attempts to publish an opportunity and returns an error if the opportunity is unpublishable.
	PublishOpportunity(ctx context.Context, opportunityID int64) error
	// UnpublishOpportunity unpublishes an opportunity.
	UnpublishOpportunity(ctx context.Context, opportunityID int64) error
	// InviteVolunteer invites a volunteer by user email to an opportunity.
	InviteVolunteer(ctx context.Context, inviterID, opportunityID int64, userEmail string) error
	// GetOpportunityFromInvite gets an opportunity view from an invite for UI use.
	GetOpportunityFromInvite(ctx context.Context, opportunityID int64, userID, inviteID int64, inviteKey string) (*OpportunityView, error)
	// AcceptInvite accepts an invite.
	AcceptInvite(ctx context.Context, opportunityID int64, userID, inviteID int64, inviteKey string) error
	// DeclineInvite declines an invite.
	DeclineInvite(ctx context.Context, opportunityID int64, userID, inviteID int64, inviteKey string) error
	// GetOpportunityMembership returns the permissions level of a single user's relationship with an opportunity.
	GetOpportunityMembership(ctx context.Context, opportunityID, userID int64) (int, error)
	// Search searches opportunities by a query struct.
	Search(ctx context.Context, query opportunitiesSearch.Query) ([]OpportunityView, error)
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
	userRepository                         models.UserRepository
	organizationRepository                 models.OrganizationRepository
	config                                 *config.Config
	logger                                 *zerolog.Logger
	snowflakeService                       snowflakes.SnowflakeService
	emailService                           email.Service
	cdnClient                              *cdn.Client
	searchStore                            opportunitiesSearch.Store
	locationService                        location.Service
}

// NewService creates and returns a new Opportunities service with the provifded dependencies.
func NewService(opportunityRepository models.OpportunityRepository, opportunityRequirementsRepository models.OpportunityRequirementsRepository, opportunityLimitsRepository models.OpportunityLimitsRepository, opportunityTagRepository models.OpportunityTagRepository, opportunityMembershipRepository models.OpportunityMembershipRepository, opportunityMembershipRequestRepository models.OpportunityMembershipRequestRepository, opportunityMembershipInviteRepository models.OpportunityMembershipInviteRepository, tagRepository models.TagRepository, userRepository models.UserRepository, organizationRepository models.OrganizationRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service, searchStore opportunitiesSearch.Store, locationService location.Service) Service {
	return &service{
		opportunityRepository,
		opportunityRequirementsRepository,
		opportunityLimitsRepository,
		opportunityTagRepository,
		opportunityMembershipRepository,
		opportunityMembershipRequestRepository,
		opportunityMembershipInviteRepository,
		tagRepository,
		userRepository,
		organizationRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		cdn.NewCDNClient(config),
		searchStore,
		locationService,
	}
}

// GetOrganizationOpportunities gets all opportunities by organization ID.
func (s *service) GetOrganizationOpportunities(ctx context.Context, organizationID int64) ([]OpportunityView, error) {
	views := []OpportunityView{}

	opportunities, err := s.opportunityRepository.FindByOrganizationID(ctx, organizationID)
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

// GetVolunteerOpportunities gets all opportunities by user ID where
// user is a registered volunteer.
func (s *service) GetVolunteerOpportunities(ctx context.Context, userID int64) ([]OpportunityView, error) {
	views := []OpportunityView{}

	memberships, err := s.opportunityMembershipRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	for _, membership := range memberships {
		view, err := s.GetOpportunity(ctx, membership.OpportunityID)
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

	opportunity, err := s.opportunityRepository.FindByID(ctx, id)
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

	_, view.Publishable = isPublishable(*opportunity)

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

	organization, err := s.organizationRepository.FindByID(opportunity.OrganizationID)
	if err == nil {
		view.OpportunityOrganizationProfile = &OpportunityOrganizationProfile{
			ID:             organization.ID,
			Name:           organization.Name,
			ProfilePicture: organization.ProfilePicture,
		}

		// Location
		if organization.LocationLatitude != 0.0 || organization.LocationLongitude != 0.0 {
			coordinates := &location.Coordinates{
				Latitude:  organization.LocationLatitude,
				Longitude: organization.LocationLongitude,
			}

			location, err := s.locationService.CoordinatesToCity(coordinates)
			if err == nil {
				view.Location = location
			}
		}
	}

	view.Tags, _ = s.GetOpportunityTags(ctx, opportunity.ID)

	return view, nil
}

// GetMinimalOpportunity returns an opportunity without tags or a profile.
func (s *service) GetMinimalOpportunity(ctx context.Context, id int64) (*OpportunityView, error) {
	view := &OpportunityView{}
	view.Requirements = &Requirements{}
	view.Limits = &Limits{}

	opportunity, err := s.opportunityRepository.FindByID(ctx, id)
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
	err := s.opportunityRepository.Create(ctx, opportunity)
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

	err = s.opportunityRequirementsRepository.Create(requirements)
	if err != nil {
		return 0, NewErrServerError()
	}

	s.searchStore.Save(opportunity.ID)

	return opportunity.ID, nil
}

// UpdateOpportunity updates changed fields on an opportunity entity.
func (s *service) UpdateOpportunity(ctx context.Context, view OpportunityView) error {
	existingOpportunity, err := s.opportunityRepository.FindByID(ctx, view.ID)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	opportunity := models.Opportunity{}

	opportunity.ID = existingOpportunity.ID
	opportunity.ProfilePicture = view.ProfilePicture
	opportunity.Title = view.Title
	opportunity.Description = view.Description
	opportunity.Public = view.Public

	err = s.opportunityRepository.Update(ctx, opportunity)
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

	s.searchStore.Save(opportunity.ID)

	return nil
}

// DeleteOpportunity deletes a single opportunity by ID.
func (s *service) DeleteOpportunity(ctx context.Context, id int64) error {
	err := s.opportunityRepository.DeleteByID(ctx, id)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	return nil
}

// GetOpportunityTags gets all of an opportunity's tags.
func (s *service) GetOpportunityTags(ctx context.Context, opportunityID int64) ([]models.Tag, error) {
	_, err := s.opportunityRepository.FindByID(ctx, opportunityID)
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

	_, err := s.opportunityRepository.FindByID(ctx, opportunityID)
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

	s.searchStore.Save(opportunityID)

	return successfulTags, nil
}

// RemoveOpportunityTag removes a tag from an organization by id.
func (s *service) RemoveOpportunityTag(ctx context.Context, opportunityID, tagID int64) error {
	opportunityTag, err := s.opportunityTagRepository.FindOpportunityTagByID(opportunityID, tagID)
	if err != nil {
		return NewErrTagNotFound()
	}

	s.searchStore.Save(opportunityID)

	return s.opportunityTagRepository.DeleteByID(opportunityTag.ID)
}

// UploadProfilePicture uploads a profile picture to the CDN and adds it to the opportunity.
func (s *service) UploadProfilePicture(ctx context.Context, opportunityID int64, fileReader io.Reader) (string, error) {
	url, err := s.cdnClient.UploadImage(fmt.Sprintf("opportunity-picture-%d-%d.png", opportunityID, time.Now().UTC().Unix()), fileReader)
	if err != nil {
		return "", err
	}

	return url, s.opportunityRepository.Update(ctx, models.Opportunity{
		Model: models.Model{
			ID: opportunityID,
		},
		ProfilePicture: url,
	})
}

// RequestOpportunityMembership creates a membership request (as a volunteer) to join an opportunity.
func (s *service) RequestOpportunityMembership(ctx context.Context, opportunityID int64, volunteerID int64) (int64, error) {
	if err := s.CanRequestOpportunityMembership(ctx, opportunityID, volunteerID); err != nil {
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
	err := s.opportunityMembershipRequestRepository.Create(opportunityMembershipRequest)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error creating opportunity membership request")
		return 0, NewErrServerError()
	}

	return id, nil
}

// CanRequestOpportunityMembership checks if a user can request membership or not.
func (s *service) CanRequestOpportunityMembership(ctx context.Context, opportunityID, volunteerID int64) error {
	_, err := s.opportunityMembershipRepository.FindUserInOpportunity(ctx, opportunityID, volunteerID)
	if err == nil {
		return NewErrMembershipAlreadyRequested()
	}

	_, err = s.opportunityMembershipRequestRepository.FindInOpportunityByVolunteerID(opportunityID, volunteerID)
	if err == nil {
		return NewErrMembershipAlreadyRequested()
	}

	return nil
}

// AcceptOpportunityMembershipRequest accepts a membership request from a volunteer by user ID.
func (s *service) AcceptOpportunityMembershipRequest(ctx context.Context, opportunityID, volunteerID, userID int64) error {
	// Check that a valid request exists.
	membershipRequest, err := s.opportunityMembershipRequestRepository.FindInOpportunityByVolunteerID(opportunityID, volunteerID)
	if err != nil {
		return NewErrRequestNotFound()
	}

	// Create the volunteer membership.
	if err := s.createVolunteerMembership(ctx, userID, membershipRequest.OpportunityID, membershipRequest.VolunteerID); err != nil {
		s.logger.Error().Err(err).Msg("Error creating volunteer membership")
		return NewErrServerError()
	}

	// Delete the membership request.
	if err := s.opportunityMembershipRequestRepository.DeleteByID(membershipRequest.ID); err != nil {
		s.logger.Error().Err(err).Msg("Error deleting membership request")
		return NewErrServerError()
	}

	return nil
}

// GetOpportunityVolunteers returns an array of OpportunityMembership volunteer objects for a specified opportunity by ID.
func (s *service) GetOpportunityVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembership, error) {
	// Get all memberships.
	memberships, err := s.opportunityMembershipRepository.FindByOpportunityID(ctx, opportunityID)
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

// GetOpportunityInvitedVolunteers returns an array of OpportunityMembershipRequest objects for a specified opportunity by ID.
func (s *service) GetOpportunityInvitedVolunteers(ctx context.Context, opportunityID int64) ([]models.OpportunityMembershipInvite, error) {
	// Get all membership invites by opportunity ID.
	invites, err := s.opportunityMembershipInviteRepository.FindByOpportunityID(opportunityID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return invites, nil
}

// PublishOpportunity attempts to publish an opportunity and returns an error if the opportunity is unpublishable.
func (s *service) PublishOpportunity(ctx context.Context, opportunityID int64) error {
	opportunity, err := s.opportunityRepository.FindByID(ctx, opportunityID)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	// Validate that the opportunity is publishable.
	if invalidFields, ok := isPublishable(*opportunity); !ok {
		return NewErrOpportunityNotPublishable(invalidFields)
	}

	opportunity.Public = true

	err = s.opportunityRepository.Save(ctx, *opportunity)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error publishing opportunity")
		return NewErrServerError()
	}

	s.searchStore.Save(opportunityID)

	return nil
}

// UnpublishOpportunity unpublishes an opportunity.
func (s *service) UnpublishOpportunity(ctx context.Context, opportunityID int64) error {
	opportunity, err := s.opportunityRepository.FindByID(ctx, opportunityID)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	opportunity.Public = false

	err = s.opportunityRepository.Save(ctx, *opportunity)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error unpublishing opportunity")
		return NewErrServerError()
	}

	s.searchStore.Save(opportunityID)

	return nil
}

// InviteVolunteer invites a volunteer by user email to an opportunity.
func (s *service) InviteVolunteer(ctx context.Context, inviterID, opportunityID int64, userEmail string) error {
	_, err := s.opportunityMembershipInviteRepository.FindInOpportunityByEmail(opportunityID, userEmail)
	if err == nil {
		return NewErrUserAlreadyInvited()
	}

	opportunity, err := s.opportunityRepository.FindByID(ctx, opportunityID)
	if err != nil {
		return NewErrOpportunityNotFound()
	}

	organization, err := s.organizationRepository.FindByID(opportunity.OrganizationID)
	if err != nil {
		return NewErrServerError()
	}

	invite := models.OpportunityMembershipInvite{}

	// Generate an ID for the invite.
	invite.ID = s.snowflakeService.GenerateID()
	invite.InviteeEmail = userEmail

	userFirstName := "Impact"
	userLastName := "User"

	user, err := s.userRepository.FindByEmail(userEmail)
	if err == nil {
		// If a user was found, add their ID to the invite.
		invite.InviteeID = user.ID
		userFirstName = user.FirstName
		userLastName = user.LastName
	}

	invite.InviterID = inviterID
	invite.OpportunityID = opportunityID

	// Generate a key.
	invite.Key = generateKey()

	err = s.opportunityMembershipInviteRepository.Create(invite)
	if err != nil {
		return NewErrServerError()
	}

	salutationName := "friend"

	if userFirstName != "Impact" {
		salutationName = userFirstName
	}

	// Create a new email with the reset password template.
	email := s.emailService.NewEmail(
		email.NewRecipient(fmt.Sprintf("%s %s", userFirstName, userLastName), userEmail),
		fmt.Sprintf("You've been invited to join %s on Impact!", opportunity.Title),
		templates.OpportunityInvitationTemplate(salutationName, opportunity.Title, organization.Name, opportunity.ID, invite.ID, invite.Key),
	)
	err = s.emailService.Send(email)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error sending opportunity invite email")
		return NewErrServerError()
	}

	return nil
}

// GetOpportunityFromInvite gets an opportunity view from an invite for UI use.
func (s *service) GetOpportunityFromInvite(ctx context.Context, opportunityID int64, userID, inviteID int64, inviteKey string) (*OpportunityView, error) {
	// Get the invite by ID.
	invite, err := s.opportunityMembershipInviteRepository.FindByID(inviteID)
	if err != nil {
		return nil, NewErrInviteInvalid()
	}

	// Get the user by ID to check the email.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, NewErrInviteInvalid()
	}

	// Check if the invite is valid.
	if invite.InviteeEmail != user.Email || invite.Key != inviteKey || invite.OpportunityID != opportunityID {
		return nil, NewErrInviteInvalid()
	}

	return s.GetOpportunity(ctx, invite.OpportunityID)
}

// AcceptInvite accepts an invite.
func (s *service) AcceptInvite(ctx context.Context, opportunityID int64, userID, inviteID int64, inviteKey string) error {
	// Get the invite by ID.
	invite, err := s.opportunityMembershipInviteRepository.FindByID(inviteID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Get the user by ID to check the email.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Check if the invite is valid.
	if invite.InviteeEmail != user.Email || invite.Key != inviteKey || invite.OpportunityID != opportunityID {
		return NewErrInviteInvalid()
	}

	if err := s.createVolunteerMembership(ctx, invite.InviterID, invite.OpportunityID, userID); err != nil {
		s.logger.Error().Err(err).Msg("Error creating volunteer membership")
		return NewErrServerError()
	}

	if err := s.invalidateInvite(ctx, invite.ID); err != nil {
		s.logger.Error().Err(err).Msg("Error invalidating invite")
		return NewErrServerError()
	}

	return nil
}

// DeclineInvite declines an invite.
func (s *service) DeclineInvite(ctx context.Context, opportunityID int64, userID, inviteID int64, inviteKey string) error {
	// Get the invite by ID.
	invite, err := s.opportunityMembershipInviteRepository.FindByID(inviteID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Get the user by ID to check the email.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Check if the invite is valid.
	if invite.InviteeEmail != user.Email || invite.Key != inviteKey || invite.OpportunityID != opportunityID {
		return NewErrInviteInvalid()
	}

	if err := s.invalidateInvite(ctx, invite.ID); err != nil {
		return NewErrServerError()
	}

	return nil
}

// createVolunteerMembership creates a volunteer membership in an opportunity.
func (s *service) createVolunteerMembership(ctx context.Context, inviterID int64, opportunityID, userID int64) error {
	membership := models.OpportunityMembership{}

	membership.Active = true
	membership.ID = s.snowflakeService.GenerateID()
	membership.UserID = userID
	membership.JoinedAt = time.Now()
	membership.OpportunityID = opportunityID
	membership.PermissionsFlag = models.OpportunityPermissionsMember
	membership.InviterID = inviterID

	return s.opportunityMembershipRepository.Create(ctx, membership)
}

// invalidateInvite invalidates an opportunity invite by ID.
func (s *service) invalidateInvite(ctx context.Context, inviteID int64) error {
	return s.opportunityMembershipInviteRepository.DeleteByID(inviteID)
}

// GetOpportunityMembership returns the permissions level of a single user's relationship with an opportunity.
func (s *service) GetOpportunityMembership(ctx context.Context, opportunityID, userID int64) (int, error) {
	membership, err := s.opportunityMembershipRepository.FindUserInOpportunity(ctx, opportunityID, userID)
	if err != nil {
		return 0, err
	}

	return membership.PermissionsFlag, nil
}

// Search searches opportunities by a query struct.
func (s *service) Search(ctx context.Context, query opportunitiesSearch.Query) ([]OpportunityView, error) {
	views := []OpportunityView{}

	search, err := s.searchStore.Search(query)
	if err != nil {
		s.logger.Error().Err(err).Msg("Error searching opportunities")
		return nil, NewErrServerError()
	}

	for _, item := range search {
		view, err := s.GetOpportunity(ctx, item.ID)
		if err != nil {
			continue
		}

		views = append(views, *view)
	}

	return views, nil
}
