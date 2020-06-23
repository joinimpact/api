package models

// ThirdPartyIdentity contains information about a user's integrations with third-party services.
type ThirdPartyIdentity struct {
	Model
	UserID                 int64  `json:"userId"`
	User                   User   `foreignkey:"UserID"`
	ThirdPartyServiceName  string `json:"thirdPartyService"`
	ThirdPartyID           string
	ThirdPartyAccessToken  string
	ThirdPartyRefreshToken string
}

// ThirdPartyIdentityRepository represents the repository for the ThirdPartyIdentity.
type ThirdPartyIdentityRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*ThirdPartyIdentity, error)
	// FindByUserID finds a single entity by UserID.
	FindByUserID(userID int64) (*ThirdPartyIdentity, error)
	// FindUserIdentityByServiceName finds a single entity by user ID and service name.
	FindUserIdentityByServiceName(userID int64, serviceName string) (*ThirdPartyIdentity, error)
	// Create creates a new entity.
	Create(thirdPartyIdentity ThirdPartyIdentity) error
	// Update updates an entity with the ID in the provided parameter.
	Update(thirdPartyIdentity ThirdPartyIdentity) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
