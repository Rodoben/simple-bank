package token

import (
	"errors"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreatePasteoToken(t *testing.T) {

	maker, err := NewPasteoMaker(util.RandomString(32))
	assert.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, payload)
	assert.NotZero(t, payload.Id)
	assert.Equal(t, username, payload.Username)

	assert.WithinDuration(t, issuedAt, payload.CreatedAt, time.Second)
	assert.WithinDuration(t, expiresAt, payload.ExpiredAt, time.Second)
}

func TestPasteoExpiredToken(t *testing.T) {

	maker, err := NewPasteoMaker(util.RandomString(32))
	assert.NoError(t, err)
	username := util.RandomOwner()
	duration := -time.Minute
	token, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrorExpiredToken.Error())
	assert.Nil(t, payload)

}

func TestPasteoInvalidKeyToken(t *testing.T) {

	_, err := NewPasteoMaker(util.RandomString(10))
	assert.Error(t, err)
	assert.EqualError(t, err, errors.New("invalid key size: must be exactly 32 characters").Error())

}
