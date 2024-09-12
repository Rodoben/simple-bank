package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorInvalidToken = errors.New("token is invalid")
	ErrorExpiredToken = errors.New("token has expired")
)

type Payload struct {
	Id        uuid.UUID
	Username  string
	createdAt time.Time
	ExpiredAt time.Time
}

func Newpayload(username string, duration time.Duration) (*Payload, error) {

	return &Payload{
		Id:        uuid.New(),
		Username:  username,
		createdAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil

}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}

	return nil

}
