package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), DefaultBcryptCost)
	if err != nil {
		return "", err
	}

	return string(pass), err
}

// compare the password to its hash
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
