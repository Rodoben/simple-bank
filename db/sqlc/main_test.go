package db

import (
	"database/sql"
	"fmt"
	"log"
	"simple-bank/util"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	config, err := util.LoadConfig("./../..")
	if err != nil {
		log.Fatalf(err.Error())
	}

	connPool, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	fmt.Println(connPool)
	testStore = NewStore(connPool)
	m.Run()
}
