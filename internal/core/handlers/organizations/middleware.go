package organizations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/joinimpact/api/internal/authentication"
	"github.com/joinimpact/api/pkg/resp"
)

type key int

const (
	keyUserID key = iota
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

			ctx = context.WithValue(ctx, keyUserID, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
