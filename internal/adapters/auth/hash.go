package auth

import "golang.org/x/crypto/bcrypt"

func Hash_password(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(pass), err
}
