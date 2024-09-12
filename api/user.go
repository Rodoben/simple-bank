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

type CreateUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Contact           string    `json:"contact"`
	Password          string    `json:"password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) CreateUserResponse {
	return CreateUserResponse{
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
	username, err := server.generateUsername(ctx, createUserRequest.Username)
	if err != nil {
		return
	}

	args := db.CreateuserParams{
		Username:       username,
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
