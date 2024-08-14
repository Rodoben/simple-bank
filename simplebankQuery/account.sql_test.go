package sqlcSimpleBank

import (
	"context"
	"fmt"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func Test_CreateAccount(t *testing.T) {

	var context = context.Background()

	args := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	result, err := testStore.CreateAccount(context, args)
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
		id              int64
		exprectedoutput Account
	}{
		{id: 1, exprectedoutput: Account{
			ID:       1,
			Owner:    "ronald",
			Balance:  100,
			Currency: "USD",
		}},
		{id: 2, exprectedoutput: Account{
			ID:       2,
			Owner:    "ronald",
			Balance:  100,
			Currency: "USD",
		}},
		{id: 16, exprectedoutput: Account{
			ID:       16,
			Owner:    "ycjiwz",
			Balance:  100,
			Currency: "USD",
		}},
	}

	for _, test := range tests {
		result, err := testStore.GetAccount(ctx, test.id)

		assert.NoError(t, err, "error")
		assert.Equal(t, result.ID, test.exprectedoutput.ID)
		assert.Equal(t, result.Owner, test.exprectedoutput.Owner)
		fmt.Println(result)
	}

}

func Test_ListAccounts(t *testing.T) {

	var ctx = context.Background()
	args := []struct {
		params         ListAccountsParams
		expectedresult []Account
	}{
		{params: ListAccountsParams{
			Limit:  2,
			Offset: 0,
		},
			expectedresult: []Account{
				{ID: 1,
					Owner:    "ronald",
					Balance:  236,
					Currency: "INR",
					CreatedAt: pgtype.Timestamp{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
						InfinityModifier: 0,
						Valid:            true},
				},
				{ID: 2,
					Owner:    "ronald",
					Balance:  236,
					Currency: "INR",
					CreatedAt: pgtype.Timestamp{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
						InfinityModifier: 0,
						Valid:            true},
				},
			},
		},
		{params: ListAccountsParams{
			Limit:  2,
			Offset: 3,
		},
			expectedresult: []Account{
				{ID: 4,
					Owner:    "ronald",
					Balance:  236,
					Currency: "INR",
					CreatedAt: pgtype.Timestamp{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
						InfinityModifier: 0,
						Valid:            true},
				},
				{ID: 5,
					Owner:    "ronald",
					Balance:  236,
					Currency: "INR",
					CreatedAt: pgtype.Timestamp{Time: time.Date(2024, time.August, 12, 7, 31, 16, 190834000, time.UTC),
						InfinityModifier: 0,
						Valid:            true},
				},
			},
		},
	}

	for _, v := range args {
		result, err := testStore.ListAccounts(ctx, v.params)
		assert.NoError(t, err, "error")
		assert.Equal(t, result, v.expectedresult)
	}

}
