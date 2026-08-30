package app

import (
	"context"
	"database/sql"

	"github.com/ModstDev/Pokerer/internal/auth"
	"github.com/ModstDev/Pokerer/internal/auth/token"
	database "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/poker"
	"github.com/ModstDev/Pokerer/internal/poker/table"
	"github.com/ModstDev/Pokerer/internal/repository"
	"github.com/ModstDev/Pokerer/internal/wallet"
)

// We create seperate app package to avoid mess in main.go
// This is our dependency container.

type App struct {
	DB *sql.DB

	Users   *repository.UserRepository
	Wallets *repository.WalletRepository

	Wallet *wallet.Service
	Auth   *auth.Service
	Poker  *poker.Service

	Token *token.JWT
}

func New(db *sql.DB, jwtSecret, jwtIssuer string) *App {
	queries := database.New(db)

	userRepository := repository.NewUserRepository(queries)
	walletRepository := repository.NewWalletRepository(queries)
	tableRepository := repository.NewPokerTableRepository(queries)

	pokerService := poker.NewService(db, tableRepository, walletRepository)
	tableManager := table.NewManager(context.Background(), pokerService)

	pokerService.SetTableManager(tableManager)

	return &App{
		DB: db,

		Users:   repository.NewUserRepository(queries),
		Wallets: walletRepository,

		Wallet: wallet.NewService(db, walletRepository),
		Auth:   auth.NewService(db, userRepository, walletRepository),
		Poker:  pokerService,

		Token: token.NewJWT(jwtSecret, jwtIssuer),
	}
}
