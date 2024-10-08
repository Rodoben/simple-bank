package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	mockdb "simple-bank/db/mock"
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransfer(t *testing.T) {
	user1 := RandomUser()
	user2 := RandomUser()

	account1 := CreateRandomAccount(user1)

	account2 := CreateRandomAccount(user2)

	tests := []struct {
		name               string
		expectedStatusCode int
		payloadBody        string
		setupAuth          func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		setupStub          func(store *mockdb.MockStore, tranferargs transferRequest)
		checkResponse      func(t *testing.T, transferRecord db.TransferTxResult, body *bytes.Buffer)
	}{

		{
			name:               "Ok",
			expectedStatusCode: http.StatusAccepted,
			payloadBody:        `{"fromaccount_id": 2,"toaccount_id": 5,"currency": "INR","amount" :14}`,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.Username, time.Minute)
			},
			setupStub: func(store *mockdb.MockStore, tranferargs transferRequest) {

				args := db.TransferTxParams{
					FromAccountID: tranferargs.FromAccountId,
					ToAccountID:   tranferargs.ToAccountId,
					Amount:        tranferargs.Amount,
				}
				account1.Currency = tranferargs.Currency
				account2.Currency = tranferargs.Currency
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(args.FromAccountID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(args.ToAccountID)).Times(1).Return(account2, nil)

				transferResult := db.TransferTxResult{
					Transfer: db.Transfer{
						ID:            util.RandomInt(1, 1000),
						FromAccountID: args.FromAccountID,
						ToAccountID:   args.ToAccountID,
						Amount:        args.Amount,
					},
					FromAccount: db.Account{
						ID:    account1.ID,
						Owner: account1.Owner,
					},
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(args)).Times(1).Return(transferResult, nil)
			},
			checkResponse: func(t *testing.T, transferRecord db.TransferTxResult, body *bytes.Buffer) {
				var transferResult db.TransferTxResult
				err := json.NewDecoder(body).Decode(&transferResult)
				assert.NoError(t, err)
				assert.NotEmpty(t, transferResult.Transfer)
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			var payload transferRequest

			err := json.Unmarshal([]byte(test.payloadBody), &payload)
			assert.NoError(t, err)

			test.setupStub(store, payload)

			url := "/transfer"
			server := newTestServer(t, store)
			req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(test.payloadBody)))
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			test.setupAuth(t, req, server.token)

			server.router.ServeHTTP(recorder, req)

			assert.Equal(t, test.expectedStatusCode, recorder.Code)

			fmt.Println(recorder.Body.String())

			test.checkResponse(t, db.TransferTxResult{}, recorder.Body)

		})
	}
}

// func TestValidAccount(t *testing.T) {

// 	type args struct {
// 		ctx       *gin.Context
// 		accountId int64
// 		currency  string
// 	}
// 	tests := []struct {
// 		name           string
// 		expectedresult bool
// 		args           args
// 		setupStubs     func(store *mockdb.MockStore)
// 	}{

// 		{
// 			name:           "valid account",
// 			expectedresult: true,
// 			args: args{
// 				ctx:       &gin.Context{},
// 				accountId: 2,
// 				currency:  "USD",
// 			},
// 			setupStubs: func(store *mockdb.MockStore) {

// 				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{
// 					ID: 2,
// 				}, nil)
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {

// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			store := mockdb.NewMockStore(ctrl)
// 			test.setupStubs(store)
// 			server := NewServer(store)

// 			_, ok := server.validAccount(test.args.ctx, test.args.accountId, test.args.currency)
// 			assert.Equal(t, ok, test.expectedresult)
// 		})
// 	}
// }

func CreateRandomAccount(user db.User) db.Account {
	id := util.RandomInt(1, 1000)
	return db.Account{
		ID:        id,
		Owner:     user.Username,
		Balance:   util.RandomMoney(),
		Currency:  util.RandomCurrency(),
		CreatedAt: time.Now(),
	}

}
