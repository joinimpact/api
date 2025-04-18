package organizations

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
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/pkg/location"
	"github.com/rs/zerolog"
)

// Service represents a provider of Organization services.
type Service interface {
	// GetUserOrganizations gets all of a user's organizations as
	// organization profile.
	GetUserOrganizations(userID int64) ([]OrganizationProfile, error)
	// GetOrganizationProfile retrieves a single organization's profile.
	GetOrganizationProfile(organizationID int64) (*OrganizationProfile, error)
	// UpdateOrganizationProfile updates a user's profile.
	UpdateOrganizationProfile(organizationID int64, profile OrganizationProfile) error
	// UpdateOrganizationLocation updates an organization's location.
	UpdateOrganizationLocation(organizationID int64, location *location.Coordinates) error
	// SetOrganizationProfileField sets an organization's profile field by name.
	SetOrganizationProfileField(organizationID int64, profileField models.OrganizationProfileField) error
	// CreateOrganization creates a new organization and returns the ID on success.
	CreateOrganization(organization models.Organization) (int64, error)
	// DeleteOrganization deletes a single organization by ID.
	DeleteOrganization(id int64) error
	// GetOrganizationTags gets all of a user's tags.
	GetOrganizationTags(organizationID int64) ([]models.Tag, error)
	// AddOrganizationTags adds tags to a user by tag name.
	AddOrganizationTags(organizationID int64, tags []string) (int, error)
	// RemoveOrganizationTag removes a tag from an organization by id.
	RemoveOrganizationTag(organizationID, tagID int64) error
	// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
	UploadProfilePicture(organizationID int64, fileReader io.Reader) (string, error)
	// InviteUser invites a user by user email to an organization.
	InviteUser(inviterID, organizationID int64, userEmail string, permissionsFlag int) error
	// GetOrganizationInvitedVolunteers returns an array of OrganizationMembershipRequest objects for a specified organization by ID.
	GetOrganizationInvitedVolunteers(ctx context.Context, organizationID int64) ([]models.OrganizationMembershipInvite, error)
	// GetOrganizationMemberships gets all users in an organization.
	GetOrganizationMemberships(organizationID int64) ([]models.OrganizationMembership, error)
	// GetOrganizationMembership returns the membership level of a user in an organization.
	// Returns an error if no membership is found.
	GetOrganizationMembership(organizationID, userID int64) (int, error)
	// GetOrganizationFromInvite gets an organization profile from an invite for UI use.
	GetOrganizationFromInvite(ctx context.Context, organizationID int64, userID, inviteID int64, inviteKey string) (*OrganizationProfile, error)
	// AcceptInvite accepts an invite.
	AcceptInvite(ctx context.Context, organizationID int64, userID, inviteID int64, inviteKey string) error
	// DeclineInvite declines an invite.
	DeclineInvite(ctx context.Context, organizationID int64, userID, inviteID int64, inviteKey string) error
}

// service represents the internal implementation of the organizations Service.
type service struct {
	organizationRepository                 models.OrganizationRepository
	organizationMembershipRepository       models.OrganizationMembershipRepository
	organizationMembershipInviteRepository models.OrganizationMembershipInviteRepository
	organizationProfileFieldRepository     models.OrganizationProfileFieldRepository
	organizationTagRepository              models.OrganizationTagRepository
	userRepository                         models.UserRepository
	tagRepository                          models.TagRepository
	config                                 *config.Config
	logger                                 *zerolog.Logger
	snowflakeService                       snowflakes.SnowflakeService
	emailService                           email.Service
	cdnClient                              *cdn.Client
	locationService                        location.Service
}

// NewService creates and returns a new Users service with the provifded dependencies.
func NewService(organizationRepository models.OrganizationRepository, organizationMembershipRepository models.OrganizationMembershipRepository, organizationMembershipInviteRepository models.OrganizationMembershipInviteRepository, organizationProfileFieldRepository models.OrganizationProfileFieldRepository, organizationTagRepository models.OrganizationTagRepository,
	userRepository models.UserRepository, tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, emailService email.Service, locationService location.Service) Service {
	return &service{
		organizationRepository,
		organizationMembershipRepository,
		organizationMembershipInviteRepository,
		organizationProfileFieldRepository,
		organizationTagRepository,
		userRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		emailService,
		cdn.NewCDNClient(config),
		locationService,
	}
}

