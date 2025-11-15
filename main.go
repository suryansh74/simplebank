package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/simplebank/api"
	"github.com/suryansh74/simplebank/db"
	"github.com/suryansh74/simplebank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("unable to connect database:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	server.Start(config.ServerAddress)
}
