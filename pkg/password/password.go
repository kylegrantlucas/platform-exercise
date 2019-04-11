package password

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(plaintextPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePlaintextWithEncypted(plaintextPassword string, encryptedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(plaintextPassword))
	if err != nil {
		return false
	}

	return true
}
