package authentication

import (
	"math/rand"
	"time"
)

// PasswordResetValidation contains fields about a password reset.
type PasswordResetValidation struct {
	FirstName string `json:"firstName"`
	Email     string `json:"email"`
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// TODO: move global variable to another solution.
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// stringWithCharset generates a random string using a charset.
func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// generatePasswordResetKey generates and returns a random string to use as a key.
func generatePasswordResetKey() string {
	return stringWithCharset(24, charset)
}
