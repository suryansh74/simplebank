package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryansh74/simplebank/db"
	"github.com/suryansh74/simplebank/db/sqlc"
	"github.com/suryansh74/simplebank/token"
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

	fromAccount, valid := server.validAccount(context, req.FromAccountID, req.Currency)
	// TODO: currency should be same From and To account
	if !valid {
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		context.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, valid = server.validAccount(context, req.ToAccountID, req.Currency)
	if !valid {
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

func (server *Server) validAccount(context *gin.Context, accountID int64, currency sqlc.Currency) (sqlc.Account, bool) {
	// check wheater account is exist or not by id
	account, err := server.store.GetAccount(context, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	// check currency is same as given
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		context.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
