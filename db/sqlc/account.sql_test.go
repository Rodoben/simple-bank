package db

import (
	"context"
	"fmt"
	"simple-bank/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func CreateRandomAccount(t *testing.T) (Account, CreateAccountParams, error) {
	var context = context.Background()

	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testStore.CreateAccount(context, args)

	if err != nil {
		return Account{}, CreateAccountParams{}, err
	}
	return account, args, nil
}

func Test_CreateAccount(t *testing.T) {
	result, args, err := CreateRandomAccount(t)

	assert.NoError(t, err, "they should not be any error!")
	assert.Equal(t, result.Owner, args.Owner)
	assert.Equal(t, result.Currency, args.Currency)
	assert.Equal(t, result.Currency, args.Currency)
	assert.NotZero(t, result.ID)
	assert.NotZero(t, result.CreatedAt)
}

func Test_GetAccount(t *testing.T) {

	var ctx = context.Background()

	tests := []struct {
		name              string
		createAccountArgs CreateAccountParams
	}{
		{
			name: "GetAccounts",
			createAccountArgs: CreateAccountParams{
				Owner:    util.RandomOwner(),
				Balance:  util.RandomMoney(),
				Currency: util.RandomCurrency(),
			},
		},
	}

	for _, test := range tests {

		account, err := testStore.CreateAccount(ctx, test.createAccountArgs)
		assert.NoError(t, err, "error")
		result, err := testStore.GetAccount(ctx, account.ID)

		assert.NoError(t, err, "error")
		assert.Equal(t, result.ID, account.ID)
		assert.Equal(t, result.Owner, account.Owner)
		assert.Equal(t, result.Currency, account.Currency)
		fmt.Println(result)
	}

}

// func Test_ListAccounts(t *testing.T) {

// 	var ctx = context.Background()
// 	args := []struct {
// 		params         ListAccountsParams
// 		expectedresult []Account
// 	}{
// 		{params: ListAccountsParams{
// 			Limit:  2,
// 			Offset: 0,
// 		},
// 			expectedresult: []Account{
// 				{ID: 1,
// 					Owner:    "ronald",
// 					Balance:  236,
// 					Currency: "INR",
// 					CreatedAt: time.Time{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
// 						InfinityModifier: 0,
// 						Valid:            true},
// 				},
// 				{ID: 2,
// 					Owner:    "ronald",
// 					Balance:  236,
// 					Currency: "INR",
// 					CreatedAt: pgtype.Timestamp{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
// 						InfinityModifier: 0,
// 						Valid:            true},
// 				},
// 			},
// 		},
// 		{params: ListAccountsParams{
// 			Limit:  2,
// 			Offset: 3,
// 		},
// 			expectedresult: []Account{
// 				{ID: 4,
// 					Owner:    "ronald",
// 					Balance:  236,
// 					Currency: "INR",
// 					CreatedAt: pgtype.Timestamp{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
// 						InfinityModifier: 0,
// 						Valid:            true},
// 				},
// 				{ID: 5,
// 					Owner:    "ronald",
// 					Balance:  236,
// 					Currency: "INR",
// 					CreatedAt: pgtype.Timestamp{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
// 						InfinityModifier: 0,
// 						Valid:            true},
// 				},
// 			},
// 		},
// 	}

// 	for _, v := range args {
// 		result, err := testStore.ListAccounts(ctx, v.params)
// 		assert.NoError(t, err, "error")
// 		assert.Equal(t, result, v.expectedresult)
// 	}

// }
