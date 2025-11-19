package utils

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)

	hashedPassword1, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	// check password matching
	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	// check for wrong password
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// check for 2 different hashedPassword
	hashedPassword2, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)

	log.Println("hashedPassword1:", hashedPassword1)
	log.Println("hashedPassword2:", hashedPassword2)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
