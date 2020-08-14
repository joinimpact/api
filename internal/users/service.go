package users

import (
	"fmt"
	"io"
	"time"

	"github.com/joinimpact/api/internal/cdn"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/joinimpact/api/pkg/location"
	"github.com/rs/zerolog"
)

// Service represents a provider of User services (excluding authentication).
type Service interface {
	// GetUserProfile retrieves a single user's profile.
	GetUserProfile(userID int64, self bool) (*UserProfile, error)
	// GetMinimalUserProfile retrieves a single user's profile but skips extra fields such as tags and profile.
	GetMinimalUserProfile(userID int64) (*UserProfile, error)
	// UpdateUserProfile updates a user's profile.
	UpdateUserProfile(userID int64, profile UserProfile) error
	// UpdateUserLocation updates a user's location.
	UpdateUserLocation(userID int64, location *location.Coordinates) error
	// GetUserTags gets all of a user's tags.
	GetUserTags(userID int64) ([]models.Tag, error)
	// AddUserTags adds tags to a user by tag name.
	AddUserTags(userID int64, tags []string) (int, error)
	// RemoveUserTag removes a tag from a user by id.
	RemoveUserTag(userID, tagID int64) error
	// SetUserProfileField sets a user's profile field by name.
	SetUserProfileField(userID int64, profileField models.UserProfileField) error
	// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
	UploadProfilePicture(userID int64, fileReader io.Reader) (string, error)
}

// service represents the internal implementation of the Service interface.
type service struct {
	userRepository             models.UserRepository
	userProfileFieldRepository models.UserProfileFieldRepository
	userTagRepository          models.UserTagRepository
	tagRepository              models.TagRepository
	config                     *config.Config
	logger                     *zerolog.Logger
	snowflakeService           snowflakes.SnowflakeService
	cdnClient                  *cdn.Client
	locationService            location.Service
}

// NewService creates and returns a new Users service with the provifded dependencies.
func NewService(userRepository models.UserRepository, userProfileFieldRepository models.UserProfileFieldRepository, userTagRepository models.UserTagRepository,
	tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService, locationService location.Service) Service {
	return &service{
		userRepository,
		userProfileFieldRepository,
		userTagRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		cdn.NewCDNClient(config),
		locationService,
	}
}

// GetUserProfile retrieves a single user's profile. If self is true, it will also add sensitive fields
// such as email.
func (s *service) GetUserProfile(userID int64, self bool) (*UserProfile, error) {
	profile := &UserProfile{}
	// Find the user to verify that it is active.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, NewErrUserNotFound()
	}

	profile.ID = user.ID
	profile.FirstName = user.FirstName
	profile.LastName = user.LastName
	profile.ProfilePicture = user.ProfilePicture
	profile.LastOnline = user.LastOnline

	if self {
		profile.Email = user.Email
		profile.DateOfBirth = user.DateOfBirth
		profile.CreatedAt = user.CreatedAt
	}

	tags, err := s.GetUserTags(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	profile.Tags = tags

	// Location
	if user.LocationLatitude != 0.0 || user.LocationLongitude != 0.0 {
		coordinates := &location.Coordinates{
			Latitude:  user.LocationLatitude,
			Longitude: user.LocationLongitude,
		}

		location, err := s.locationService.CoordinatesToCity(coordinates)
		if err == nil {
			profile.Location = location
		}
	}

	// Profile fields
	fields, err := s.userProfileFieldRepository.FindByUserID(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	profile.ProfileFields = []models.UserProfileField{}

	// TODO: replace with scoped struct tags
	if self {
		profile.ProfileFields = fields
	} else {
		for _, field := range fields {
			if field.Privacy == 1 {
				profile.ProfileFields = append(profile.ProfileFields, field)
			}
		}
	}

	return profile, nil
}

// GetMinimalUserProfile retrieves a single user's profile but skips extra fields such as tags and profile.
func (s *service) GetMinimalUserProfile(userID int64) (*UserProfile, error) {
	profile := &UserProfile{}
	// Find the user to verify that it is active.
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, NewErrUserNotFound()
	}

	profile.ID = user.ID
	profile.FirstName = user.FirstName
	profile.LastName = user.LastName
	profile.ProfilePicture = user.ProfilePicture
	profile.LastOnline = user.LastOnline

	return profile, nil
}

