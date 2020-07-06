package oauth

import (
	"context"

	"github.com/joinimpact/api/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth "google.golang.org/api/oauth2/v2"
)

// GoogleClient is an interface for interacting with the Google Oauth APIs.
type GoogleClient interface {
	// GetToken gets an oauth2 token by code.
	GetToken(code string) (*oauth2.Token, error)
	// GetTokenFromAcessToken returns an oauth2 token from the access token string.'
	GetTokenFromAccessToken(accessToken string) *oauth2.Token
	// GetProfile gets a user's profile from an oauth token.
	GetProfile(token *oauth2.Token) (Profile, error)
}

// googleClient acts as the internal implementation of the GoogleClient interface.
type googleClient struct {
	config      *config.Config
	oauthConfig *oauth2.Config
}

// GoogleProfile represents a profile based on a Google account.
type GoogleProfile struct {
	FirstName string
	LastName  string
	Email     string
}

// NewGoogleClient creates and returns a new GoogleClient using the provided parameters.
func NewGoogleClient(config *config.Config) GoogleClient {
	return &googleClient{
		config,
		buildGoogleOauthConfig(config.GoogleClientID, config.GoogleClientSecret, config.GoogleCallbackURL),
	}
}

// buildGoogleOauthConfig creates and returns an oauth2 config from the provided client ID, secret, and callback URL.
func buildGoogleOauthConfig(clientID, clientSecret, callback string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  callback,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// GetToken gets an oauth2 token by code.
func (g *googleClient) GetToken(code string) (*oauth2.Token, error) {
	return g.oauthConfig.Exchange(context.Background(), code)
}

// GetTokenFromAcessToken returns an oauth2 token from the access token string.
func (g *googleClient) GetTokenFromAccessToken(accessToken string) *oauth2.Token {
	return &oauth2.Token{
		AccessToken: accessToken,
	}
}

// GetProfile gets a user's profile from an oauth token.
func (g *googleClient) GetProfile(token *oauth2.Token) (Profile, error) {
	oauthService, err := googleOauth.New(g.oauthConfig.Client(context.Background(), token))
	if err != nil {
		return nil, err
	}

	profile, err := oauthService.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return nil, err
	}

	return &GoogleProfile{
		FirstName: profile.GivenName,
		LastName:  profile.FamilyName,
		Email:     profile.Email,
	}, nil
}

// GetFirstName gets the first name of the profile.
func (p *GoogleProfile) GetFirstName() string {
	return p.FirstName
}

// GetLastName gets the last name of the profile.
func (p *GoogleProfile) GetLastName() string {
	return p.LastName
}

// GetEmail gets the email of the profile.
func (p *GoogleProfile) GetEmail() string {
	return p.Email
}
