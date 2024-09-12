package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}
func CreateRandomUser(t *testing.T) {

	// Create a random user

	tests := []struct {
		name               string
		expectedStatusCode int
		payloadBody        string
		setupStub          func(store *mockdb.MockStore, createUserargs CreateUserRequest)
		checkResponse      func(t *testing.T, user db.User, body *bytes.Buffer)
	}{

		{
			name: "Bad payload",
			payloadBody: `{
                    "username": "roro
                    "full_name": "ronald benjamin",
                    "email":  "ronald.benjamin008@gmail.com",
                    "contact" : "9986398896"
                }`,
			expectedStatusCode: http.StatusBadRequest,
			setupStub: func(store *mockdb.MockStore, createuserargs CreateUserRequest) {

				store.EXPECT().GetUser(gomock.Any(), gomock.Any()).Times(0)

				store.EXPECT().Createuser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, user db.User, body *bytes.Buffer) {

			},
		},

		{
			name: "Internal Server error",
			payloadBody: `{
                    "username": "roro",
                    "full_name": "ronald benjamin",
                    "email":  "ronald.benjamin008@gmail.co",
                    "contact" : "9986398896"
                }`,
			expectedStatusCode: http.StatusInternalServerError,
			setupStub: func(store *mockdb.MockStore, createuserargs CreateUserRequest) {

				createUserRequestargs := db.CreateuserParams{
					Username: createuserargs.Username,
					Email:    createuserargs.Email,
					FullName: createuserargs.FullName,

					Contact: createuserargs.Contact,
				}

				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(createUserRequestargs.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)

				// Simulate internal server error during CreateUser call
				store.EXPECT().Createuser(gomock.Any(), gomock.AssignableToTypeOf(createUserRequestargs)).Times(1).Return(db.User{}, errors.New("internal server error"))
			},
			checkResponse: func(t *testing.T, user db.User, body *bytes.Buffer) {

				var createdUser db.User
				err := json.NewDecoder(body).Decode(&createdUser)
				assert.NoError(t, err)
				fmt.Println("check", user.Username, createdUser.Username)
				assert.Equal(t, user.Username, createdUser.Username)
				fmt.Println("check", user.Email, createdUser.Email)
				assert.Equal(t, user.Email, createdUser.Email)
				assert.Equal(t, user.FullName, createdUser.FullName)
				assert.Equal(t, user.Contact, createdUser.Contact)

			},
		},

		{
			name: "OK",
			payloadBody: `{
                    "username": "roro",
                    "full_name": "ronald benjamin",
                    "email":  "ronald.benjamin008@gmail.com",
                    "contact" : "9986398896"
                }`,
			expectedStatusCode: http.StatusAccepted,
			setupStub: func(store *mockdb.MockStore, createuserargs CreateUserRequest) {

				createUserRequestargs := db.CreateuserParams{
					Username: createuserargs.Username,
					Email:    createuserargs.Email,
					FullName: createuserargs.FullName,

					Contact: createuserargs.Contact,
				}

				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(createUserRequestargs.Username)).Times(1).Return(db.User{}, sql.ErrNoRows)
				user := RandomUser1(createUserRequestargs)
				store.EXPECT().Createuser(gomock.Any(), gomock.AssignableToTypeOf(createUserRequestargs)).Times(1).Return(user, nil)
			},
			checkResponse: func(t *testing.T, user db.User, body *bytes.Buffer) {

				var createdUser db.User
				err := json.NewDecoder(body).Decode(&createdUser)
				assert.NoError(t, err)
				fmt.Println("check", user.Username, createdUser.Username)
				assert.Equal(t, user.Username, createdUser.Username)
				fmt.Println("check", user.Email, createdUser.Email)
				assert.Equal(t, user.Email, createdUser.Email)
				assert.Equal(t, user.FullName, createdUser.FullName)
				fmt.Println("check", user.FullName, createdUser.FullName)
				assert.Equal(t, user.Contact, createdUser.Contact)

			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			var payload CreateUserRequest

			json.Unmarshal([]byte(test.payloadBody), &payload)

			test.setupStub(store, payload)

			server := NewServer(store)
			res := httptest.NewRecorder()
			url := "/user"
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(test.payloadBody)))
			assert.NoError(t, err)
			server.router.ServeHTTP(res, req)

			assert.Equal(t, test.expectedStatusCode, res.Code)

			user := db.User{}

			fmt.Println(res.Body)
			err = json.Unmarshal(res.Body.Bytes(), &user)
			fmt.Println("3", user.Username)
			assert.NoError(t, err)
			test.checkResponse(t, user, res.Body)

		})
	}

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

func RandomUser1(args db.CreateuserParams) db.User {

	return db.User{
		Username:          args.Username,
		HashedPassword:    args.HashedPassword,
		FullName:          args.FullName,
		Email:             args.Email,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		Contact:           args.Contact,
	}
}
