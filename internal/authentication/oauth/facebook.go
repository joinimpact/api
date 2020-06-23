package oauth

import (
	"context"
	"errors"

	fb "github.com/huandu/facebook/v2"
	"github.com/joinimpact/api/internal/config"
	"golang.org/x/oauth2"
)

// FacebookClient is an interface for interacting with the Facebook Oauth APIs.
type FacebookClient interface {
	// GetToken gets an oauth2 token by code.
	GetToken(code string) (*oauth2.Token, error)
	// GetProfile gets a user's profile from an oauth token.
	GetProfile(token *oauth2.Token) (Profile, error)
}

// facebookClient acts as the internal implementation of the FacebookClient interface.
type facebookClient struct {
	config      *config.Config
	oauthConfig *oauth2.Config
}

// FacebookProfile represents a profile based on a Facebook account.
type FacebookProfile struct {
	FirstName string
	LastName  string
	Email     string
}

// NewFacebookClient creates and returns a new FacebookClient using the provided parameters.
func NewFacebookClient(config *config.Config) FacebookClient {
	return &facebookClient{
		config,
		buildFacebookOauthConfig(config.FacebookAppID, config.FacebookAppSecret, config.FacebookCallbackURL),
	}
}

// buildFacebookOauthConfig creates and returns an oauth2 config from the provided client ID, secret, and callback URL.
func buildFacebookOauthConfig(clientID, clientSecret, callback string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  callback,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://graph.facebook.com/oauth/authorize",
			TokenURL: "https://graph.facebook.com/oauth/access_token",
		},
	}
}

// GetToken gets an oauth2 token by code.
func (f *facebookClient) GetToken(code string) (*oauth2.Token, error) {
	return f.oauthConfig.Exchange(context.Background(), code)
}

// GetProfile gets a user's profile from an oauth token.
func (f *facebookClient) GetProfile(token *oauth2.Token) (Profile, error) {
	res, err := fb.Get("/me", fb.Params{
		"fields":       "email,first_name,last_name",
		"access_token": token.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	firstName, ok := res["first_name"].(string)
	if !ok {
		return nil, errors.New("error with facebook api")
	}
	lastName, ok := res["last_name"].(string)
	if !ok {
		return nil, errors.New("error with facebook api")
	}
	email, ok := res["email"].(string)
	if !ok {
		return nil, errors.New("error with facebook api")
	}

	return &FacebookProfile{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}

// GetFirstName gets the first name of the profile.
func (p *FacebookProfile) GetFirstName() string {
	return p.FirstName
}

// GetLastName gets the last name of the profile.
func (p *FacebookProfile) GetLastName() string {
	return p.LastName
}

// GetEmail gets the email of the profile.
func (p *FacebookProfile) GetEmail() string {
	return p.Email
}
