package users

import (
	"github.com/BenjaminArmijo3/bank/internal/app"
)

type UsersApp struct {
	server *app.Server
}

func NewUsersApp(server *app.Server) *UsersApp {
	return &UsersApp{
		server: server,
	}
}
