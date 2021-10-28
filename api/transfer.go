package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/xiahezzz/simplebank/db/sqlc"
)

type createTransferRequest struct {
	//binding:"required" 值后台会自动对输入进行验证
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	//测试输入参数是否有效
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorRequest(err))
		return
	}

	if !server.validAccount(ctx, req.FromAccountID, req.Currency) ||
		!server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}

	if !server.validBalance(ctx, req.FromAccountID, req.Amount) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorRequest(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorRequest(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorRequest(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] mismactch currency : %s vs %s ", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorRequest(err))
		return false
	}
	return true
}

func (server *Server) validBalance(ctx *gin.Context, accountID int64, amount int64) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorRequest(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorRequest(err))
		return false
	}
	if account.Balance < amount {
		err := fmt.Errorf("from Account [%d] does not have enough money: request:%d,now:%d", accountID, amount, account.Balance)
		ctx.JSON(http.StatusBadRequest, errorRequest(err))
		return false
	}
	return true
}
