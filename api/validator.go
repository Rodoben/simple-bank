package api

import (
	"simple-bank/util"

	"github.com/go-playground/validator/v10"
)

var validCurrrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.IsSuppportedCurrency(currency)
	}
	return false
}

func (server *Server) generateTempPassword() (string, error) {

	password := util.RandomString(10)
	return password, nil

}
