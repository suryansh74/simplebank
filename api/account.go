// Package api
package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/suryansh74/simplebank/db/sqlc"
	"github.com/suryansh74/simplebank/token"
)

type createAccountRequest struct {
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

	authPayload := context.MustGet(authroizatoinPayloadKey).(*token.Payload)
	args := sqlc.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(context, args)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23503": // foreign_key_violation
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			case "23505": // unique_violation
				context.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
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

	authPayload := context.MustGet(authroizatoinPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to authenticated user")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
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

	authPayload := context.MustGet(authroizatoinPayloadKey).(*token.Payload)
	args := sqlc.ListAccountsParams{
		Owner:  authPayload.Username,
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
