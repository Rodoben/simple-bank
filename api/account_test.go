package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateAccount(t *testing.T) {

	account := RandomAccount()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	args := db.CreateAccountParams{
		Owner:    account.Owner,
		Balance:  0,
		Currency: account.Currency,
	}

	store.EXPECT().CreateAccount(gomock.Any(), args).Times(1).Return(account, nil)
	server := NewServer(store)

	body1 := fmt.Sprintf(`{"owner": "%s", "currency":"%s"}`, account.Owner, account.Currency)
	recorder := httptest.NewRecorder()
	url := "/account"

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body1))
	assert.NoError(t, err)

	server.router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response db.Account
	json.Unmarshal(recorder.Body.Bytes(), &response)
	fmt.Println(recorder.Body.String())
	fmt.Println(account.Owner)
	assert.Equal(t, account.Owner, response.Owner)
	assert.Equal(t, account.Currency, response.Currency)
	assert.Equal(t, account.Balance, response.Balance)

}

func TestGetAccount(t *testing.T) {

	account := RandomAccount()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)
	store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/account/%d", account.ID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	server.router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

}

func Test_ListAccounts(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var accounts []db.Account
	for i := 0; i < 10; i++ {
		account := RandomAccount()
		accounts = append(accounts, account)
	}

	args := db.ListAccountsParams{
		Limit:  10,
		Offset: 5,
	}
	store := mockdb.NewMockStore(ctrl)

	store.EXPECT().ListAccounts(gomock.Any(), args).
		Times(1).
		Return(accounts, nil)
	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts?page_id=%v&page_size=%v", args.Limit, args.Offset)
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	server.router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
	fmt.Println(recorder.Body)
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
