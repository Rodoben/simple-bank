package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TransferTX(t *testing.T) {

	Account1, _, err := CreateRandomAccount(t)
	if err != nil {
		t.Fatalf(err.Error())
	}
	Account2, _, err := CreateRandomAccount(t)
	if err != nil {
		t.Fatalf(err.Error())
	}

	n := 5
	amount := int64(30)

	results := make(chan TransferTxResult)
	errors := make(chan error)
	for i := 0; i < n; i++ {
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: Account1.ID,
				ToAccountID:   Account2.ID,
				Amount:        int64(amount),
			})
			errors <- err
			results <- result

		}()

	}
	existed := make(map[int]bool)

	for i := 0; i < 5; i++ {

		err := <-errors
		assert.NoError(t, err)
		fmt.Println(err)
		result := <-results
		assert.NotEmpty(t, result)
		fmt.Println(result)

		transfer := result.Transfer

		assert.NotNil(t, transfer)
		assert.Equal(t, Account1.ID, transfer.FromAccountID)
		assert.Equal(t, Account2.ID, transfer.ToAccountID)
		assert.Equal(t, int64(amount), transfer.Amount)
		assert.NotEmpty(t, transfer.CreatedAt)
		assert.NotZero(t, transfer.ID)

		_, err = testStore.GetTransfer(context.Background(), transfer.ID)
		assert.NoError(t, err)

		fromEntry := result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, Account1.ID, fromEntry.AccountID)
		assert.Equal(t, int64(amount), fromEntry.Amount)
		assert.NotZero(t, fromEntry.ID)
		assert.NotEmpty(t, fromEntry.CreatedAt)

		toEntry := result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, Account2.ID, toEntry.AccountID)
		assert.Equal(t, int64(amount), toEntry.Amount)
		assert.NotZero(t, toEntry.ID)
		assert.NotEmpty(t, toEntry.CreatedAt)

		_, err = testStore.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, Account1.ID, fromAccount.ID)
		assert.NotEmpty(t, fromAccount.CreatedAt)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, Account2.ID, toAccount.ID)
		assert.NotEmpty(t, toAccount.CreatedAt)

		diff1 := Account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - Account2.Balance
		require.Equal(t, diff1, diff2)

		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // 1 * amount, 2 * amount, 3 * amount, ..., n * amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

}
