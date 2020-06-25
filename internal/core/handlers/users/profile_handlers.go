package users

import (
	"net/http"

	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

type userTagsResponse struct {
	Tags []models.Tag `json:"tags"`
}

// GetUserTags gets all tags by User ID.
func GetUserTags(usersService users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqCtx, ok := ctx.Value(keyUserRequestContext).(*userRequestContext)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		tags, err := usersService.GetUserTags(reqCtx.userID)
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

		resp.OK(w, r, userTagsResponse{tags})
	}
}

type postUserTagsRequest struct {
	Name string `json:"name" validate:"min=2,max=24"`
}

type postUserTagsResponse struct {
	NumberAdded int `json:"numAdded"`
}

// PostUserTags adds tags to a user's profile.
func PostUserTags(usersService users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Tags []postUserTagsRequest `json:"tags"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		ctx := r.Context()
		reqCtx, ok := ctx.Value(keyUserRequestContext).(*userRequestContext)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}
		if !reqCtx.isSelf {
			resp.Forbidden(w, r, resp.Error(403, "forbidden user"))
			return
		}

		tags := []string{}
		for _, tag := range req.Tags {
			tags = append(tags, tag.Name)
		}

		numAdded, err := usersService.AddUserTags(reqCtx.userID, tags)
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

		resp.OK(w, r, postUserTagsResponse{numAdded})
	}
}
