package main

import (
	"database/sql"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"
	"simple-bank/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf(util.ErrorUnableToLoadConfig, err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("unable to connect to db %v", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("unable to start server %v", err)
	}
	server.Start(config.HTTPServer)
}
