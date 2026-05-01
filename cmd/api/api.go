package main

import (
	"context"
	"log"

	"github.com/floffah/catena/internal/app/api"
	"github.com/floffah/catena/internal/pkg/auth"
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/environment"
)

func main() {
	config, err := environment.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.Connect(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	authService := auth.NewAuthService(config.ClerkSecretKey)

	server, err := api.NewServer(conn, authService, config.CORSAllowedOrigins)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
