package users

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
	"github.com/oliamb/cutter"
)

// GetUserProfile gets a user's public profile.
func GetUserProfile(usersService users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqCtx, ok := ctx.Value(keyUserRequestContext).(*userRequestContext)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		profile, err := usersService.GetUserProfile(reqCtx.userID, reqCtx.isSelf)
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

		resp.OK(w, r, profile)
	}
}

// UpdateUserProfile updates a user's profile.
func UpdateUserProfile(usersService users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		req := struct {
			FirstName   string    `json:"firstName" validate:"omitempty,min=2,max=48"`
			LastName    string    `json:"lastName" validate:"omitempty,min=2,max=48"`
			DateOfBirth time.Time `json:"dateOfBirth" validate:"omitempty"`
			ZIPCode     string    `json:"zipCode" validate:"omitempty,min=5,max=8"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = usersService.UpdateUserProfile(reqCtx.userID, users.UserProfile{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			DateOfBirth: req.DateOfBirth,
			ZIPCode:     req.ZIPCode,
		})
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

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}

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

// DeleteUserTag deletes a single user tag by ID.
func DeleteUserTag(usersService users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Get the tagID from the URL.
		tagIDString := chi.URLParam(r, "tagID")
		tagID, err := strconv.ParseInt(tagIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid tag ID"))
			return
		}

		err = usersService.RemoveUserTag(reqCtx.userID, tagID)
		if err != nil {
			switch err.(type) {
			case *users.ErrUserNotFound, *users.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *users.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}

// UploadProfilePicture handles the file upload of a new profile picture.
func UploadProfilePicture(usersService users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid file"))
			return
		}
		defer file.Close()

		tmpfile, err := ioutil.TempFile("", "image-upload.*.png")
		if err != nil {
			fmt.Println(err)
			resp.ServerError(w, r, resp.Error(500, "server error"))
			return
		}
		defer os.Remove(tmpfile.Name())

		switch handler.Header.Get("Content-Type") {
		case "image/png", "image/jpeg":
			image, _, err := image.Decode(file)
			if err != nil {
				resp.BadRequest(w, r, resp.Error(400, "invalid file"))
				return
			}

			cropped, err := cutter.Crop(image, cutter.Config{
				Width:   1,
				Height:  1,
				Mode:    cutter.Centered,
				Options: cutter.Ratio,
			})
			if err != nil {
				resp.ServerError(w, r, resp.Error(500, err.Error()))
				return
			}

			err = png.Encode(tmpfile, cropped)
			if err != nil {
				resp.ServerError(w, r, resp.Error(500, "error encoding image"))
				return
			}
			fmt.Println("encoded")
		default:
			resp.BadRequest(w, r, resp.Error(400, "invalid file"))
			return
		}

		f, err := os.Open(tmpfile.Name())
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, "error encoding image"))
			return
		}

		err = usersService.UploadProfilePicture(reqCtx.userID, f)
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, err.Error()))
			return
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}
