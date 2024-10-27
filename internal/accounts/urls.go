package accounts

import (
	"github.com/gin-gonic/gin"
)

func (a *AccountApp) RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(a.server.AuthenticatedMiddleware())
	rg.POST("/create", a.createAccount)
	rg.POST("/transfer", a.transfer)
	rg.GET("/my_account", a.myAccount)
	rg.GET("/my_transfers", a.myTransfers)
}
