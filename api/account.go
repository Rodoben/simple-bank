package api

import (
	"errors"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/token"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type CreateRequestParams struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateAccount(ctx *gin.Context) {

	var req CreateRequestParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	args := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}
	account, err := server.DbStore.CreateAccount(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccount(ctx *gin.Context) {
	var id getAccountRequest
	if err := ctx.ShouldBindUri(&id); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.DbStore.GetAccount(ctx, id.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		errForbidden := errors.New("account doesn't belong to authenticated user")
		ctx.JSON(http.StatusForbidden, errorResponse(errForbidden))
		return
	}

	ctx.JSON(http.StatusOK, account)

}

type ListAccountsParams struct {
	PageId   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListAccounts(ctx *gin.Context) {

	var listAccountParams ListAccountsParams
	if err := ctx.ShouldBindQuery(&listAccountParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Offset: listAccountParams.PageSize,
		Limit:  listAccountParams.PageId,
	}
	accounts, err := server.DbStore.ListAccounts(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)

}
