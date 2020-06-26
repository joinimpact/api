package organizations

import (
	"fmt"
	"io"
	"time"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
)

// Service represents a provider of Organization services.
type Service interface {
	// GetOrganizationProfile retrieves a single user's profile.
	// GetOrganizationProfile(organizationID int64) (*OrganizationProfile, error)
	// UpdateOrganizationProfile updates a user's profile.
	// UpdateOrganizationProfile(userID int64, profile OrganizationProfile) error

	// CreateOrganization creates a new organization and returns the ID on success.
	CreateOrganization(organization models.Organization) (int64, error)
	// GetOrganizationTags gets all of a user's tags.
	GetOrganizationTags(organizationID int64) ([]models.Tag, error)
	// AddOrganizationTags adds tags to a user by tag name.
	AddOrganizationTags(organizationID int64, tags []string) (int, error)
	// RemoveOrganizationTag removes a tag from an organization by id.
	RemoveOrganizationTag(organizationID, tagID int64) error
	// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
	UploadProfilePicture(organizationID int64, fileReader io.Reader) error
}

// service represents the internal implementation of the organizations Service.
type service struct {
	organizationRepository           models.OrganizationRepository
	organizationMembershipRepository models.OrganizationMembershipRepository
	organizationTagRepository        models.OrganizationTagRepository
	tagRepository                    models.TagRepository
	config                           *config.Config
	logger                           *zerolog.Logger
	snowflakeService                 snowflakes.SnowflakeService
	cdnClient                        *cdn.Client
}

// NewService creates and returns a new Users service with the provifded dependencies.
func NewService(organizationRepository models.OrganizationRepository, organizationMembershipRepository models.OrganizationMembershipRepository, organizationTagRepository models.OrganizationTagRepository,
	tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService) Service {
	return &service{
		organizationRepository,
		organizationMembershipRepository,
		organizationTagRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		cdn.NewCDNClient(config),
	}
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
		PermissionsFlag: 3,
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
func (s *service) UploadProfilePicture(organizationID int64, fileReader io.Reader) error {
	url, err := s.cdnClient.UploadImage(fmt.Sprintf("organization-picture-%d-%d.png", organizationID, time.Now().UTC().Unix()), fileReader)
	if err != nil {
		return err
	}

	return s.organizationRepository.Update(models.Organization{
		Model: models.Model{
			ID: organizationID,
		},
		ProfilePicture: url,
	})
}
