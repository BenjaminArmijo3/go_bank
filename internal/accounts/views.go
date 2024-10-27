package accounts

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/BenjaminArmijo3/bank/internal/auth"
	"github.com/BenjaminArmijo3/bank/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	UserID  int32   `json:"user_id" binding:"required"`
	Balance float64 `json:"balance" binding:"required"`
}

type createAccountResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *AccountApp) createAccount(c *gin.Context) {
	var account createAccountRequest

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userDb, err := a.server.Store.Queries.GetUserById(context.Background(), int64(account.UserID))
	if err == sql.ErrNoRows {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exists"})
		return
	}
	log.Printf("account:%v", userDb.ID)

	// check if user has account
	_, err = a.server.Store.Queries.GetAccountByUserId(context.Background(), account.UserID)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the user already has an account"})
		return
	}

	arg := sqlc.CreateAccountParams{
		UserID:  account.UserID,
		Balance: account.Balance,
	}

	accountDb, err := a.server.Store.Queries.CreateAccount(context.Background(), arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, accountDb)
}

type transferRequest struct {
	ToAccountID   int32 `json:"to_account_id" binding:"required"`
	Amount        int32 `json:"amount" binding:"required"`
	FromAccountID int32 `json:"from_account_id" binding:"required"`
}

func (a *AccountApp) transfer(c *gin.Context) {
	userId, err := auth.GetActiveUser(c)
	if err != nil {
		return
	}

	tr := new(transferRequest)

	if err := c.ShouldBind(&tr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account, err := a.server.Store.Queries.GetAccountById(context.Background(), int64(tr.FromAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not get account"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account.UserID != int32(userId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not get account"})
		return
	}

	_, err = a.server.Store.Queries.GetAccountById(context.Background(), int64(tr.ToAccountID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not get account to send to"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if account.Balance < float64(tr.Amount) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You dont have enough balance"})
		return
	}

	txArg := sqlc.CreateTransferParams{
		FromAccountID: tr.FromAccountID,
		ToAccountID:   tr.ToAccountID,
		Amount:        float64(tr.Amount),
	}
	tx, err := a.server.Store.TransferTx(context.Background(), txArg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Encountered issue with transaction"})
		return
	}

	c.JSON(http.StatusCreated, tx)

}

func (a *AccountApp) myAccount(c *gin.Context) {
	userId, err := auth.GetActiveUser(c)
	if err != nil {
		return
	}

	account, err := a.server.Store.GetAccountByUserId(context.Background(), int32(userId))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you dont have and account"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, account)

}

func (a *AccountApp) myTransfers(c *gin.Context) {
	userId, err := auth.GetActiveUser(c)
	if err != nil {
		return
	}

	dir := c.DefaultQuery("dir", "out")

	if dir == "out" {
		transfers, err := a.server.Store.GetTransfersByFromAccountID(context.Background(), int32(userId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, transfers)

	} else if dir == "in" {
		transfers, err := a.server.Store.GetTransfersByToAccountID(context.Background(), int32(userId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, transfers)
	} else {

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type of direction"})
		return
	}

}
