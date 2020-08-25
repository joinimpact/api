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
	"github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/opportunities"
	"github.com/joinimpact/api/internal/organizations"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/location"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
	"github.com/oliamb/cutter"
)

// GetUserOrganizations gets organizations by User ID.
func GetUserOrganizations(organizationsService organizations.Service) http.HandlerFunc {
	type response struct {
		Organizations []organizations.OrganizationProfile `json:"organizations"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := idctx.Get(r, "userID")
		if err != nil {
			return
		}

		res, err := organizationsService.GetUserOrganizations(userID)
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

		resp.OK(w, r, response{
			Organizations: res,
		})
	}
}

// GetOrganizationProfile gets a profile for an organization by ID.
func GetOrganizationProfile(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
			return
		}

		profile, err := organizationsService.GetOrganizationProfile(organizationID)
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

		resp.OK(w, r, profile)
	}
}

// UpdateOrganizationProfile updates an organization's profile.
func UpdateOrganizationProfile(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationID, err := idctx.Get(r, "organizationID")
		if err != nil {
			return
		}

		req := struct {
			Name        string                            `json:"name" validate:"min=4,max=72"`
			WebsiteURL  string                            `json:"websiteURL" validate:"url,omitempty"`
			Location    *location.Coordinates             `json:"location"`
			Description string                            `json:"description" validate:"max=800,omitempty"`
			Profile     []models.OrganizationProfileField `json:"profile"`
		}{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = organizationsService.UpdateOrganizationProfile(organizationID, organizations.OrganizationProfile{
			Name:        req.Name,
			WebsiteURL:  req.WebsiteURL,
			Description: req.Description,
		})
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

		if req.Location != nil {
			err = organizationsService.UpdateOrganizationLocation(organizationID, req.Location)
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
		}

		for _, field := range req.Profile {
			err := organizationsService.SetOrganizationProfileField(organizationID, field)
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
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}

type createOrganizationResponse struct {
	Success        bool  `json:"success"`
	OrganizationID int64 `json:"organizationId"`
}

// CreateOrganization creates a new organization.
func CreateOrganization(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		req := struct {
			Name        string                            `json:"name" validate:"min=4,max=72"`
			WebsiteURL  string                            `json:"websiteURL" validate:"omitempty,url"`
			Location    *location.Coordinates             `json:"location"`
			Description string                            `json:"description" validate:"max=800,omitempty"`
			Profile     []models.OrganizationProfileField `json:"profile"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		id, err := organizationsService.CreateOrganization(models.Organization{
			CreatorID:   userID,
			Name:        req.Name,
			WebsiteURL:  req.WebsiteURL,
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

		if req.Location != nil {
			err = organizationsService.UpdateOrganizationLocation(id, req.Location)
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
		}

		for _, field := range req.Profile {
			err := organizationsService.SetOrganizationProfileField(id, field)
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
		}

		resp.OK(w, r, createOrganizationResponse{
			true,
			id,
		})
	}
}

// DeleteOrganization deletes a single organization by ID.
func DeleteOrganization(organizationsService organizations.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
			return
		}

		err = organizationsService.DeleteOrganization(organizationID)
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

		resp.OK(w, r, map[string]bool{
			"success": true,
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

		url, err := organizationsService.UploadProfilePicture(organizationID, f)
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, err.Error()))
			return
		}

		resp.OK(w, r, map[string]interface{}{
			"success":        true,
			"profilePicture": url,
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
		userID, ok := ctx.Value(auth.KeyUserID).(int64)
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

// MembersGet gets all members in a single opportunity.
func MembersGet(organizationsService organizations.Service, usersService users.Service) http.HandlerFunc {
	type membership struct {
		models.OrganizationMembership
		users.UserProfile
	}
	type invitedMember struct {
		EmailOnly bool `json:"emailOnly"`
		models.OrganizationMembershipInvite
		users.UserProfile
	}
	type response struct {
		Members []membership    `json:"members"`
		Invited []invitedMember `json:"invited"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		organizationIDString := chi.URLParam(r, "organizationID")
		organizationID, err := strconv.ParseInt(organizationIDString, 10, 64)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, "invalid organization id"))
			return
		}

		memberships, err := organizationsService.GetOrganizationMemberships(organizationID)
		if err != nil {
			resp.ServerError(w, r, resp.ErrorRef(500, "error getting members", "generic.server_error", nil))
		}

		members := []membership{}
		for _, member := range memberships {
			user, err := usersService.GetMinimalUserProfile(member.UserID)
			if err != nil {
				// resp.ServerError(w, r, resp.ErrorRef(500, "error getting user", "generic.server_error", nil))
				continue
			}

			members = append(members, membership{member, *user})
		}

		invites, err := organizationsService.GetOrganizationInvitedVolunteers(r.Context(), organizationID)
		if err != nil {
			resp.ServerError(w, r, resp.ErrorRef(500, "error getting members", "generic.server_error", nil))
		}

		invitedMembers := []invitedMember{}
		for _, invite := range invites {
			user, err := usersService.GetMinimalUserProfile(invite.InviteeID)
			if err != nil {
				invitedMembers = append(invitedMembers, invitedMember{true, invite, users.UserProfile{}})
				continue
			}

			invitedMembers = append(invitedMembers, invitedMember{false, invite, *user})
		}

		resp.OK(w, r, response{members, invitedMembers})
	}
}

// InviteValidatePost validates an invite and returns an organization profile on success.
func InviteValidatePost(organizationsService organizations.Service) http.HandlerFunc {
	type request struct {
		Key string `json:"key" validate:"min=8,max=128"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		organizationID, err := idctx.Get(r, "organizationID")
		if err != nil {
			return
		}

		inviteID, err := idctx.Get(r, "inviteID")
		if err != nil {
			return
		}

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		res, err := organizationsService.GetOrganizationFromInvite(ctx, organizationID, userID, inviteID, req.Key)
		if err != nil {
			switch err.(type) {
			case *organizations.ErrInviteInvalid:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *organizations.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, res)
	}
}

// InviteAcceptPost attempts to accept an invite.
func InviteAcceptPost(organizationsService organizations.Service) http.HandlerFunc {
	type request struct {
		Key string `json:"key" validate:"min=8,max=128"`
	}
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		organizationID, err := idctx.Get(r, "organizationID")
		if err != nil {
			return
		}

		inviteID, err := idctx.Get(r, "inviteID")
		if err != nil {
			return
		}

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = organizationsService.AcceptInvite(ctx, organizationID, userID, inviteID, req.Key)
		if err != nil {
			switch err.(type) {
			case *organizations.ErrInviteInvalid:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *organizations.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{true})
	}
}

// InviteDeclinePost attempts to decline an invite.
func InviteDeclinePost(organizationsService organizations.Service) http.HandlerFunc {
	type request struct {
		Key string `json:"key" validate:"min=8,max=128"`
	}
	type response struct {
		Success bool `json:"success"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(auth.KeyUserID).(int64)
		if !ok {
			resp.ServerError(w, r, resp.UnknownError)
			return
		}

		organizationID, err := idctx.Get(r, "organizationID")
		if err != nil {
			return
		}

		inviteID, err := idctx.Get(r, "inviteID")
		if err != nil {
			return
		}

		req := request{}
		err = parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = organizationsService.DeclineInvite(ctx, organizationID, userID, inviteID, req.Key)
		if err != nil {
			switch err.(type) {
			case *organizations.ErrInviteInvalid:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *organizations.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		resp.OK(w, r, response{true})
	}
}

// OrganizationVolunteersGet gets all volunteers in all opportunities inside an organization.
func OrganizationVolunteersGet(opportunitiesService opportunities.Service, usersService users.Service) http.HandlerFunc {
	type OpportunitySummary struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	// OpportunityVolunteer represents a volunteer in an opportunity.
	type OpportunityVolunteer struct {
		OpportunitySummary OpportunitySummary `json:"opportunity"`
		models.OpportunityMembership
		users.UserProfile
	}
	// OpportunityPendingVolunteer represents a pending volunteer in an opportunity.
	type OpportunityPendingVolunteer struct {
		OpportunitySummary OpportunitySummary `json:"opportunity"`
		models.OpportunityMembershipRequest
		users.UserProfile
	}
	// OpportunityInvitedVolunteer represents a pending volunteer in an opportunity.
	type OpportunityInvitedVolunteer struct {
		OpportunitySummary OpportunitySummary `json:"opportunity"`
		EmailOnly          bool               `json:"emailOnly"`
		models.OpportunityMembershipInvite
		users.UserProfile
	}
	type response struct {
		Volunteers []OpportunityVolunteer        `json:"volunteers"`
		Pending    []OpportunityPendingVolunteer `json:"pending"`
		Invited    []OpportunityInvitedVolunteer `json:"invited"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		organizationID, err := idctx.Get(r, "organizationID")
		if err != nil {
			return
		}

		memberships, err := opportunitiesService.GetOrganizationOpportunityVolunteers(ctx, organizationID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrRequestNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		volunteers := []OpportunityVolunteer{}

		for _, membership := range memberships {
			profile, err := usersService.GetMinimalUserProfile(membership.UserID)
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

			volunteers = append(volunteers, OpportunityVolunteer{OpportunitySummary{membership.Opportunity.ID, membership.Opportunity.Title}, membership, *profile})
		}

		invited, err := opportunitiesService.GetOrganizationOpportunityInvitedVolunteers(ctx, organizationID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrRequestNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		invitedVolunteers := []OpportunityInvitedVolunteer{}

		for _, membership := range invited {
			profile, err := usersService.GetMinimalUserProfile(membership.InviteeID)
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

			if profile == nil {
				invitedVolunteers = append(invitedVolunteers, OpportunityInvitedVolunteer{OpportunitySummary{membership.Opportunity.ID, membership.Opportunity.Title}, false, membership, users.UserProfile{}})
				continue
			}
			invitedVolunteers = append(invitedVolunteers, OpportunityInvitedVolunteer{OpportunitySummary{membership.Opportunity.ID, membership.Opportunity.Title}, true, membership, *profile})
		}

		pending, err := opportunitiesService.GetOrganizationOpportunityRequestedVolunteers(ctx, organizationID)
		if err != nil {
			switch err.(type) {
			case *opportunities.ErrRequestNotFound:
				resp.NotFound(w, r, resp.APIError(err, nil))
			case *opportunities.ErrServerError:
				resp.ServerError(w, r, resp.APIError(err, nil))
			default:
				resp.ServerError(w, r, resp.UnknownError)
			}
			return
		}

		pendingVolunteers := []OpportunityPendingVolunteer{}

		for _, membership := range pending {
			profile, err := usersService.GetMinimalUserProfile(membership.VolunteerID)
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

			pendingVolunteers = append(pendingVolunteers, OpportunityPendingVolunteer{OpportunitySummary{membership.Opportunity.ID, membership.Opportunity.Title}, membership, *profile})
		}

		resp.OK(w, r, response{volunteers, pendingVolunteers, invitedVolunteers})
	}
}
