package db

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dsn = "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable"

	//	dsn      = "postgresql:/username1:strongpassword@localhost:5432/simplebank?sslmode=disable"
	dbDriver = "postgres"
)

var testStore Store

func TestMain(m *testing.M) {

	connPool, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	fmt.Println(connPool)
	testStore = NewStore(connPool)
	m.Run()
}
