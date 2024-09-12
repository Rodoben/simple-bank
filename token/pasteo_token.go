package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasteoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasteoMaker(secretkey string) (Maker, error) {
	if len(secretkey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)

	}

	maker := PasteoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(secretkey),
	}

	return &maker, nil
}

// CreateToken
func (maker *PasteoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := Newpayload(username, duration)
	if err != nil {
		return "", err
	}
	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", err
	}
	return token, err

}

// VerifyToken
func (maker *PasteoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)

	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	return payload, nil

}