// GetUserOrganizations gets all of a user's organizations as
// organization profile.
func (s *service) GetUserOrganizations(userID int64) ([]OrganizationProfile, error) {
	memberships, err := s.organizationMembershipRepository.FindByUserID(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	organizations := []OrganizationProfile{}
	for _, membership := range memberships {
		organization, err := s.getMinimumOrganizationProfile(membership.OrganizationID)
		if err != nil {
			// If there is an error, consider the organization to be
			// unreachable and skip it.
			continue
		}

		organizations = append(organizations, *organization)
	}

	return organizations, nil
}

// getMinimumOrganizationProfile gets an organization profile by ID without
// extra profile fields, tags, or location.
func (s *service) getMinimumOrganizationProfile(organizationID int64) (*OrganizationProfile, error) {
	// Find the organization to verify that it is active.
	organization, err := s.organizationRepository.FindByID(organizationID)
	if err != nil {
		return nil, NewErrOrganizationNotFound()
	}

	return s.organizationToProfile(organization), nil
}

// organizationToProfile gets an organization without extra profile fields,
// tags, or location.
func (s *service) organizationToProfile(organization *models.Organization) *OrganizationProfile {
	profile := &OrganizationProfile{}

	profile.ID = organization.ID
	profile.CreatorID = organization.CreatorID
	profile.Name = organization.Name
	profile.Description = organization.Description
	profile.ProfilePicture = organization.ProfilePicture
	profile.WebsiteURL = organization.WebsiteURL

	return profile
}

// GetOrganizationProfile retrieves a single organization's profile.
func (s *service) GetOrganizationProfile(organizationID int64) (*OrganizationProfile, error) {
	// Find the organization to verify that it is active.
	organization, err := s.organizationRepository.FindByID(organizationID)
	if err != nil {
		return nil, NewErrOrganizationNotFound()
	}

	profile := s.organizationToProfile(organization)

	tags, err := s.GetOrganizationTags(organizationID)
	if err != nil {
		return nil, NewErrServerError()
	}

	profile.Tags = tags

	// Location
	if organization.LocationLatitude != 0.0 || organization.LocationLongitude != 0.0 {
		coordinates := &location.Coordinates{
			Latitude:  organization.LocationLatitude,
			Longitude: organization.LocationLongitude,
		}

		location, err := s.locationService.CoordinatesToCity(coordinates)
		if err == nil {
			profile.Location = location
		}
	}

	// Profile fields
	fields, err := s.organizationProfileFieldRepository.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, NewErrServerError()
	}

	profile.ProfileFields = []models.OrganizationProfileField{}
	profile.ProfileFields = fields

	return profile, nil
}

// UpdateOrganizationProfile updates a user's profile.
func (s *service) UpdateOrganizationProfile(organizationID int64, profile OrganizationProfile) error {
	return s.organizationRepository.Update(models.Organization{
		Model: models.Model{
			ID: organizationID,
		},
		Name:        profile.Name,
		Description: profile.Description,
		WebsiteURL:  formatURL(profile.WebsiteURL),
	})
}

// UpdateOrganizationLocation updates an organization's location.
func (s *service) UpdateOrganizationLocation(organizationID int64, location *location.Coordinates) error {
	return s.organizationRepository.Update(models.Organization{
		Model: models.Model{
			ID: organizationID,
		},
		LocationLatitude:  location.Latitude,
		LocationLongitude: location.Longitude,
	})
}

