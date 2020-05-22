package interfaces

import "github.com/joinimpact/api/internal/models"

// UserService defines all methods to be implemented in the user service.
type UserService interface {
	FindByID(id int64) (*models.User, error)
}
