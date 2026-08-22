package main

import (
	"log"

	"github.com/ModstDev/Pokerer/internal/config"
	"github.com/ModstDev/Pokerer/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	db, err := database.Connect(database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
	})
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer db.Close()

	log.Println("database connection established")
}
