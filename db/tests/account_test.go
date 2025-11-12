package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suryansh74/simplebank/db/sqlc"
)

func TestAccount(t *testing.T) {
	arg := sqlc.CreateAccountParams{
		Owner:    "Tom",
		Balance:  100,
		Currency: sqlc.CurrencyUSD,
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
