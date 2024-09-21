package api

import (
	db "simple-bank/db/sqlc"
	"simple-bank/token"
	"simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Config  util.Config
	DbStore db.Store
	token   token.Maker
	router  *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	token, err := token.NewJwtMaker(config.AuthTokenKey)
	if err != nil {
		return nil, err
	}
	server := Server{
		Config:  config,
		DbStore: store,
		token:   token,
		router:  gin.Default(),
	}

	server.SetRoutes()

	return &server, nil

}

func (server *Server) SetRoutes() {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrrency)
	}

	server.router.POST("/user", server.CreateUser)
	server.router.POST("/user/login", server.LoginUser)

	authRoutes := server.router.Group("/").Use(authMiddleware(server.token))
	authRoutes.POST("/account", server.CreateAccount)
	authRoutes.GET("/account/:id", server.GetAccount)
	authRoutes.GET("/accounts", server.ListAccounts)
	authRoutes.POST("/transfer", server.CreateTransfer)

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
