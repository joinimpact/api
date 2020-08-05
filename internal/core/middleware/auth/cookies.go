package auth

import (
	"net/http"
	"time"

	"github.com/joinimpact/api/internal/authentication"
)

// SetAuthCookies sets token cookies.
func SetAuthCookies(w http.ResponseWriter, r *http.Request, tokenPair *authentication.TokenPair) {
	// Set auth token cookie.
	authTokenCookie := http.Cookie{
		Name:     "auth_token",
		Value:    tokenPair.AuthToken,
		Expires:  time.Unix(tokenPair.AuthExpiry, 0),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &authTokenCookie)

	// Set refresh token cookie.
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    tokenPair.RefreshToken,
		Expires:  time.Unix(tokenPair.RefreshExpiry, 0),
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &refreshTokenCookie)
}

// ClearAuthCookies clears token cookies.
func ClearAuthCookies(w http.ResponseWriter, r *http.Request) {
	// Set auth token cookie.
	authTokenCookie := http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // a MaxAge of less than 0 will effectively delete the cookie immediately.
	}
	http.SetCookie(w, &authTokenCookie)

	// Set refresh token cookie.
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // a MaxAge of less than 0 will effectively delete the cookie immediately.
	}
	http.SetCookie(w, &refreshTokenCookie)
}
