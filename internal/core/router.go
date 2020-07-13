package core

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/core/handlers/auth"
	"github.com/joinimpact/api/internal/core/handlers/opportunities"
	"github.com/joinimpact/api/internal/core/handlers/organizations"
	"github.com/joinimpact/api/internal/core/handlers/tags"
	"github.com/joinimpact/api/internal/core/handlers/users"
	authm "github.com/joinimpact/api/internal/core/middleware/auth"
	"github.com/joinimpact/api/internal/core/middleware/db"
	"github.com/joinimpact/api/internal/core/middleware/permissions"
	"github.com/joinimpact/api/pkg/idctx"
	"github.com/joinimpact/api/pkg/scopes"
)

// Router assembles and returns a *chi.Mux router with all the API routes.
func (app *App) Router() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthcheck", healthcheckHandler)
	router.Route("/auth", func(r chi.Router) {
		r.Post("/login", auth.Login(app.authenticationService))
		r.Post("/validate-email", auth.ValidateEmail(app.authenticationService))
		r.Post("/register", auth.Register(app.authenticationService, app.usersService))
		r.Post("/logout", auth.Logout(app.authenticationService))

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
		router.Use(scopes.Middleware(func(ctx context.Context) scopes.Scope {
			return scopes.ScopeAuthenticated
		}))
		// Gets limit and other database query parameters from the URL.
		router.Use(db.ContextMiddleware())

		router.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				// For processing the userID param.
				r.Use(users.Middleware(app.authenticationService))
				r.Use(scopes.Middleware(users.ScopeProviderUsers()))

				r.Get("/", users.GetUserProfile(app.usersService))
				r.With(permissions.Require(scopes.ScopeOwner)).Patch("/", users.UpdateUserProfile(app.usersService))

				r.Get("/tags", users.GetUserTags(app.usersService))
				r.With(permissions.Require(scopes.ScopeOwner)).Post("/tags", users.PostUserTags(app.usersService))
				r.With(permissions.Require(scopes.ScopeOwner)).Delete("/tags/{tagID}", users.DeleteUserTag(app.usersService))

				r.With(permissions.Require(scopes.ScopeOwner)).Post("/profile-picture", users.UploadProfilePicture(app.usersService))

				r.With(permissions.Require(scopes.ScopeOwner)).Get("/organizations", organizations.GetUserOrganizations(app.organizationsService))
			})
		})

		router.Route("/organizations", func(r chi.Router) {
			r.Post("/", organizations.CreateOrganization(app.organizationsService))

			r.Route("/{organizationID}", func(r chi.Router) {
				r.Use(idctx.Prepare("organizationID"))
				r.Use(scopes.Middleware(organizations.ScopeProviderOrganizations(app.organizationsService)))

				r.Get("/", organizations.GetOrganizationProfile(app.organizationsService))
				r.With(permissions.Require(scopes.ScopeAdmin)).Patch("/", organizations.UpdateOrganizationProfile(app.organizationsService))
				r.With(permissions.Require(scopes.ScopeAdmin)).Delete("/", organizations.DeleteOrganization(app.organizationsService))

				r.Get("/tags", organizations.GetOrganizationTags(app.organizationsService))
				r.With(permissions.Require(scopes.ScopeAdmin)).Post("/tags", organizations.PostOrganizationTags(app.organizationsService))
				r.With(permissions.Require(scopes.ScopeAdmin)).Delete("/tags/{tagID}", organizations.DeleteOrganizationTag(app.organizationsService))

				r.With(permissions.Require(scopes.ScopeAdmin)).Post("/profile-picture", organizations.UploadProfilePicture(app.organizationsService))
				r.With(permissions.Require(scopes.ScopeAdmin)).Post("/invite", organizations.PostInvite(app.organizationsService))
				r.With(permissions.Require(scopes.ScopeAdmin)).Post("/invites", organizations.PostInvite(app.organizationsService))

				r.Route("/opportunities", func(r chi.Router) {
					r.With(permissions.Require(scopes.ScopeManager)).Post("/", opportunities.Post(app.opportunitiesService))
					r.With(permissions.Require(scopes.ScopeAuthenticated)).Get("/", opportunities.Get(app.opportunitiesService))
				})
			})
		})

		router.Route("/opportunities", func(r chi.Router) {
			r.Route("/{opportunityID}", func(r chi.Router) {
				r.Use(idctx.Prepare("opportunityID"))

				r.Get("/", opportunities.GetOne(app.opportunitiesService))
				r.
					With(permissions.Require(scopes.ScopeAuthenticated)).
					Patch("/", opportunities.Patch(app.opportunitiesService))
				r.
					With(permissions.Require(scopes.ScopeAuthenticated)).
					Delete("/", opportunities.Delete(app.opportunitiesService))

				r.
					With(permissions.Require(scopes.ScopeAuthenticated)).
					Post("/profile-picture", opportunities.ProfilePicturePost(app.opportunitiesService))

				r.
					Post("/request", opportunities.RequestPost(app.opportunitiesService, app.conversationsService))

				r.Route("/volunteers", func(r chi.Router) {
					r.
						With(permissions.Require(scopes.ScopeAuthenticated)).
						Get("/", opportunities.VolunteersGet(app.opportunitiesService, app.usersService))
				})

				r.
					With(permissions.Require(scopes.ScopeAuthenticated)).
					Post("/publish", opportunities.PublishPost(app.opportunitiesService))
				r.
					With(permissions.Require(scopes.ScopeAuthenticated)).
					Post("/unpublish", opportunities.UnpublishPost(app.opportunitiesService))

				r.Route("/tags", func(r chi.Router) {
					r.Get("/", opportunities.TagsGet(app.opportunitiesService))
					r.
						With(permissions.Require(scopes.ScopeAuthenticated)).
						Post("/", opportunities.TagsPost(app.opportunitiesService))

					r.Route("/{tagID}", func(r chi.Router) {
						r.Use(idctx.Prepare("tagID"))

						r.
							With(permissions.Require(scopes.ScopeAuthenticated)).
							Delete("/", opportunities.TagsDelete(app.opportunitiesService))
					})
				})
			})
		})
	})

	router.Route("/tags", func(r chi.Router) {
		r.Get("/", tags.GetTags(app.tagsService))
	})

	return router
}