// SetOrganizationProfileField sets an organization's profile field by name.
func (s *service) SetOrganizationProfileField(organizationID int64, profileField models.OrganizationProfileField) error {
	field, err := s.organizationProfileFieldRepository.FindOrganizationFieldByName(organizationID, profileField.Name)
	if err == nil {
		profileField.ID = field.ID
		return s.organizationProfileFieldRepository.Update(profileField)
	}

	profileField.OrganizationID = organizationID

	// Create an ID and assign it to the profile field.
	id := s.snowflakeService.GenerateID()
	profileField.ID = id

	// Create the profile field.
	return s.organizationProfileFieldRepository.Create(profileField)
}

// CreateOrganization creates a new organization and returns the ID on success.
func (s *service) CreateOrganization(organization models.Organization) (int64, error) {
	if len(organization.Name) < 1 {
		return 0, NewErrServerError()
	}

	// Generate an ID for the new organization.
	organization.ID = s.snowflakeService.GenerateID()

	// Make the organization active if not previously true.
	organization.Active = true

	// Create the organization.
	err := s.organizationRepository.Create(organization)
	if err != nil {
		return 0, NewErrServerError()
	}

	// Create the organization membership.
	organizationMembership := models.OrganizationMembership{
		Active:          true,
		UserID:          organization.CreatorID,
		OrganizationID:  organization.ID,
		PermissionsFlag: models.OrganizationPermissionsCreator,
	}

	// Generate an ID for the membership.
	organizationMembership.ID = s.snowflakeService.GenerateID()

	// Create the membership and add it to the repository.
	err = s.organizationMembershipRepository.Create(organizationMembership)
	if err != nil {
		return 0, NewErrServerError()
	}

	return organization.ID, nil
}

// DeleteOrganization deletes a single organization by ID.
func (s *service) DeleteOrganization(id int64) error {
	err := s.organizationRepository.DeleteByID(id)
	if err != nil {
		return NewErrOrganizationNotFound()
	}

	return nil
}

// GetOrganizationTags gets all of a user's tags.
func (s *service) GetOrganizationTags(organizationID int64) ([]models.Tag, error) {
	_, err := s.organizationRepository.FindByID(organizationID)
	if err != nil {
		return nil, NewErrOrganizationNotFound()
	}

	// Find all OrganizationTag objects by organization ID.
	organizationTags, err := s.organizationTagRepository.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, NewErrServerError()
	}

	tags := []models.Tag{}
	for _, organizationTag := range organizationTags {
		// Get the tag by ID.
		tag, err := s.tagRepository.FindByID(organizationTag.TagID)
		if err != nil {
			// Tag not found, skip.
			s.logger.Error().Err(err).Msg("Error in GetOrganizationTags: OrganizationTag object missing valid Tag")
			continue
		}

		// Append the tag to the tags array.
		tags = append(tags, *tag)
	}

	return tags, nil
}

// AddOrganizationTags adds tags to a user by tag name.
func (s *service) AddOrganizationTags(organizationID int64, tags []string) (int, error) {
	// successfulTags counts how many tags were inserted correctly.
	successfulTags := 0

	_, err := s.organizationRepository.FindByID(organizationID)
	if err != nil {
		return successfulTags, NewErrOrganizationNotFound()
	}

	for _, tag := range tags {
		tag, err := s.tagRepository.FindByName(tag)
		if err != nil {
			// Log the error and skip the tag.
			s.logger.Error().Err(err).Msg("Error in AddOrganizationTags finding a tag")
			continue
		}

		// Increment the successful tags value as the tag was found.
		successfulTags++

		// Check to see if the organization already has this tag.
		_, err = s.organizationTagRepository.FindOrganizationTagByID(organizationID, tag.ID)
		if err == nil {
			// The organization already has this tag, skip.
			continue
		}

		// Create a new UserTag entity.
		id := s.snowflakeService.GenerateID()
		err = s.organizationTagRepository.Create(models.OrganizationTag{
			Model: models.Model{
				ID: id,
			},
			OrganizationID: organizationID,
			TagID:          tag.ID,
		})
		if err != nil {
			s.logger.Error().Err(err).Msg("Error in AddOrganizationTags creating a OrganizationTag")
			return successfulTags - 1, NewErrServerError()
		}
	}

	return successfulTags, nil
}

