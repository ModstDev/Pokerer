package app

import (
	"database/sql"

	"github.com/ModstDev/Pokerer/internal/auth"
	"github.com/ModstDev/Pokerer/internal/auth/token"
	database "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
)

// We create seperate app package to avoid mess in main.go
// This is our dependency container.

type App struct {
	DB *sql.DB

	Users *repository.UserRepository
	Auth  *auth.Service
	Token *token.JWT
}

func New(db *sql.DB, jwtSecret, jwtIssuer string) *App {
	queries := database.New(db)

	userRepository := repository.NewUserRepository(queries)

	return &App{
		DB:    db,
		Users: repository.NewUserRepository(queries),
		Auth:  auth.NewService(userRepository),
		Token: token.NewJWT(jwtSecret, jwtIssuer),
	}
}
