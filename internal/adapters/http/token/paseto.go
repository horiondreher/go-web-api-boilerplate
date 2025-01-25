package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(email string, role string, duration time.Duration) (string, *Payload, *domainerr.DomainError) {
	payload, payloadErr := NewPayload(email, role, duration)
	if payloadErr != nil {
		return "", payload, payloadErr
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return token, nil, domainerr.NewInternalError(err)
	}

	return token, payload, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, *domainerr.DomainError) {
	payload := &Payload{}

	if maker == nil || maker.paseto == nil {
		return nil, domainerr.NewDomainError(http.StatusInternalServerError, domainerr.UnexpectedError, "internal server error", ErrInvalidInstance)
	}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, domainerr.NewDomainError(http.StatusUnauthorized, domainerr.UnauthorizedError, "invalid token", ErrInvalidToken)
	}

	validationErr := payload.Valid()
	if validationErr != nil {
		return nil, validationErr
	}

	return payload, nil
}
