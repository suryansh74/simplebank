package tests

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/simplebank/db/sqlc"
)

var (
	testQueries *sqlc.Queries
	testDB      *pgxpool.Pool
)

func TestMain(m *testing.M) {
	var err error

	// Get database URL from environment variable, with fallback to local
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://root:secret@192.168.29.20:5432/simple_bank?sslmode=disable"
	}

	testDB, err = pgxpool.New(context.Background(), dbSource)
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
