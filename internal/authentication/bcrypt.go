package authentication

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptCostFactor represents how hard bcrypt has to work to hash the password.
const BcryptCostFactor = 12

// compareHashAndPassword takes a plaintext password and compares it to a hashed password.
// If the password matches the hash, it returns true.
func compareHashAndPassword(password, hash string) (bool, error) {
	// Compare the hash and password.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil, nil
}

// generateFromPassword generates a bcrypt hash from a password string.
func generateFromPassword(password string) (string, error) {
	// Generate a hash from the password and the bcrypt cost factor defined as a constant.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCostFactor)
	if err != nil {
		return "", err
	}

	// Convert the []byte to a string.
	return string(hash), nil
}
