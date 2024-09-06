package api

import (
	"fmt"
	"net/http"
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountId int64  `json:"fromaccount_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"toaccount_id" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,currency"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var transfer transferRequest
	if err := ctx.ShouldBindJSON(&transfer); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, ok := server.validAccount(ctx, transfer.FromAccountId, transfer.Currency)
	if !ok {
		return
	}

	_, ok = server.validAccount(ctx, transfer.ToAccountId, transfer.Currency)
	if !ok {
		return
	}
	// TODO: implement transfer logic
	args := db.TransferTxParams{
		FromAccountID: transfer.FromAccountId,
		ToAccountID:   transfer.ToAccountId,
		Amount:        transfer.Amount,
	}

	result, err := server.DbStore.TransferTx(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {

	account, err := server.DbStore.GetAccount(ctx, accountId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
