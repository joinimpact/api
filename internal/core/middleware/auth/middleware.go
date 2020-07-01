package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/joinimpact/api/internal/authentication"
	"github.com/joinimpact/api/pkg/resp"
)

// Key represents a context key
type Key int

// Keys
const (
	KeyUserID Key = iota
)

// getToken attempts to get the token from the Authorization HTTP header.
func getToken(r *http.Request) (string, error) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {
		return "", errors.New("can not get header")
	}

	return auth[1], nil
}

// AuthMiddleware adds necessary context to the handler for auth.
// TODO: move to another package.
func AuthMiddleware(authService authentication.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := getToken(r)
			if err != nil {
				resp.Unauthorized(w, r, resp.Error(401, "no authorization header present"))
				return
			}

			userID, err := authService.GetUserIDFromToken(token)
			if err != nil {
				fmt.Println(token)
				resp.Unauthorized(w, r, resp.Error(401, "invalid auth token"))
				return
			}

			ctx = context.WithValue(ctx, KeyUserID, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CookieMiddleware takes authentication information from the cookies and injects it into the Authorization header for later consumption.
func CookieMiddleware(authService authentication.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(r.Header.Get("Authorization")) > 0 {
				// Authorization headers set by the client will take priority.
				// If one exists, skip the injection.
				next.ServeHTTP(w, r)
				return
			}

			// Get the auth_token cookie from the request.
			token, err := r.Cookie("auth_token")
			if err != nil {
				// No cookie found, skip.
				next.ServeHTTP(w, r)
				return
			}

			authToken := token.Value

			_, err = authService.GetUserIDFromToken(authToken)
			if err != nil {
				// Attempt to refresh token.
				// Get the refresh_token cookie from the request.
				refreshToken, err := r.Cookie("refresh_token")
				if err != nil {
					// Clear cookies on failure.
					ClearAuthCookies(w, r)

					next.ServeHTTP(w, r)
					return
				}

				// Attempt to refresh the token.
				tokenPair, err := authService.RefreshToken(r.Context(), refreshToken.Value)
				if err != nil {
					// Clear cookies on failure.
					ClearAuthCookies(w, r)

					next.ServeHTTP(w, r)
					return
				}

				SetAuthCookies(w, r, tokenPair)
				authToken = tokenPair.AuthToken
			}

			// Set the header.
			r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

			next.ServeHTTP(w, r)
		})
	}
}
