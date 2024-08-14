package sqlcSimpleBank

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore *Queries

func TestMain(m *testing.M) {

	dsn := "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable"
	connPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = New(connPool)
	m.Run()
}
