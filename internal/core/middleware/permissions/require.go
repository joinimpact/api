package permissions

import (
	"net/http"

	"github.com/joinimpact/api/pkg/resp"
	"github.com/joinimpact/api/pkg/scopes"
)

// Require requires an authentication level for a user to access an endpoint.
func Require(scope scopes.Scope) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userScope := scopes.ScopeFromContext(r.Context())
			if userScope < scope {
				resp.Forbidden(w, r, resp.Error(403, "you do not have sufficient permissions to this resource"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
