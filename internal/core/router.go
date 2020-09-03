package core

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/core/handlers/auth"
	"github.com/joinimpact/api/internal/core/handlers/browse"
	"github.com/joinimpact/api/internal/core/handlers/conversations"
	"github.com/joinimpact/api/internal/core/handlers/events"
	"github.com/joinimpact/api/internal/core/handlers/hours"
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

		router.Route("/browse", func(r chi.Router) {
			r.Get("/", browse.Get(app.opportunitiesService))
			r.Post("/query", browse.QueryPost(app.opportunitiesService))
		})

		router.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				// For processing the userID param.
				r.Use(users.Middleware(app.authenticationService))
				r.Use(scopes.Middleware(users.ScopeProviderUsers()))

				r.Get("/", users.GetUserProfile(app.usersService))
				r.With(permissions.Require(scopes.ScopeOwner)).Delete("/", users.DeletePost(app.usersService))
				r.With(permissions.Require(scopes.ScopeOwner)).Patch("/", users.UpdateUserProfile(app.usersService))

				r.Get("/tags", users.GetUserTags(app.usersService))
				r.With(permissions.Require(scopes.ScopeOwner)).Post("/tags", users.PostUserTags(app.usersService))
				r.With(permissions.Require(scopes.ScopeOwner)).Delete("/tags/{tagID}", users.DeleteUserTag(app.usersService))

				r.With(permissions.Require(scopes.ScopeOwner)).Post("/profile-picture", users.UploadProfilePicture(app.usersService))

				r.With(permissions.Require(scopes.ScopeOwner)).Get("/organizations", organizations.GetUserOrganizations(app.organizationsService))
				r.Get("/opportunities", opportunities.GetByVolunteer(app.opportunitiesService))
				r.With(permissions.Require(scopes.ScopeOwner)).Get("/events", events.GetByVolunteer(app.eventsService))

				r.Route("/conversations", func(r chi.Router) {
					r.Use(permissions.Require(scopes.ScopeOwner))
					r.Get("/", conversations.GetByUser(app.conversationsService))

					r.Route("/{conversationID}", func(r chi.Router) {
						r.Use(idctx.Prepare("conversationID"))
						r.Use(scopes.Middleware(conversations.ScopeProviderConversations(app.conversationsService)))
						r.Use(permissions.Require(scopes.ScopeCollaborator))
						r.Get("/", conversations.Get(app.conversationsService, false))

						r.Route("/messages", func(r chi.Router) {
							r.Get("/", conversations.MessagesGet(app.conversationsService))
							r.Post("/", conversations.MessagesPost(app.conversationsService, false))
						})
					})
				})

				r.Route("/hours", func(r chi.Router) {
					r.Use(permissions.Require(scopes.ScopeOwner))
					r.Get("/", hours.GetByUser(app.hoursService))
				})
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

				r.Route("/invites", func(r chi.Router) {
					r.With(permissions.Require(scopes.ScopeAdmin)).Post("/", organizations.PostInvite(app.organizationsService))

					r.Route("/{inviteID}", func(r chi.Router) {
						r.Use(permissions.Require(scopes.ScopeAuthenticated))
						r.Use(idctx.Prepare("inviteID"))

						r.Post("/validate", organizations.InviteValidatePost(app.organizationsService))
						r.Post("/accept", organizations.InviteAcceptPost(app.organizationsService))
						r.Post("/decline", organizations.InviteDeclinePost(app.organizationsService))
					})
				})

				r.With(permissions.Require(scopes.ScopeManager)).Get("/members", organizations.MembersGet(app.organizationsService, app.usersService))
				r.With(permissions.Require(scopes.ScopeManager)).Get("/volunteers", organizations.OrganizationVolunteersGet(app.opportunitiesService, app.usersService))

				r.Route("/opportunities", func(r chi.Router) {
					r.With(permissions.Require(scopes.ScopeManager)).Post("/", opportunities.Post(app.opportunitiesService))
					r.With(permissions.Require(scopes.ScopeAuthenticated)).Get("/", opportunities.Get(app.opportunitiesService))
				})

				r.Route("/conversations", func(r chi.Router) {
					r.Use(permissions.Require(scopes.ScopeManager))

					r.Get("/", conversations.GetByOrganization(app.conversationsService))
					r.Route("/{conversationID}", func(r chi.Router) {
						r.Use(idctx.Prepare("conversationID"))
						r.Get("/", conversations.Get(app.conversationsService, true))

						r.Route("/messages", func(r chi.Router) {
							r.Get("/", conversations.MessagesGet(app.conversationsService))
							r.Post("/", conversations.MessagesPost(app.conversationsService, true))
						})
					})
				})

				r.Route("/hours", func(r chi.Router) {
					r.Route("/requests", func(r chi.Router) {
						r.Post("/", hours.OrganizationRequestsPost(app.hoursService, app.conversationsService))

						r.Route("/{requestID}", func(r chi.Router) {
							r.Use(idctx.Prepare("requestID"))
							r.With(permissions.Require(scopes.ScopeManager)).Post("/accept", hours.OrganizationRequestAcceptPost(app.hoursService, app.conversationsService))
							r.With(permissions.Require(scopes.ScopeManager)).Post("/decline", hours.OrganizationRequestDeclinePost(app.hoursService, app.conversationsService))
						})
					})
				})
			})

		})

		router.Route("/opportunities", func(r chi.Router) {
			r.Route("/{opportunityID}", func(r chi.Router) {
				r.Use(idctx.Prepare("opportunityID"))
				r.Use(scopes.Middleware(opportunities.ScopeProviderOpportunities(app.organizationsService, app.opportunitiesService)))

				r.Get("/", opportunities.GetOne(app.opportunitiesService))
				r.
					With(permissions.Require(scopes.ScopeManager)).
					Patch("/", opportunities.Patch(app.opportunitiesService))
				r.
					With(permissions.Require(scopes.ScopeManager)).
					Delete("/", opportunities.Delete(app.opportunitiesService))

				r.
					With(permissions.Require(scopes.ScopeManager)).
					Post("/profile-picture", opportunities.ProfilePicturePost(app.opportunitiesService))

				r.
					Post("/request", opportunities.RequestPost(app.opportunitiesService, app.conversationsService))

				r.Route("/volunteers", func(r chi.Router) {
					r.Use(permissions.Require(scopes.ScopeManager))
					r.Get("/", opportunities.VolunteersGet(app.opportunitiesService, app.usersService))

					r.Route("/{userID}", func(r chi.Router) {
						r.Use(idctx.Prepare("userID"))

						r.Post("/accept", opportunities.VolunteersAcceptPost(app.opportunitiesService, app.conversationsService))
						r.Post("/decline", opportunities.VolunteersDeclinePost(app.opportunitiesService))
					})
				})

				r.Get("/status", opportunities.StatusGet(app.opportunitiesService))

				r.Route("/invites", func(r chi.Router) {
					r.With(permissions.Require(scopes.ScopeAuthenticated)).Post("/", opportunities.InvitesPost(app.opportunitiesService))

					r.Route("/{inviteID}", func(r chi.Router) {
						r.Use(idctx.Prepare("inviteID"))

						r.Post("/validate", opportunities.InviteValidatePost(app.opportunitiesService))
						r.Post("/accept", opportunities.InviteAcceptPost(app.opportunitiesService))
						r.Post("/decline", opportunities.InviteDeclinePost(app.opportunitiesService))
					})
				})

				r.
					With(permissions.Require(scopes.ScopeManager)).
					Post("/publish", opportunities.PublishPost(app.opportunitiesService))
				r.
					With(permissions.Require(scopes.ScopeManager)).
					Post("/unpublish", opportunities.UnpublishPost(app.opportunitiesService))

				r.Route("/tags", func(r chi.Router) {
					r.Get("/", opportunities.TagsGet(app.opportunitiesService))
					r.
						With(permissions.Require(scopes.ScopeManager)).
						Post("/", opportunities.TagsPost(app.opportunitiesService))

					r.Route("/{tagID}", func(r chi.Router) {
						r.Use(idctx.Prepare("tagID"))

						r.
							With(permissions.Require(scopes.ScopeManager)).
							Delete("/", opportunities.TagsDelete(app.opportunitiesService))
					})
				})

				r.With(permissions.Require(scopes.ScopeManager)).Post("/events", events.Post(app.eventsService))
				r.With(permissions.Require(scopes.ScopeCollaborator)).Get("/events", events.GetByOpportunity(app.eventsService))
			})
		})

		router.Route("/events", func(r chi.Router) {
			r.Route("/{eventID}", func(r chi.Router) {
				r.Use(idctx.Prepare("eventID"))
				r.Use(scopes.Middleware(events.ScopeProviderEvents(app.eventsService, app.organizationsService, app.opportunitiesService)))
				r.Use(permissions.Require(scopes.ScopeCollaborator))

				r.Get("/", events.GetOne(app.eventsService))
				r.With(permissions.Require(scopes.ScopeManager)).Patch("/", events.Patch(app.eventsService))
				r.With(permissions.Require(scopes.ScopeManager)).Delete("/", events.Delete(app.eventsService))

				r.Route("/response", func(r chi.Router) {
					r.Get("/", events.ResponseGet(app.eventsService))
					r.Put("/", events.ResponsePut(app.eventsService))
				})

				r.Get("/responses", events.ResponsesGet(app.eventsService, app.usersService))
			})
		})
	})

	router.Route("/tags", func(r chi.Router) {
		r.Get("/", tags.GetTags(app.tagsService))
	})

	return router
}
