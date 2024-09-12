package token

import "time"

type Maker interface {
	//CreateToken
	CreateToken(username string, duration time.Duration) (token string, err error)
	//VerifyToken
	VerifyToken(token string) (*Payload, error)
}
