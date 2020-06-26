package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/authentication"
	"github.com/joinimpact/api/pkg/resp"
)

// userRequestContext contains contextual information about a user needed
// for the /users routes.
type userRequestContext struct {
	userID int64
	isSelf bool
}

type key int

const (
	keyUserRequestContext key = iota
)

// getToken attempts to get the token from the Authorization HTTP header.
func getToken(r *http.Request) (string, error) {
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {
		return "", errors.New("can not get header")
	}

	return auth[1], nil
}

// Middleware adds necessary context to the /users handler.
func Middleware(authService authentication.Service) func(next http.Handler) http.Handler {
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

			queryUserID := userID
			userIDString := chi.URLParam(r, "userID")

			if userIDString != "me" {
				queryUserID, err = strconv.ParseInt(userIDString, 10, 64)
				if err != nil {
					resp.BadRequest(w, r, resp.Error(400, "invalid user ID"))
					return
				}
			}

			ctx = context.WithValue(ctx, keyUserRequestContext, &userRequestContext{
				userID: queryUserID,
				isSelf: queryUserID == userID,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
