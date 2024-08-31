package main

import (
	"database/sql"
	"fmt"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"
	"simple-bank/util"

	_ "github.com/lib/pq"
)

const (
	dsn      = "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable"
	dbDriver = "postgres"
)

func main() {

	config, err := util.LoadConfig()
	if err != nil {
		log.Fatalf("unable to load configuration %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("unable to connect to db %v", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(*store)
	fmt.Println("__", config.HTTPServer)
	server.Start(config.HTTPServer)
}
