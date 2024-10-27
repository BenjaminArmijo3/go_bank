package auth

import (
	"github.com/gin-gonic/gin"
)

func (a *AuthApp) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("sign-up", a.register)
	rg.POST("sign-in", a.login)
}
