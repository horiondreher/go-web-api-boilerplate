package utils

import "golang.org/x/crypto/bcrypt"

type HashError struct {
	msg string
}

func (e *HashError) Error() string {
	return e.msg
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", &HashError{msg: "Error hashing password"}
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
