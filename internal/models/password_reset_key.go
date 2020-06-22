package models

import "time"

// PasswordResetKey represents a key that allows a user to reset their
// password.
type PasswordResetKey struct {
	Model
	UserID    int64 `json:"userID"`
	User      User
	Key       string    `json:"key"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// PasswordResetKeyRepository represents the repository for the
// PasswordResetKey.
type PasswordResetKeyRepository interface {
	// FindByID finds a single PasswordResetKey by ID.
	FindByID(id int64) (*PasswordResetKey, error)
	// FindByKey finds a single PasswordResetKey by Key.
	FindByKey(key string) (*PasswordResetKey, error)
	// Create creates a new PasswordResetKey.
	Create(passwordResetKey PasswordResetKey) error
	// Update updates a PasswordResetKey with the ID in the provided PasswordResetKey.
	Update(passwordResetKey PasswordResetKey) error
	// DeleteByID deletes a PasswordResetKey by ID.
	DeleteByID(id int64) error
}
