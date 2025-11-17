package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryansh74/simplebank/db"
	"github.com/suryansh74/simplebank/db/sqlc"
)

type transferRequest struct {
	FromAccountID int64         `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64         `json:"to_account_id" binding:"required,min=1"`
	Amount        int64         `json:"amount" binding:"required,gt=0"`
	Currency      sqlc.Currency `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createTransfer(context *gin.Context) {
	// validating incoming req
	var req transferRequest
	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// TODO: currency should be same From and To account
	if !server.validAccount(context, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(context, req.ToAccountID, req.Currency) {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(context, args)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	context.JSON(http.StatusCreated, transfer)
}

func (server *Server) validAccount(context *gin.Context, accountID int64, currency sqlc.Currency) bool {
	// check wheater account is exist or not by id
	account, err := server.store.GetAccount(context, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	// check currency is same as given
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
