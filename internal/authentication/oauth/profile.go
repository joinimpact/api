package oauth

// Profile includes methods for interacting with an Oauth-generated profile.
type Profile interface {
	GetFirstName() string
	GetLastName() string
	GetEmail() string
}
