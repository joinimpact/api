package authentication

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthTokenLifespanDays represents how long (in days) an auth token lasts before expiring.
const AuthTokenLifespanDays = 3

// RefreshTokenLifespanDays represents how long (in days) a refresh token lasts before expiring.
const RefreshTokenLifespanDays = 21

// TokenPair represents a pair of an auth token and refresh token.
type TokenPair struct {
	AuthToken     string `json:"authToken"`
	AuthExpiry    int64  `json:"authExpiry"`
	RefreshToken  string `json:"refreshToken"`
	RefreshExpiry int64  `json:"refreshExpiry"`
	UserID        int64  `json:"userId,omitempty"`
}

const (
	// AuthTokenType represents the type of an auth token.
	AuthTokenType = iota
	// RefreshTokenType represents the type of a refresh token.
	RefreshTokenType = iota
)

// jwtClaims contains the claims that are used in generated JWT tokens.
type jwtClaims struct {
	UserID int64 `json:"userId"`
	Type   int   `json:"type"`
	jwt.StandardClaims
}

// generateTokenPair generates a new TokenPair which includes an auth and refresh token using the user's ID.
func (s *service) generateTokenPair(userID int64) (*TokenPair, error) {
	// Ensure that there is a jwt secret present.
	if len(s.config.JWTSecret) <= 32 {
		s.logger.Error().Msg("no jwt secret present")
		return nil, errors.New("no jwt secret present")
	}

	now := time.Now()

	// Append the date to the secret for rolling secrets.
	secret := fmt.Sprintf("%s::%d", s.config.JWTSecret, now.UTC().Unix())

	authTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaims{
		UserID: userID,
		Type:   AuthTokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.UTC().Unix(),
			ExpiresAt: now.Add(AuthTokenLifespanDays * 24 * time.Hour).UTC().Unix(),
			Issuer:    "impact-prod-01",
		},
	})
	authToken, err := authTokenObject.SignedString([]byte(secret))
	if err != nil {
		s.logger.Error().Err(err).Msg("error generating auth jwt")
		return nil, errors.New("error generating auth jwt")
	}

	refreshTokenObject := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtClaims{
		UserID: userID,
		Type:   RefreshTokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.UTC().Unix(),
			ExpiresAt: now.Add(RefreshTokenLifespanDays * 24 * time.Hour).UTC().Unix(),
			Issuer:    "impact-prod-01",
		},
	})
	refreshToken, err := refreshTokenObject.SignedString([]byte(secret))
	if err != nil {
		s.logger.Error().Err(err).Msg("error generating refresh jwt")
		return nil, errors.New("error generating refresh jwt")
	}

	return &TokenPair{
		authToken,
		now.Add(AuthTokenLifespanDays * 24 * time.Hour).UTC().Unix(),
		refreshToken,
		now.Add(RefreshTokenLifespanDays * 24 * time.Hour).UTC().Unix(),
		userID,
	}, nil
}

// claimsFromToken validates and parses a token. If the token is valid, it
// will return a jwtClaims.
func (s *service) claimsFromToken(token string) (*jwtClaims, error) {
	jwt, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		claims, ok := t.Claims.(*jwtClaims)
		if !ok {
			return nil, errors.New("error processing jwt claims")
		}

		// Append the date to the secret for rolling secrets.
		secret := fmt.Sprintf("%s::%d", s.config.JWTSecret, claims.IssuedAt)

		return []byte(secret), nil
	})

	if claims, ok := jwt.Claims.(*jwtClaims); ok && jwt.Valid {
		return claims, nil
	}

	return nil, err
}

// GetUserIDFromToken gets a user's ID from a JWT token.
func (s *service) GetUserIDFromToken(token string) (int64, error) {
	claims, err := s.claimsFromToken(token)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
