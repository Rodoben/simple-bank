package api

import (
	"net/http"
	db "simple-bank/db/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username string `json:"user_name"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Contact  string `json:"contact"`
}

type CreateUserResponse struct {
	Username          string    `json:"user_name"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Contact           string    `json:"contact"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) CreateUserResponse {
	return CreateUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		Contact:           user.Contact,
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

	// generatepassword
	password, err := server.generateTempPassword(ctx)
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
		HashedPassword: password,
		FullName:       createUserRequest.FullName,
		Email:          createUserRequest.Email,
		Contact:        createUserRequest.Contact,
	}

	// call to create  a user record
	user, err := server.DbStore.Createuser(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := newUserResponse(user)

	ctx.JSON(http.StatusAccepted, resp)

}
