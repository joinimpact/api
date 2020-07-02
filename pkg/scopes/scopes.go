// Package scopes encapsulates methods for conditionally hiding certain
// fields based on scopes, which are increasing levels of permissions a
// certain user holds.
package scopes

// Scope represents the scope of what a user can view.
type Scope int

// Scope constants
const (
	ScopeUnauthenticated Scope = iota
	ScopeAuthenticated   Scope = iota
	ScopeCollaborator    Scope = iota
	ScopeManager         Scope = iota
	ScopeAdmin           Scope = iota
	ScopeOwner           Scope = iota
)

// stringToScope converts a string (ex: "admin") to a Scope type.
func stringToScope(scope string) Scope {
	switch scope {
	case "unauthenticated":
		return ScopeUnauthenticated
	case "authenticated":
		return ScopeAuthenticated
	case "collaborator":
		return ScopeCollaborator
	case "manager":
		return ScopeManager
	case "admin":
		return ScopeAdmin
	case "owner":
		return ScopeOwner
	}

	// Fallback to unauthenticated.
	return ScopeUnauthenticated
}
