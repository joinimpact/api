package users

import (
	"net/http"

	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/resp"
)

// DeletePost deletes a user by ID.
func DeletePost(usersService users.Service) http.HandlerFunc {
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		err = usersService.DeleteUser(ctx, userID)
		if err != nil {
			switch err.(type) {
			case *users.ErrUserNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *users.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{true})
	}
}
