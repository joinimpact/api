package models

// UserTag defines a single user's area of interest.
type UserTag struct {
	Model
	UserID int64 `json:"-"`
	User   User  `json:"-"`
	TagID  int64 `json:"-"`
	Tag
}

// UserTagRepository represents a repository of UserTag.
type UserTagRepository interface {
	// FindByID finds a single entity by ID.
	FindByID(id int64) (*UserTag, error)
	// FindByUserID finds entities by UserID.
	FindByUserID(userID int64) ([]UserTag, error)
	// FindUserTagByID finds a single entity by UserID and tag ID.
	FindUserTagByID(userID int64, tagID int64) (*UserTag, error)
	// Create creates a new entity.
	Create(userTag UserTag) error
	// Update updates an entity with the ID in the provided entity.
	Update(userTag UserTag) error
	// DeleteByID deletes an entity by ID.
	DeleteByID(id int64) error
}
