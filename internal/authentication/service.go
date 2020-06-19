package authentication

import (
	"errors"
	"strings"

	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
)

// Service represents a provider of authentication services.
type Service interface {
	// Login attempts to login to a user account with a email and password, and returns a TokenPair on success.
	Login(email, password string) (*TokenPair, error)
	// Register attempts to create a new user and returns a token pair on success.
	Register(user models.User, password string) (*TokenPair, error)
}

// service represents the default authentication service of this package.
type service struct {
	userRepository   models.UserRepository
	config           *config.Config
	logger           *zerolog.Logger
	snowflakeService snowflakes.SnowflakeService
}

// NewService creates and returns a new Service with the provided UserRepository, Config, Logger, and SnowflakeService.
func NewService(userRepository models.UserRepository, config *config.Config, logger *zerolog.Logger, snowflakeService snowflakes.SnowflakeService) Service {
	return &service{
		userRepository,
		config,
		logger,
		snowflakeService,
	}
}

// Login attempts to login to a user account with an email and password, and returns a TokenPair on success.
func (s *service) Login(email, password string) (*TokenPair, error) {
	// Find the user by email.
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if len(user.Password) < 50 {
		// Invalid/null password on user, block login.
		return nil, errors.New("user password not set")
	}

	// Compare the provided password with the user password.
	ok, err := compareHashAndPassword(password, user.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		// TODO: Custom type errors
		return nil, errors.New("invalid password")
	}

	// Successful login, generate and return a token pair.
	return s.generateTokenPair(user.ID)
}

// Register attempts to register a new user account with an email and password, and returns a TokenPair on success.
func (s *service) Register(user models.User, password string) (*TokenPair, error) {
	if user.ID > 0 {
		// Throw an error if the user already has an ID,
		// possible indication that it already exists and this
		// function was called in error.
		return nil, errors.New("user already has id")
	}

	if !validateEmail(user.Email) {
		return nil, errors.New("invalid email")
	}

	// Make the user's email lowercase for standardization.
	user.Email = strings.ToLower(user.Email)

	if !validatePassword(password) {
		return nil, errors.New("invalid password")
	}

	// Generate a new ID for the new User.
	user.ID = s.snowflakeService.GenerateID()

	// Hash the password and add it to the User.
	hash, err := generateFromPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hash

	// Create the user in the repository.
	err = s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	// Successful user creation, generate and return a token pair.
	return s.generateTokenPair(user.ID)
}
