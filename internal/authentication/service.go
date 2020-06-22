package authentication

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/email/templates"
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
	// RequestPasswordReset creates a new PasswordResetKey and emails the user a link to it.
	RequestPasswordReset(userEmail string) error
	// CheckPasswordReset gets a PasswordResetKey by its key for validation purposes.
	CheckPasswordReset(key string) (*PasswordResetValidation, error)
	// ResetPassword resets a user's password from a PasswordResetKey's key.
	ResetPassword(key string, newPassword string) error
}

// service represents the default authentication service of this package.
type service struct {
	userRepository             models.UserRepository
	passwordResetKeyRepository models.PasswordResetKeyRepository
	config                     *config.Config
	logger                     *zerolog.Logger
	snowflakeService           snowflakes.SnowflakeService
	emailService               email.Service
}

// NewService creates and returns a new Service with the provided UserRepository, Config, Logger, and SnowflakeService.
func NewService(userRepository models.UserRepository, passwordResetKeyRepository models.PasswordResetKeyRepository, config *config.Config, logger *zerolog.Logger,
	snowflakeService snowflakes.SnowflakeService, emailService email.Service) Service {
	return &service{
		userRepository,
		passwordResetKeyRepository,
		config,
		logger,
		snowflakeService,
		emailService,
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

// RequestPasswordReset creates a new PasswordResetKey and emails the user a link to it.
func (s *service) RequestPasswordReset(userEmail string) error {
	// Find the user by email.
	user, err := s.userRepository.FindByEmail(userEmail)
	if err != nil {
		return errors.New("invalid user")
	}

	// Generate an id and a key for the PasswordResetKey.
	id := s.snowflakeService.GenerateID()
	key := generatePasswordResetKey()

	// Create the PasswordResetKey.
	err = s.passwordResetKeyRepository.Create(models.PasswordResetKey{
		Model: models.Model{
			ID: id,
		},
		UserID:    user.ID,
		Key:       key,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour),
	})
	if err != nil {
		return errors.New("could not create reset key")
	}

	// Create a new email with the reset password template.
	email := s.emailService.NewEmail(
		email.NewRecipient(fmt.Sprintf("%s %s", user.FirstName, user.LastName), user.Email),
		"Your password reset link",
		templates.ResetPasswordTemplate(user.FirstName, user.Email, key),
	)

	err = s.emailService.Send(email)
	if err != nil {
		return errors.New("error sending email")
	}

	return nil
}

// CheckPasswordReset gets a PasswordResetKey by its key for validation purposes.
func (s *service) CheckPasswordReset(key string) (*PasswordResetValidation, error) {
	resetKey, err := s.passwordResetKeyRepository.FindByKey(key)
	if err != nil {
		return nil, errors.New("invalid key")
	}

	user, err := s.userRepository.FindByID(resetKey.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &PasswordResetValidation{
		FirstName: user.FirstName,
		Email:     user.Email,
	}, nil
}

// ResetPassword resets a user's password from a PasswordResetKey's key.
func (s *service) ResetPassword(key string, newPassword string) error {
	resetKey, err := s.passwordResetKeyRepository.FindByKey(key)
	if err != nil {
		return errors.New("invalid key")
	}

	// Hash the new password.
	hash, err := generateFromPassword(newPassword)
	if err != nil {
		return errors.New("could not hash password")
	}

	if err = s.userRepository.Update(models.User{
		Model: models.Model{
			ID: resetKey.UserID,
		},
		Password: hash,
	}); err != nil {
		return err
	}

	err = s.passwordResetKeyRepository.DeleteByID(resetKey.ID)

	return err
}
