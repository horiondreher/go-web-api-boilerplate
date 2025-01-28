package utils

import (
	"net/http"

	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"golang.org/x/crypto/bcrypt"
)

type HashError struct {
	msg string
}

func (e *HashError) Error() string {
	return e.msg
}

func HashPassword(password string) (string, *domainerr.DomainError) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", domainerr.NewDomainError(http.StatusInternalServerError, domainerr.InternalError, err.Error(), err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) *domainerr.DomainError {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return domainerr.MatchHashError(err)
	}

	return nil
}
