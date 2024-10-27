package main

import (
	"database/sql"
	"log"

	"github.com/BenjaminArmijo3/bank/internal/accounts"
	"github.com/BenjaminArmijo3/bank/internal/app"
	"github.com/BenjaminArmijo3/bank/internal/auth"
	"github.com/BenjaminArmijo3/bank/internal/config"
	"github.com/BenjaminArmijo3/bank/internal/db/store"
	"github.com/BenjaminArmijo3/bank/internal/pkg/utils/token"
	"github.com/BenjaminArmijo3/bank/internal/users"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {

	cfg := config.MustLoad()

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	// Set migration directory
	dir := cfg.MigrationsPath

	// Run migrations
	if err := goose.Up(db, dir); err != nil {
		log.Fatal(err)
	}

	token := token.NewJWTToken(cfg)
	store := store.NewStore(db)
	server := app.NewServer(store, cfg, token)

	authApp := auth.NewAuthApp(server)
	server.RegisterRoutes("api/auth", authApp)

	usersApp := users.NewUsersApp(server)
	server.RegisterRoutes("api/users", usersApp)

	accountsApp := accounts.NewAccountApp(server)
	server.RegisterRoutes("api/accounts", accountsApp)

	if err := server.Start(cfg.Address); err != nil {
		panic(err)
	}
}
