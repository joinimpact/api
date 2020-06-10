package users

import (
	"github.com/joinimpact/api/internal/models"
	"github.com/rs/zerolog"
)

// UserService represents a service for the usecases regarding users.
type UserService struct {
	repository UserRepository
	logger     *zerolog.Logger
}

// FindByID finds a single user by ID and returns it.
func (s *UserService) FindByID(id int64) (*models.User, error) {
	return s.repository.FindByID(id)
}
