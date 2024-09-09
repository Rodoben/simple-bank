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

	user, err := server.DbStore.GetUser(ctx, username)

	if err != nil {

		if err == sql.ErrNoRows {

			return username + util.RandomOwner(), nil
		}

		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return "", err
	}

	return user.Username, nil
}

func (server *Server) generateTempPassword() (string, error) {

	password := util.RandomString(10)
	return password, nil

}
