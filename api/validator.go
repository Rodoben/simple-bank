package api

import (
	"database/sql"
	"net/http"
	"simple-bank/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validCurrrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSuppportedCurrency(currency)
	}
	return false
}

func (server *Server) generateUsername(ctx *gin.Context, username string) (string, error) {
	// Try to get the user by username
	_, err := server.DbStore.GetUser(ctx, username)

	// If the user does not exist, err should be a "not found" error
	if err != nil {
		// Check if the error indicates the user was not found
		if err == sql.ErrNoRows {
			// Append a random string to the username and return
			return username + util.RandomOwner(), nil
		}
		// Otherwise, return an error
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return "", err
	}

	// If the user exists, return an empty string (or you can handle it as needed)
	return "", nil
}

func (server *Server) generateTempPassword(ctx *gin.Context) (string, error) {

	password := util.RandomString(10)
	return password, nil

}
