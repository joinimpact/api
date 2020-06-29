package organizations

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
	"github.com/oliamb/cutter"
)

type createOrganizationResponse struct {
	OrganizationID int64 `json:"organizationId"`
}

// CreateOrganization creates a new organization.
func CreateOrganization(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := ctx.Value(keyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		req := struct {
			Name        string `json:"name" validate:"min=4,max=72"`
			WebsiteURL  string `json:"websiteURL" validate:"url,omitempty"`
			Location    string `json:"location" validate:"omitempty"`
			Description string `json:"description" validate:"max=800,omitempty"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		id, err := organizationsService.CreateOrganization(models.Organization{
			CreatorID:   userID,
			Name:        req.Name,
			WebsiteURL:  req.WebsiteURL,
			Location:    req.Location,
			Description: req.Description,
		})
		if err != nil {
			switch err.(type) {
			case *organizations.ErrOrganizationNotFound, *organizations.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *organizations.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, createOrganizationResponse{
			id,
		})
	}
}

type organizationTagsResponse struct {
	Tags []models.Tag `json:"tags"`
}

// GetOrganizationTags gets all tags by Organization ID.
func GetOrganizationTags(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
			return
		}

		tags, err := organizationsService.GetOrganizationTags(organizationID)
		if err != nil {
			switch err.(type) {
			case *organizations.ErrOrganizationNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *organizations.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, organizationTagsResponse{tags})
	}
}

type postOrganizationTagsRequest struct {
	Name string `json:"name" validate:"min=2,max=24"`
}

type postOrganizationTagsResponse struct {
	NumberAdded int `json:"numAdded"`
}

// PostOrganizationTags adds tags to an organization's profile.
func PostOrganizationTags(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Tags []postOrganizationTagsRequest `json:"tags"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
			return
		}

		tags := []string{}
		for _, tag := range req.Tags {
			tags = append(tags, tag.Name)
		}

		numAdded, err := organizationsService.AddOrganizationTags(organizationID, tags)
		if err != nil {
			switch err.(type) {
			case *organizations.ErrOrganizationNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *organizations.ErrServerError:
				resp.ServerError(w, r, resp.Error(500, err.Error()))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, postOrganizationTagsResponse{numAdded})
	}
}

// DeleteOrganizationTag deletes a single organization tag by ID.
func DeleteOrganizationTag(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
			return
		}

		// Get the tagID from the URL.
		tagIDString := chi.URLParam(r, "tagID")
		tagID, err := strconv.ParseInt(tagIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid tag ID"))
			return
		}

		err = organizationsService.RemoveOrganizationTag(organizationID, tagID)
		if err != nil {
			switch err.(type) {
			case *organizations.ErrOrganizationNotFound, *organizations.ErrTagNotFound:
				resp.NotFound(w, r, resp.Error(404, err.Error()))
			case *organizations.ErrServerError:
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
func UploadProfilePicture(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
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

		err = organizationsService.UploadProfilePicture(organizationID, f)
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, err.Error()))
			return
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}

type postInviteEmail struct {
	Email string `json:"email" validate:"email"`
}

// PostInvite creates new invites from user emails.
func PostInvite(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
			return
		}

		ctx := r.Context()
		userID, ok := ctx.Value(keyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		req := struct {
			Invites []postInviteEmail `json:"invites"`
		}{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		errors := []error{}

		for _, invite := range req.Invites {
			err := organizationsService.InviteUser(userID, organizationID, invite.Email, models.OrganizationPermissionsMember)
			if err != nil {
				errors = append(errors, err)
			}
		}

		if len(errors) > 0 {
			resp.BadRequest(w, r, resp.ErrorData(32, "errors while sending invites", errors))
			return
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}
