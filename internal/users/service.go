package users

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

// Service represents a provider of User services (excluding authentication).
type Service interface {
	// GetUserProfile retrieves a single user's profile.
	GetUserProfile(userID int64, self bool) (*UserProfile, error)
	// UpdateUserProfile updates a user's profile.
	UpdateUserProfile(userID int64, profile UserProfile) error
	// GetUserTags gets all of a user's tags.
	GetUserTags(userID int64) ([]models.Tag, error)
	// AddUserTags adds tags to a user by tag name.
	AddUserTags(userID int64, tags []string) (int, error)
	// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
	UploadProfilePicture(userID int64, fileReader io.Reader) error
	// RemoveUserTag removes a tag from a user by id.
	RemoveUserTag(userID, tagID int64) error
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
}

// NewService creates and returns a new Users service with the provifded dependencies.
func NewService(userRepository models.UserRepository, userProfileFieldRepository models.UserProfileFieldRepository, userTagRepository models.UserTagRepository,
	tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService) Service {
	return &service{
		userRepository,
		userProfileFieldRepository,
		userTagRepository,
		tagRepository,
		config,
		logger,
		snowflakeService,
		cdn.NewCDNClient(config),
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

	profile.FirstName = user.FirstName
	profile.LastName = user.LastName
	profile.ProfilePicture = user.ProfilePicture

	if self {
		profile.Email = user.Email
		profile.DateOfBirth = user.DateOfBirth
		profile.ZIPCode = user.ZIPCode
	}

	tags, err := s.GetUserTags(userID)
	if err != nil {
		return nil, NewErrServerError()
	}

	profile.Tags = tags

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
		ZIPCode:     profile.ZIPCode,
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

// UploadProfilePicture uploads a profile picture to the CDN and adds it to the user.
func (s *service) UploadProfilePicture(userID int64, fileReader io.Reader) error {
	url, err := s.cdnClient.UploadImage(fmt.Sprintf("profile-picture-%d-%d.png", userID, time.Now().UTC().Unix()), fileReader)
	if err != nil {
		return err
	}

	return s.userRepository.Update(models.User{
		Model: models.Model{
			ID: userID,
		},
		ProfilePicture: url,
	})
}
