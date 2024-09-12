package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	minimumCost = 32
)

type JwtMaker struct {
	secretkey string
}

func NewJwtMaker(secretKey string) (Maker, error) {

	if len(secretKey) < minimumCost {
		return nil, nil

	}

	return &JwtMaker{secretkey: secretKey}, nil

}

func (maker *JwtMaker) CreateToken(username string, duration time.Duration) (string, error) {
	//   initialize an new payload
	fmt.Println("1")
	payload, err := Newpayload(username, duration)
	if err != nil {
		return "", nil
	}
	fmt.Println("2")

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// create a jwt token with claims
	fmt.Println("3")

	fmt.Println("4")
	// create a token with signedString
	return jwtToken.SignedString([]byte(maker.secretkey))
}

func (maker *JwtMaker) VerifyToken(token string) (*Payload, error) {
	Keyfunc := func(t *jwt.Token) (any, error) {

		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretkey), nil

	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, Keyfunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrorInvalidToken) {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}

	Payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return Payload, nil
}
