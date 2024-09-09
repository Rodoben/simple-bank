package api

import (
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"testing"
	"time"
)

// func CreateRandomUser(t *testing.T) db.User {

// 	// Create a random user

// 	user := RandomUser()

// 	test := []struct {
// 		name               string
// 		expectedStatusCode int
// 		setupStub          func(store *mockdb.MockStore)
// 	}{
// 		{},
// 	}

// }

func TestCreateUser(t *testing.T) {
	//CreateRandomUser(t)
}

func RandomUser() db.User {
	username := util.RandomOwner()
	return db.User{
		Username:          username,
		HashedPassword:    "secret",
		FullName:          username,
		Email:             util.RandomEmail(),
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		Contact:           util.RandomContactNumber(),
	}
}
