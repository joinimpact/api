package authentication

import "regexp"

const emailRegex = `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}`

// validateEmail checks whether or not a string is an email.
func validateEmail(email string) bool {
	// Compile the regex.
	regex, err := regexp.Compile(emailRegex)
	if err != nil {
		return false
	}

	// Run MatchString on the email and return the result.
	return regex.MatchString(email)
}

// validatePassword checks whether or not a password matches the set criteria.
func validatePassword(password string) bool {
	// TODO: better password validation.
	return len(password) >= 8 && len(password) < 512
}
