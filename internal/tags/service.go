package tags

import (
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
)

// Service defines the interface for a Tags service.
type Service interface {
	// GetTags returns tags with a specified search string (for name matching).
	// If the string is blank, random tags will be returned.
	GetTags(query string, limit int) ([]models.Tag, error)
	// CreateTag creates a new tag with the specified name and category and returns the id.
	CreateTag(tag models.Tag) (int64, error)
}

// service represents the internal implementation of the Tags Service.
type service struct {
	tagRepository    models.TagRepository
	config           *config.Config
	logger           *zerolog.Logger
	snowflakeService snowflakes.SnowflakeService
}

// NewService creates and returns a new tags Service.
func NewService(tagRepository models.TagRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService) Service {
	return &service{
		tagRepository,
		config,
		logger,
		snowflakeService,
	}
}

// GetTags returns tags with a specified search string (for name matching).
// If the string is blank, random tags will be returned.
func (s *service) GetTags(query string, limit int) ([]models.Tag, error) {
	return s.tagRepository.SearchTags(query, limit)
}

// CreateTag creates a new tag with the specified name and category and returns the id.
func (s *service) CreateTag(tag models.Tag) (int64, error) {
	tag.ID = s.snowflakeService.GenerateID()
	err := s.tagRepository.Create(tag)
	if err != nil {
		return 0, err
	}

	return tag.ID, nil
}
