package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suryansh74/simplebank/db/sqlc"
	"github.com/suryansh74/simplebank/utils"
)

func createRandomUser(t *testing.T) sqlc.User {
	hashedPassword, err := utils.HashedPassword(utils.RandomString(6))
	require.NoError(t, err)
	arg := sqlc.CreateUserParams{
		Username:       utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.True(t, user.PasswordChangedAt.Time.IsZero())

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)
	returnedUser, err := testQueries.GetUser(context.Background(), user.Username)

	require.NoError(t, err)
	require.NotEmpty(t, returnedUser)

	require.Equal(t, user.Username, returnedUser.Username)
	require.Equal(t, user.HashedPassword, returnedUser.HashedPassword)
	require.Equal(t, user.FullName, returnedUser.FullName)
	require.Equal(t, user.Email, returnedUser.Email)
	require.WithinDuration(t, user.PasswordChangedAt.Time, returnedUser.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, user.CreatedAt.Time, returnedUser.CreatedAt.Time, time.Second)
}
