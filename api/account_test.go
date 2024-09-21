package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/util"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func requireBodymatcherAccount(t *testing.T, account db.Account, body *bytes.Buffer) {
	response, err := io.ReadAll(body)

	assert.NoError(t, err)

	var gotAccount db.Account
	json.Unmarshal(response, &gotAccount)
	reflect.DeepEqual(account, gotAccount)

}

func Test_CreateAccount(t *testing.T) {
	account := RandomAccount()
	args := db.CreateAccountParams{
		Owner:    account.Owner,
		Balance:  0,
		Currency: account.Currency,
	}
	tests := []struct {
		name               string
		payloadBody        string
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		expectedStatuscode int
		buildStubs         func(store *mockdb.MockStore)
		checkResponse      func(t *testing.T, account db.Account, body *bytes.Buffer)
	}{

		{
			name:               "BadPayload",
			payloadBody:        fmt.Sprintf(`{"owner: "%s", "currency":"%s"}`, account.Owner, account.Currency),
			expectedStatuscode: http.StatusBadRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), args).Times(0).Return(db.Account{}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, account db.Account, body *bytes.Buffer) {

			},
		},

		{
			name:               "Internal Server error",
			payloadBody:        fmt.Sprintf(`{"owner": "%s", "currency":"%s"}`, account.Owner, account.Currency),
			expectedStatuscode: http.StatusInternalServerError,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), args).Times(1).Return(db.Account{}, errors.New("internal server error"))
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, account db.Account, body *bytes.Buffer) {
				requireBodymatcherAccount(t, account, body)
			},
		},
		{
			name:               "OK",
			payloadBody:        fmt.Sprintf(`{"owner": "%s", "currency":"%s"}`, account.Owner, account.Currency),
			expectedStatuscode: http.StatusCreated,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), args).Times(1).Return(account, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, account db.Account, body *bytes.Buffer) {
				requireBodymatcherAccount(t, account, body)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			server := newTestServer(t, store)

			test.buildStubs(store)

			recorder := httptest.NewRecorder()
			url := "/account"

			req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(test.payloadBody))
			assert.NoError(t, err)

			test.setupAuth(t, req, server.token)

			server.router.ServeHTTP(recorder, req)
			assert.Equal(t, test.expectedStatuscode, recorder.Code)

			test.checkResponse(t, account, recorder.Body)

		})

	}

}

func TestGetAccount(t *testing.T) {
	account := RandomAccount()
	tests := []struct {
		name               string
		accountId          int64
		buildStubs         func(store *mockdb.MockStore)
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		expectedStatuscode int
		checkResponse      func(t *testing.T, account db.Account, body *bytes.Buffer)
	}{

		{
			name:               "BadRequest",
			accountId:          0,
			expectedStatuscode: http.StatusBadRequest,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), 0).Times(0).Return(db.Account{}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, account db.Account, body *bytes.Buffer) {
				requireBodymatcherAccount(t, account, body)
			},
		},

		{
			name:               "Internal Server Error",
			accountId:          account.ID,
			expectedStatuscode: http.StatusInternalServerError,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(db.Account{}, errors.New("internal server error"))
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, account db.Account, body *bytes.Buffer) {
				requireBodymatcherAccount(t, account, body)
			},
		},
		{
			name:               "OK",
			accountId:          account.ID,
			expectedStatuscode: http.StatusOK,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), account.ID).Times(1).Return(account, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, account db.Account, body *bytes.Buffer) {
				requireBodymatcherAccount(t, account, body)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			test.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/%d", test.accountId)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)
			test.setupAuth(t, req, server.token)
			server.router.ServeHTTP(recorder, req)
			assert.Equal(t, test.expectedStatuscode, recorder.Code)
			test.checkResponse(t, account, recorder.Body)
		})
	}

}

func Test_ListAccounts(t *testing.T) {
	var accounts []db.Account
	var account db.Account
	for i := 0; i < 10; i++ {
		account := RandomAccount()
		accounts = append(accounts, account)
	}
	tests := []struct {
		name               string
		args               db.ListAccountsParams
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs         func(store *mockdb.MockStore, args db.ListAccountsParams)
		expectedStatuscode int
		checkResponse      func(t *testing.T, accounts []db.Account, body *bytes.Buffer)
	}{

		{
			name: "Bad Request",
			args: db.ListAccountsParams{
				Limit:  0,
				Offset: 0,
			},
			expectedStatuscode: http.StatusBadRequest,
			buildStubs: func(store *mockdb.MockStore, args db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), args).Times(0).Return([]db.Account{}, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, accounts []db.Account, body *bytes.Buffer) {
				for _, account := range accounts {
					requireBodymatcherAccount(t, account, body)
				}

			},
		},

		{
			name: "Internal server error",
			args: db.ListAccountsParams{
				Limit:  10,
				Offset: 10,
			},
			expectedStatuscode: http.StatusInternalServerError,
			buildStubs: func(store *mockdb.MockStore, args db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), args).Times(1).Return([]db.Account{}, errors.New("internal server error"))
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, accounts []db.Account, body *bytes.Buffer) {
				for _, account := range accounts {
					requireBodymatcherAccount(t, account, body)
				}

			},
		},
		{
			name: "OK",
			args: db.ListAccountsParams{
				Limit:  10,
				Offset: 5,
			},
			expectedStatuscode: http.StatusOK,
			buildStubs: func(store *mockdb.MockStore, args db.ListAccountsParams) {
				store.EXPECT().ListAccounts(gomock.Any(), args).Times(1).Return(accounts, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, account.Owner, time.Minute)
			},
			checkResponse: func(t *testing.T, accounts []db.Account, body *bytes.Buffer) {
				for _, account := range accounts {
					requireBodymatcherAccount(t, account, body)
				}

			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			test.buildStubs(store, test.args)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?page_id=%v&page_size=%v", test.args.Limit, test.args.Offset)
			fmt.Println(url)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)
			test.setupAuth(t, req, server.token)
			server.router.ServeHTTP(recorder, req)
			assert.Equal(t, test.expectedStatuscode, recorder.Code)
			//	fmt.Println(recorder.Body)
			test.checkResponse(t, accounts, recorder.Body)
		})
	}

}

func RandomAccount() db.Account {
	return db.Account{
		ID:        util.RandomInt(1, 100),
		Owner:     util.RandomOwner(),
		Balance:   util.RandomMoney(),
		Currency:  util.RandomCurrency(),
		CreatedAt: time.Now(),
	}
}