// RemoveOrganizationTag removes a tag from an organization by id.
func (s *service) RemoveOrganizationTag(organizationID, tagID int64) error {
	organizationTag, err := s.organizationTagRepository.FindOrganizationTagByID(organizationID, tagID)
	if err != nil {
		return NewErrTagNotFound()
	}

	return s.organizationTagRepository.DeleteByID(organizationTag.ID)
}

// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
func (s *service) UploadProfilePicture(organizationID int64, fileReader io.Reader) (string, error) {
	url, err := s.cdnClient.UploadImage(fmt.Sprintf("organization-picture-%d-%d.png", organizationID, time.Now().UTC().Unix()), fileReader)
	if err != nil {
		return "", err
	}

	return url, s.organizationRepository.Update(models.Organization{
		Model: models.Model{
			ID: organizationID,
		},
		ProfilePicture: url,
	})
}

// InviteUser invites a user by user email to an organization.
func (s *service) InviteUser(inviterID, organizationID int64, userEmail string, permissionsFlag int) error {
	organization, err := s.organizationRepository.FindByID(organizationID)
	if err != nil {
		return NewErrOrganizationNotFound()
	}

	user, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return s.inviteByEmail(inviterID, organization, userEmail, permissionsFlag)
	}

	return s.inviteByID(inviterID, organization, user, permissionsFlag)
}

func (s *service) inviteByEmail(inviterID int64, organization *models.Organization, userEmail string, permissionsFlag int) error {
	// Generate an ID for the invite.
	id := s.snowflakeService.GenerateID()
	key := generateKey()
	err := s.organizationMembershipInviteRepository.Create(models.OrganizationMembershipInvite{
		Model: models.Model{
			ID: id,
		},
		Accepted:       false,
		InviteeEmail:   userEmail,
		OrganizationID: organization.ID,
		InviterID:      inviterID,
		Key:            key,
	})
	if err != nil {
		return NewErrServerError()
	}

	// Create a new email with the reset password template.
	email := s.emailService.NewEmail(
		email.NewRecipient(fmt.Sprintf("%s %s", "Impact", "User"), userEmail),
		fmt.Sprintf("You've been invited to join %s!", organization.Name),
		templates.OrganizationInvitationTemplate("friend", organization.Name, organization.ID, id, key),
	)
	err = s.emailService.Send(email)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}

func (s *service) inviteByID(inviterID int64, organization *models.Organization, user *models.User, permissionsFlag int) error {
	_, err := s.organizationMembershipRepository.FindUserInOrganization(organization.ID, user.ID)
	if err == nil {
		// OrganizationMembership exists, throw error
		return NewErrUserAlreadyInOrganization()
	}

	// Generate an ID for the invite.
	id := s.snowflakeService.GenerateID()
	key := generateKey()
	err = s.organizationMembershipInviteRepository.Create(models.OrganizationMembershipInvite{
		Model: models.Model{
			ID: id,
		},
		Accepted:       false,
		InviteeID:      user.ID,
		OrganizationID: organization.ID,
		InviterID:      inviterID,
		Key:            key,
	})
	if err != nil {
		return NewErrServerError()
	}

	// Create a new email with the reset password template.
	email := s.emailService.NewEmail(
		email.NewRecipient(fmt.Sprintf("%s %s", user.FirstName, user.LastName), user.Email),
		fmt.Sprintf("You've been invited to join %s!", organization.Name),
		templates.OrganizationInvitationTemplate(user.FirstName, organization.Name, organization.ID, id, key),
	)
	err = s.emailService.Send(email)
	if err != nil {
		return NewErrServerError()
	}

	return nil
}

// GetOrganizationInvitedVolunteers returns an array of OrganizationMembershipRequest objects for a specified organization by ID.
func (s *service) GetOrganizationInvitedVolunteers(ctx context.Context, organizationID int64) ([]models.OrganizationMembershipInvite, error) {
	// Get all membership invites by opportunity ID.
	invites, err := s.organizationMembershipInviteRepository.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, NewErrServerError()
	}

	return invites, nil
}

