package db

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateuser(t *testing.T) {
	CreateRandomUser(t)
}

func CreateRandomUser(t *testing.T) User {

	// create a random user
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	accountOwner := util.RandomOwner()

	fmt.Println("2", accountOwner)
	args := CreateuserParams{
		Username:       accountOwner,
		Email:          fmt.Sprintf("%v@gmail.com", accountOwner),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Contact:        util.RandomContactNumber(),
	}

	result, err := testStore.Createuser(ctx, args)
	assert.NoError(t, err)

	assert.Equal(t, result.Email, args.Email)

	assert.Equal(t, result.Username, args.Username)
	assert.Equal(t, result.FullName, args.FullName)
	assert.Equal(t, result.Contact, args.Contact)
	return result
}

func TestGetUser(t *testing.T) {

	tests := []struct {
		name           string
		User           User
		expectedStatus int
	}{
		{
			name:           "existing user",
			User:           CreateRandomUser(t),
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			Userdetails, err := testStore.GetUser(context.Background(), test.User.Username)
			assert.NoError(t, err)
			assert.NotNil(t, Userdetails)
			assert.Equal(t, Userdetails.Username, test.User.Username)
			assert.Equal(t, Userdetails.Email, test.User.Email)
			assert.Equal(t, Userdetails.FullName, test.User.FullName)
			reflect.DeepEqual(Userdetails, test.User)

		})
	}
}
