package auth

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/authentication"
	authm "github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/models"
	"github.com/joinimpact/api/internal/users"
	"github.com/joinimpact/api/pkg/location"
	"github.com/joinimpact/api/pkg/parse"
	"github.com/joinimpact/api/pkg/resp"
)

// Login attempts to login the user.
func Login(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Email    string `json:"email" validate:"email"`
			Password string `json:"password" validate:"min=8,max=512"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		tokenPair, err := service.Login(req.Email, req.Password)
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, err.Error()))
			return
		}

		// Update cookies.
		authm.SetAuthCookies(w, r, tokenPair)

		resp.OK(w, r, tokenPair)
	}
}

// ValidateEmail checks if a user's email is taken.
func ValidateEmail(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Email string `json:"email" validate:"email"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = service.CheckEmail(req.Email)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, err.Error()))
		}

		resp.OK(w, r, nil)
	}
}

// Register attempts to register the user.
func Register(service authentication.Service, usersService users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			FirstName   string                `json:"firstName" validate:"min=2,max=48"`
			LastName    string                `json:"lastName" validate:"min=2,max=48"`
			Email       string                `json:"email" validate:"email"`
			Password    string                `json:"password" validate:"min=8,max=512"`
			DateOfBirth time.Time             `json:"dateOfBirth"`
			Location    *location.Coordinates `json:"location"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		tokenPair, err := service.Register(models.User{
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Email:       req.Email,
			DateOfBirth: req.DateOfBirth,
		}, req.Password)
		if err != nil {
			resp.ServerError(w, r, resp.Error(500, err.Error()))
			return
		}

		if req.Location != nil {
			_ = usersService.UpdateUserLocation(tokenPair.UserID, req.Location)
		}

		// Update cookies.
		authm.SetAuthCookies(w, r, tokenPair)

		resp.OK(w, r, tokenPair)
	}
}

// RequestPasswordReset requests a password reset link to be emailed to a user.
func RequestPasswordReset(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Email string `json:"email" validate:"email"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = service.RequestPasswordReset(req.Email)
		if err != nil {
			resp.NotFound(w, r, resp.Error(404, err.Error()))
			return
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}

// VerifyPasswordReset verifies a password reset request by key and returns the user's first name and email for UI purposes.
func VerifyPasswordReset(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Key string `json:"key" validate:"min=2"`
		}{
			Key: chi.URLParam(r, "passwordResetKey"),
		}

		reset, err := service.CheckPasswordReset(req.Key)
		if err != nil {
			resp.NotFound(w, r, resp.Error(404, err.Error()))
			return
		}

		resp.OK(w, r, reset)
	}
}

// ResetPassword resets a user's password using a reset key.
func ResetPassword(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Key      string `json:"key"`
			Password string `json:"password" validate:"min=8,max=512"`
		}{
			Key: chi.URLParam(r, "passwordResetKey"),
		}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		err = service.ResetPassword(req.Key, req.Password)
		if err != nil {
			resp.NotFound(w, r, resp.Error(404, err.Error()))
			return
		}

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}

// GoogleOauth attempts to use Google to login through Oauth.
func GoogleOauth(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Token string `json:"token" validate:"min=8,max=1024"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		res, err := service.OauthLogin("google", req.Token)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, err.Error()))
			return
		}

		// Update cookies.
		authm.SetAuthCookies(w, r, res.TokenPair)

		resp.OK(w, r, res)
	}
}

// FacebookOauth attempts to use Facebook to login through Oauth.
func FacebookOauth(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Token string `json:"token" validate:"min=8,max=1024"`
		}{}
		err := parse.POST(w, r, &req)
		if err != nil {
			return
		}

		res, err := service.OauthLogin("facebook", req.Token)
		if err != nil {
			resp.BadRequest(w, r, resp.Error(400, err.Error()))
			return
		}

		// Update cookies.
		authm.SetAuthCookies(w, r, res.TokenPair)

		resp.OK(w, r, res)
	}
}

// Logout logs a user out by clearing auth and refresh token cookies.
func Logout(service authentication.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authm.ClearAuthCookies(w, r)

		resp.OK(w, r, map[string]bool{
			"success": true,
		})
	}
}
