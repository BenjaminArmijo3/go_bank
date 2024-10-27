package auth

import (
	"github.com/BenjaminArmijo3/bank/internal/app"
)

type AuthApp struct {
	server *app.Server
}

func NewAuthApp(server *app.Server) *AuthApp {
	return &AuthApp{
		server: server,
	}
}
