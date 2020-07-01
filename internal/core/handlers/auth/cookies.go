package auth

import (
	"net/http"
	"time"

	"github.com/joinimpact/api/internal/authentication"
)

// setAuthCookie sets token cookies.
func setAuthCookie(w http.ResponseWriter, r *http.Request, tokenPair *authentication.TokenPair) {
	// Set auth token cookie.
	authTokenCookie := http.Cookie{
		Name:     "auth_token",
		Value:    tokenPair.AuthToken,
		Expires:  time.Unix(tokenPair.AuthExpiry, 0),
		Path:     "/api/v1",
		HttpOnly: true,
	}
	http.SetCookie(w, &authTokenCookie)

	// Set refresh token cookie.
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Expires:  time.Unix(tokenPair.RefreshExpiry, 0),
		Path:     "/api/v1",
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshTokenCookie)
}
