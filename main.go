package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/simplebank/api"
	"github.com/suryansh74/simplebank/db"
)

const (
	dbSource      = "postgresql://root:secret@192.168.29.20:5432/simple_bank?sslmode=disable"
	serverAddress = "localhost:8000"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("unable to connect database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	server.Start(serverAddress)
}
