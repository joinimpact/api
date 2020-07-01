package core

import (
	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/core/handlers/auth"
	"github.com/joinimpact/api/internal/core/handlers/opportunities"
	"github.com/joinimpact/api/internal/core/handlers/organizations"
	"github.com/joinimpact/api/internal/core/handlers/tags"
	"github.com/joinimpact/api/internal/core/handlers/users"
	authm "github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/pkg/idctx"
)

// Router assembles and returns a *chi.Mux router with all the API routes.
func (app *App) Router() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthcheck", healthcheckHandler)
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.Login(app.authenticationService))
		r.Post("/register", auth.Register(app.authenticationService))

		r.Route("/password-resets", func(r chi.Router) {
			r.Post("/", auth.RequestPasswordReset(app.authenticationService))
			r.Route("/{passwordResetKey}", func(r chi.Router) {
				r.Get("/", auth.VerifyPasswordReset(app.authenticationService))
				r.Post("/reset", auth.ResetPassword(app.authenticationService))
			})
		})

		r.Route("/oauth", func(r chi.Router) {
			r.Post("/google", auth.GoogleOauth(app.authenticationService))
			r.Post("/facebook", auth.FacebookOauth(app.authenticationService))
		})
	})

	router.Group(func(router chi.Router) {
		router.Use(authm.CookieMiddleware(app.authenticationService))
		router.Use(authm.AuthMiddleware(app.authenticationService))

		router.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				// For processing the userID param.
				r.Use(users.Middleware(app.authenticationService))
				r.Get("/", users.GetUserProfile(app.usersService))
				r.Patch("/", users.UpdateUserProfile(app.usersService))

				r.Get("/tags", users.GetUserTags(app.usersService))
				r.Post("/tags", users.PostUserTags(app.usersService))
				r.Delete("/tags/{tagID}", users.DeleteUserTag(app.usersService))

				r.Post("/profile-picture", users.UploadProfilePicture(app.usersService))
			})
		})

		router.Route("/organizations", func(r chi.Router) {

			r.Post("/", organizations.CreateOrganization(app.organizationsService))

			r.Route("/{organizationID}", func(r chi.Router) {
				r.Use(idctx.Prepare("organizationID"))

				// r.Get("/", users.GetUserProfile(app.usersService))
				// r.Patch("/", users.UpdateUserProfile(app.usersService))

				r.Get("/tags", organizations.GetOrganizationTags(app.organizationsService))
				r.Post("/tags", organizations.PostOrganizationTags(app.organizationsService))
				r.Delete("/tags/{tagID}", organizations.DeleteOrganizationTag(app.organizationsService))

				r.Post("/profile-picture", organizations.UploadProfilePicture(app.organizationsService))

				r.Post("/invite", organizations.PostInvite(app.organizationsService))
				r.Post("/invites", organizations.PostInvite(app.organizationsService))

				r.Route("/opportunities", func(r chi.Router) {
					r.Post("/", opportunities.Post(app.opportunitiesService))

					r.Route("/{opportunityID}", func(r chi.Router) {
						r.Use(idctx.Prepare("opportunityID"))

						r.Get("/", opportunities.Get(app.opportunitiesService))
						r.Patch("/", opportunities.Patch(app.opportunitiesService))
					})
				})
			})
		})

		router.Route("/tags", func(r chi.Router) {
			r.Get("/", tags.GetTags(app.tagsService))
		})
	})

	return router
}
