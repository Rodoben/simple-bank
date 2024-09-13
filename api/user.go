package api

import (
	"fmt"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/util"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Contact           string    `json:"contact"`
	Password          string    `json:"password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) UserResponse {
	return UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		Contact:           user.Contact,
		Password:          user.HashedPassword,
		PasswordChangedAt: time.Now(),
		CreatedAt:         user.CreatedAt,
	}
}
func (server *Server) CreateUser(ctx *gin.Context) {

	var createUserRequest CreateUserRequest

	if err := ctx.ShouldBindJSON(&createUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fmt.Println("1", createUserRequest.Username)
	// generatepassword
	password, err := server.generateTempPassword()
	if err != nil {
		return
	}

	hashedpassword, err := util.Hashedpassword(password)
	if err != nil {
		return
	}

	// generateUsername

	args := db.CreateuserParams{
		Username:       createUserRequest.Username,
		HashedPassword: hashedpassword,
		FullName:       createUserRequest.FullName,
		Email:          createUserRequest.Email,
		Contact:        createUserRequest.Contact,
	}
	fmt.Println("2", args)
	// call to create  a user record
	user, err := server.DbStore.Createuser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user.HashedPassword = password
	resp := newUserResponse(user)

	ctx.JSON(http.StatusAccepted, resp)

}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken string       `json:"accesstoken"`
	User        UserResponse `json:"user"`
}

func (server *Server) LoginUser(ctx *gin.Context) {

	var loginUserRequest LoginRequest
	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.DbStore.GetUser(ctx, loginUserRequest.Username)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = util.ComparePassword(user.HashedPassword, loginUserRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := server.token.CreateToken(user.Username, server.Config.TokenExpiry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := LoginUserResponse{
		AccessToken: token,
		User:        newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, resp)

}
