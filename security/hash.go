package security

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes the given string with the specified cost.
func HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// VerifyPassword compares the password and hash and returns if they match.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
