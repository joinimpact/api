package users

import (
	"errors"
	"net/http"
	"strings"

	"github.com/joinimpact/api/internal/authentication"
)

// userRequestContext contains contextual information about a user needed
// for the /users routes.
type userRequestContext struct {
	userID int64
	isSelf bool
}

const (
	keyUserRequestContext = iota
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

			// userID := chi.URLParam(r, "userID")
			// if userID != "me" {
			// 	id, err := strconv.ParseInt(userID, 10, 64)
			// 	if err != nil {
			// 		resp.BadRequest(w, r, resp.Error(400, "invalid user ID"))
			// 		return
			// 	}
			// 	ctx = context.WithValue(ctx, keyUserRequestContext, userRequestContext{

			// 	})
			// }
			// token, err := getToken(r)
			// if err != nil {

			// }

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
