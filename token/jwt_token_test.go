package token

import (
	"simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateJwtToken(t *testing.T) {

	maker, err := NewJwtMaker(util.RandomString(32))
	assert.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotZero(t, payload.Id)
	assert.Equal(t, username, payload.Username)

	assert.WithinDuration(t, issuedAt, payload.CreatedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredToken(t *testing.T) {

	maker, err := NewJwtMaker(util.RandomString(32))
	assert.NoError(t, err)
	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrorInvalidToken.Error())
	assert.Nil(t, payload)

}
