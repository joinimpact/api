package authentication

import "testing"

const samplePassword = "123Password!"
const wrongPassword = "123Passwodd!"

// TestBcrypt tests bcrypt functionality with two passwords that match.
func TestBcryptMatch(t *testing.T) {
	// Generate a hash from the password.
	hash, err := generateFromPassword(samplePassword)
	if err != nil {
		t.Fatal("error generating password hash:", err)
	}

	// Compare the password with the original hash.
	ok, err := compareHashAndPassword(samplePassword, hash)
	if err != nil {
		t.Fatal("error comparing hash and password:", err)
	}
	if !ok {
		t.Fatal("incorrect reporting - password reported as invalid when it matches:", err)
	}
}

// TestBcrypt tests bcrypt functionality with two passwords that don't match.
func TestBcryptFail(t *testing.T) {
	// Generate a hash from the password.
	hash, err := generateFromPassword(samplePassword)
	if err != nil {
		t.Fatal("error generating password hash:", err)
	}

	// Compare the password with the original hash.
	ok, err := compareHashAndPassword(wrongPassword, hash)
	if err != nil {
		t.Fatal("error comparing hash and password:", err)
	}
	if ok {
		t.Fatal("incorrect reporting - password reported as valid when it does not match:", err)
	}
}
