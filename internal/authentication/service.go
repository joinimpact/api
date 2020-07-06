package authentication

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/joinimpact/api/internal/authentication/oauth"
	"github.com/joinimpact/api/internal/config"
	"github.com/joinimpact/api/internal/email"
	"github.com/joinimpact/api/internal/email/templates"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/snowflakes"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

// Service represents a provider of authentication services.
type Service interface {
	// Login attempts to login to a user account with a email and password, and returns a TokenPair on success.
	Login(email, password string) (*TokenPair, error)
	// CheckEmail checks if an email is available for use. If the email is taken,
	// an error will be returned.
	CheckEmail(email string) error
	// Register attempts to create a new user and returns a token pair on success.
	Register(user models.User, password string) (*TokenPair, error)
	// RequestPasswordReset creates a new PasswordResetKey and emails the user a link to it.
	RequestPasswordReset(userEmail string) error
	// CheckPasswordReset gets a PasswordResetKey by its key for validation purposes.
	CheckPasswordReset(key string) (*PasswordResetValidation, error)
	// ResetPassword resets a user's password from a PasswordResetKey's key.
	ResetPassword(key string, newPassword string) error
	// GetUserIDFromToken gets a user's ID from a JWT token.
	GetUserIDFromToken(token string) (int64, error)
	// OauthLogin authenticates using a third-party service instead of a traditional username and password.
	OauthLogin(serviceName, accessToken string) (*OauthResponse, error)
	// RefreshToken generates a new token pair from a refresh token.
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
}

// service represents the default authentication service of this package.
type service struct {
	userRepository               models.UserRepository
	passwordResetKeyRepository   models.PasswordResetKeyRepository
	thirdPartyIdentityRepository models.ThirdPartyIdentityRepository
	config                       *config.Config
	logger                       *zerolog.Logger
	snowflakeService             snowflakes.SnowflakeService
	emailService                 email.Service
}

// NewService creates and returns a new Service with the provided UserRepository, Config, Logger, and SnowflakeService.
func NewService(userRepository models.UserRepository, passwordResetKeyRepository models.PasswordResetKeyRepository, thirdPartyIdentityRepository models.ThirdPartyIdentityRepository, config *config.Config, logger *zerolog.Logger,
	snowflakeService snowflakes.SnowflakeService, emailService email.Service) Service {
	return &service{
		userRepository,
		passwordResetKeyRepository,
		thirdPartyIdentityRepository,
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

// CheckEmail checks if an email is available for use. If the email is taken,
// an error will be returned.
func (s *service) CheckEmail(email string) error {
	if !validateEmail(email) {
		return errors.New("invalid email")
	}

	_, err := s.userRepository.FindByEmail(email)
	if err == nil {
		return errors.New("email taken")
	}

	return nil
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

	_, err := s.userRepository.FindByEmail(user.Email)
	if err == nil {
		return nil, errors.New("email taken")
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

	email := s.emailService.NewEmail(email.NewRecipient(
		fmt.Sprintf("%s %s", user.FirstName, user.LastName), user.Email),
		"Welcome to Impact!",
		templates.WelcomeTemplate(user.FirstName),
	)
	err = s.emailService.Send(email)
	if err != nil {
		s.logger.Error().Err(err).Msg("error sending email to new user")
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

// OauthResponse has information relating to the autentication of users using Oauth.
type OauthResponse struct {
	UserCreated bool       `json:"userCreated"`
	TokenPair   *TokenPair `json:"token"`
}

// OauthLogin authenticates using a third-party service instead of a traditional username and password.
func (s *service) OauthLogin(serviceName, accessToken string) (*OauthResponse, error) {
	var profile oauth.Profile
	var token *oauth2.Token
	var err error

	switch serviceName {
	case "google":
		client := oauth.NewGoogleClient(s.config)
		token = client.GetTokenFromAccessToken(accessToken)
		profile, err = client.GetProfile(token)
		if err != nil {
			return nil, err
		}
	case "facebook":
		client := oauth.NewFacebookClient(s.config)
		token = client.GetTokenFromAccessToken(accessToken)
		profile, err = client.GetProfile(token)
		if err != nil {
			return nil, err
		}
	}

	if profile == nil {
		return nil, errors.New("invalid service name")
	}

	if !validateEmail(profile.GetEmail()) {
		return nil, errors.New("could not get user email")
	}

	user, created, err := s.createOauthUserIfNotExists(profile)
	if err != nil {
		return nil, err
	}

	identity, err := s.thirdPartyIdentityRepository.FindUserIdentityByServiceName(user.ID, serviceName)
	if err != nil {
		identity = &models.ThirdPartyIdentity{
			UserID:                 user.ID,
			ThirdPartyServiceName:  serviceName,
			ThirdPartyAccessToken:  token.AccessToken,
			ThirdPartyRefreshToken: token.RefreshToken,
		}

		identity.ID = s.snowflakeService.GenerateID()

		err = s.thirdPartyIdentityRepository.Create(*identity)
	}
	if err != nil {
		return nil, err
	}

	// Successful login, generate and return a token pair.
	userToken, err := s.generateTokenPair(user.ID)
	if err != nil {
		return nil, err
	}

	return &OauthResponse{
		UserCreated: created,
		TokenPair:   userToken,
	}, nil
}

// createOauthUserIfNotExists takes an oauth profile and checks to see if a user already exists.
// If one exists, it will return it with the bool false,
// and if not, it will return a newly created user with the bool true.
func (s *service) createOauthUserIfNotExists(profile oauth.Profile) (*models.User, bool, error) {
	user, err := s.userRepository.FindByEmail(profile.GetEmail())
	if err == nil {
		return user, false, nil
	}

	// Create the new user around the oauth values.
	newUser := models.User{
		Email:     strings.ToLower(profile.GetEmail()),
		FirstName: profile.GetFirstName(),
		LastName:  profile.GetLastName(),
	}

	// Generate an ID for the new user.
	newUser.ID = s.snowflakeService.GenerateID()

	// Create the user in the repository.
	err = s.userRepository.Create(newUser)
	if err != nil {
		return nil, false, err
	}

	email := s.emailService.NewEmail(email.NewRecipient(
		fmt.Sprintf("%s %s", newUser.FirstName, newUser.LastName), newUser.Email),
		"Welcome to Impact!",
		templates.WelcomeTemplate(newUser.FirstName),
	)
	err = s.emailService.Send(email)
	if err != nil {
		s.logger.Error().Err(err).Msg("error sending email to new user (oauth)")
	}

	return &newUser, true, nil
}

// RefreshToken generates a new token pair from a refresh token.
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	claims, err := s.claimsFromToken(refreshToken)
	if err != nil {
		// Error validating/parsing token.
		return nil, err
	}

	if claims.Type != RefreshTokenType {
		return nil, errors.New("not a refresh token")
	}

	// No errors, generate a token pair.
	return s.generateTokenPair(claims.UserID)
}
