package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/simplebank/db/sqlc"
	"github.com/suryansh74/simplebank/utils"
)

var (
	testQueries *sqlc.Queries
	testDB      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../../")

	// Get database URL from environment variable, with fallback to local
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("unable to connect database:", err)
	}

	testQueries = sqlc.New(testDB)

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()

	// Exit with test result code
	os.Exit(code)
}
