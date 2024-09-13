package api

import (
	"os"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func newTestServer(t *testing.T, store db.Store) *Server {

	config := util.Config{
		AuthTokenKey: util.RandomString(32),
		TokenExpiry:  time.Minute,
	}
	server, err := NewServer(config, store)
	assert.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
