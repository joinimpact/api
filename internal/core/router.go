package core

import (
	"github.com/go-chi/chi"
	"github.com/joinimpact/api/internal/core/handlers/auth"
	"github.com/joinimpact/api/internal/core/handlers/users"
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

	return router
}
