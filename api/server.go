package api

import (
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	DbStore db.Store
	router  *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := Server{
		DbStore: store,
		router:  gin.Default(),
	}

	server.SetRoutes()

	return &server

}

func (server *Server) SetRoutes() {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrrency)
	}

	server.router.POST("/account", server.CreateAccount)
	server.router.GET("/account/:id", server.GetAccount)
	server.router.GET("/accounts", server.ListAccounts)

	server.router.POST("/transfer", server.CreateTransfer)

	server.router.POST("/user", server.CreateUser)

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