// UpdateUserProfile updates a user's profile.
func (s *service) UpdateUserProfile(userID int64, profile UserProfile) error {
	return s.userRepository.Update(models.User{
		Model: models.Model{
			ID: userID,
		},
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		DateOfBirth: profile.DateOfBirth,
	})
}

// UpdateUserLocation updates a user's location.
func (s *service) UpdateUserLocation(userID int64, location *location.Coordinates) error {
	return s.userRepository.Update(models.User{
		Model: models.Model{
			ID: userID,
		},
		LocationLatitude:  location.Latitude,
		LocationLongitude: location.Longitude,
	})
}

// GetUserTags gets all of a user's tags.
func (s *service) GetUserTags(userID int64) ([]models.Tag, error) {
	// Find the user to verify that it is active.
	_, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, NewErrUserNotFound()
	}

	// Find all UserTag objects by UserID.
	userTags, err := s.userTagRepository.FindByUserID(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	tags := []models.Tag{}
	for _, userTag := range userTags {
		// Get the tag by ID.
		tag, err := s.tagRepository.FindByID(userTag.TagID)
		if err != nil {
			// Tag not found, skip.
			s.logger.Error().Err(err).Msg("Error in GetUserTags: UserTag object missing valid Tag")
			continue
		}

		// Append the tag to the tags array.
		tags = append(tags, *tag)
	}

	return tags, nil
}

// AddUserTags adds tags to a user by tag name.
func (s *service) AddUserTags(userID int64, tags []string) (int, error) {
	// successfulTags counts how many tags were inserted correctly.
	successfulTags := 0

	_, err := s.userRepository.FindByID(userID)
	if err != nil {
		return successfulTags, NewErrUserNotFound()
	}

	for _, tag := range tags {
		tag, err := s.tagRepository.FindByName(tag)
		if err != nil {
			// Log the error and skip the tag.
			s.logger.Error().Err(err).Msg("Error in AddUserTags finding a tag")
			continue
		}

		// Increment the successful tags value as the tag was found.
		successfulTags++

		// Check to see if the user already has this tag.
		_, err = s.userTagRepository.FindUserTagByID(userID, tag.ID)
		if err == nil {
			// The user already has this tag, skip.
			continue
		}

		// Create a new UserTag entity.
		id := s.snowflakeService.GenerateID()
		err = s.userTagRepository.Create(models.UserTag{
			Model: models.Model{
				ID: id,
			},
			UserID: userID,
			TagID:  tag.ID,
		})
		if err != nil {
			s.logger.Error().Err(err).Msg("Error in AddUserTags creating a UserTag")
			return successfulTags - 1, NewErrServerError()
		}
	}

	return successfulTags, nil
}

// RemoveUserTag removes a tag from a user by id.
func (s *service) RemoveUserTag(userID, tagID int64) error {
	userTag, err := s.userTagRepository.FindUserTagByID(userID, tagID)
	if err != nil {
		return NewErrTagNotFound()
	}

	return s.userTagRepository.DeleteByID(userTag.ID)
}

// SetUserProfileField sets a user's profile field by name.
func (s *service) SetUserProfileField(userID int64, profileField models.UserProfileField) error {
	field, err := s.userProfileFieldRepository.FindUserFieldByName(userID, profileField.Name)
	if err == nil {
		profileField.ID = field.ID
		return s.userProfileFieldRepository.Update(profileField)
	}

	profileField.UserID = userID

	// Create an ID and assign it to the profile field.
	id := s.snowflakeService.GenerateID()
	profileField.ID = id

	// Create the profile field.
	return s.userProfileFieldRepository.Create(profileField)
}

// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
func (s *service) UploadProfilePicture(userID int64, fileReader io.Reader) (string, error) {
	url, err := s.cdnClient.UploadImage(fmt.Sprintf("profile-picture-%d-%d.png", userID, time.Now().UTC().Unix()), fileReader)
	if err != nil {
		return "", err
	}

	return url, s.userRepository.Update(models.User{
		Model: models.Model{
			ID: userID,
		},
		ProfilePicture: url,
	})
}