// GetOrganizationMemberships gets all users in an organization.
func (s *service) GetOrganizationMemberships(organizationID int64) ([]models.OrganizationMembership, error) {
	memberships, err := s.organizationMembershipRepository.FindByOrganizationID(organizationID)
	if err != nil {
		return nil, err
	}

	return memberships, nil
}

// GetOrganizationMembership returns the membership level of a user in an organization.
// Returns an error if no membership is found.
func (s *service) GetOrganizationMembership(organizationID, userID int64) (int, error) {
	m, err := s.organizationMembershipRepository.FindUserInOrganization(organizationID, userID)
	if err != nil {
		return 0, err
	}

	return m.PermissionsFlag, nil
}

// GetOrganizationFromInvite gets an organization profile from an invite for UI use.
func (s *service) GetOrganizationFromInvite(ctx context.Context, organizationID int64, userID, inviteID int64, inviteKey string) (*OrganizationProfile, error) {
	// Get the invite by ID.
	invite, err := s.organizationMembershipInviteRepository.FindByID(inviteID)
	if err != nil {
		return nil, NewErrInviteInvalid()
	}

	// Get the user by ID to check the email.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, NewErrInviteInvalid()
	}

	// Check if the invite is valid.
	if (invite.InviteeEmail != user.Email && invite.InviteeID != user.ID) || invite.Key != inviteKey || invite.OrganizationID != organizationID {
		return nil, NewErrInviteInvalid()
	}

	return s.getMinimumOrganizationProfile(invite.OrganizationID)
}

// AcceptInvite accepts an invite.
func (s *service) AcceptInvite(ctx context.Context, organizationID int64, userID, inviteID int64, inviteKey string) error {
	// Get the invite by ID.
	invite, err := s.organizationMembershipInviteRepository.FindByID(inviteID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Get the user by ID to check the email.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Check if the invite is valid.
	if (invite.InviteeEmail != user.Email && invite.InviteeID != user.ID) || invite.Key != inviteKey || invite.OrganizationID != organizationID {
		return NewErrInviteInvalid()
	}

	if err := s.createVolunteerMembership(ctx, invite.InviterID, invite.OrganizationID, userID); err != nil {
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
func (s *service) DeclineInvite(ctx context.Context, organizationID int64, userID, inviteID int64, inviteKey string) error {
	// Get the invite by ID.
	invite, err := s.organizationMembershipInviteRepository.FindByID(inviteID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Get the user by ID to check the email.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return NewErrInviteInvalid()
	}

	// Check if the invite is valid.
	if (invite.InviteeEmail != user.Email && invite.InviteeID != user.ID) || invite.Key != inviteKey || invite.OrganizationID != organizationID {
		return NewErrInviteInvalid()
	}

	if err := s.createVolunteerMembership(ctx, invite.InviterID, invite.OrganizationID, userID); err != nil {
		s.logger.Error().Err(err).Msg("Error creating volunteer membership")
		return NewErrServerError()
	}

	if err := s.invalidateInvite(ctx, invite.ID); err != nil {
		s.logger.Error().Err(err).Msg("Error invalidating invite")
		return NewErrServerError()
	}

	return nil
}

// createVolunteerMembership creates a volunteer membership in an organization.
func (s *service) createVolunteerMembership(ctx context.Context, inviterID int64, organizationID, userID int64) error {
	membership := models.OrganizationMembership{}

	membership.Active = true
	membership.ID = s.snowflakeService.GenerateID()
	membership.UserID = userID
	membership.JoinedAt = time.Now()
	membership.OrganizationID = organizationID
	membership.PermissionsFlag = models.OpportunityPermissionsMember
	membership.InviterID = inviterID

	return s.organizationMembershipRepository.Create(membership)
}

// invalidateInvite invalidates an organization invite by ID.
func (s *service) invalidateInvite(ctx context.Context, inviteID int64) error {
	return s.organizationMembershipInviteRepository.DeleteByID(inviteID)
}

// GetOrganizationVolunteers gets all volunteers from a specified organization.
func (s *service) GetOrganizationVolunteers(ctx context.Context, organizationID int64) error {
	return nil
}
