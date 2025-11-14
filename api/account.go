// Package api
package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryansh74/simplebank/db/sqlc"
)

type createAccountRequest struct {
	Owner    string        `json:"owner" binding:"required"`
	Currency sqlc.Currency `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(context *gin.Context) {
	// validating incoming req
	var req createAccountRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := sqlc.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(context, args)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(context *gin.Context) {
	var req getAccountRequest
	err := context.ShouldBindUri(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(context, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(context *gin.Context) {
	var req listAccountRequest
	err := context.ShouldBindQuery(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := sqlc.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(context, args)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusOK, accounts)
}
