package app

import (
	"database/sql"

	database "github.com/ModstDev/Pokerer/internal/database/generated"
	"github.com/ModstDev/Pokerer/internal/repository"
)

// We create seperate app package to avoid mess in main.go
// This is our dependency container.

type App struct {
	DB *sql.DB

	Users *repository.UserRepository
}

func New(db *sql.DB) *App {
	queries := database.New(db)

	return &App{
		DB:    db,
		Users: repository.NewUserRepository(queries),
	}
}
