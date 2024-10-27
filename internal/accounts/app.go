package accounts

import (
	"github.com/BenjaminArmijo3/bank/internal/app"
)

type AccountApp struct {
	server *app.Server
}

func NewAccountApp(server *app.Server) *AccountApp {
	return &AccountApp{
		server: server,
	}
}
