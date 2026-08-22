package app

import (
	"database/sql"

	"github.com/ModstDev/Pokerer/internal/auth"
	database "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
)

// We create seperate app package to avoid mess in main.go
// This is our dependency container.

type App struct {
	DB *sql.DB

	Users *repository.UserRepository
	Auth  *auth.Service
}

func New(db *sql.DB) *App {
	queries := database.New(db)

	userRepository := repository.NewUserRepository(queries)

	return &App{
		DB:    db,
		Users: repository.NewUserRepository(queries),
		Auth:  auth.NewService(userRepository),
	}
}
