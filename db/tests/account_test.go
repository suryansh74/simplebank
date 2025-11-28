package tests

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/suryansh74/simplebank/db/sqlc"
	"github.com/suryansh74/simplebank/utils"
)

// createRandomAccount is does not have 'Test' prefix hence it don't count as unit test
func createRandomAccount(t *testing.T) sqlc.Account {
	user := createRandomUser(t)
	arg := sqlc.CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	returnedAccount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, returnedAccount)

	require.Equal(t, account.ID, returnedAccount.ID)
	require.Equal(t, account.Owner, returnedAccount.Owner)
	require.Equal(t, account.Balance, returnedAccount.Balance)
	require.Equal(t, account.Currency, returnedAccount.Currency)
	require.WithinDuration(t, account.CreatedAt.Time, returnedAccount.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	// create account
	account := createRandomAccount(t)
	// make update arguments
	args := sqlc.UpdateAccountParams{
		ID:      account.ID,
		Balance: utils.RandomMoney(),
	}
	// update account
	returnedAccount, err := testQueries.UpdateAccount(context.Background(), args)

	// check other account details with original one except balance
	require.NoError(t, err)
	require.NotEmpty(t, returnedAccount)

	require.Equal(t, account.ID, returnedAccount.ID)
	require.Equal(t, account.Owner, returnedAccount.Owner)
	require.Equal(t, account.Currency, returnedAccount.Currency)

	// check updated account balance with update argument
	require.Equal(t, args.Balance, returnedAccount.Balance)
	require.WithinDuration(t, account.CreatedAt.Time, returnedAccount.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	// error is expected since account1 is already deleted
	require.Error(t, err)
	require.Empty(t, account2)
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestListAccounts(t *testing.T) {
	var lastAccount sqlc.Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	args := sqlc.ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, a := range accounts {
		require.NotEmpty(t, a)
		require.Equal(t, lastAccount.Owner, a.Owner)
	}
}
