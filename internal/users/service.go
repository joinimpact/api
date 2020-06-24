package users

import (
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
)

// Service represents a provider of User services (excluding authentication).
type Service interface {
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
	}
}
