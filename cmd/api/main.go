package main

import (
	"log"
	"net/http"

	"github.com/ModstDev/Pokerer/internal/app"
	"github.com/ModstDev/Pokerer/internal/config"
	"github.com/ModstDev/Pokerer/internal/database"
	"github.com/ModstDev/Pokerer/internal/httpapi"
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

	application := app.New(db, cfg.JWT.Secret, cfg.JWT.Issuer)

	server := httpapi.NewServer(
		application.Auth,
		application.Token,
		application.Users,
		application.Wallet,
	)

	log.Println("server listening on :8080")

	if err := http.ListenAndServe(":8080", server.Handler()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
