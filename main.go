package main

import (
	"database/sql"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dsn      = "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable"
	dbDriver = "postgres"
)

func main() {

	conn, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalf("unable to connect to db %v", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(*store)
	server.Start(":8080")
}
