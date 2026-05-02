package main

import (
	"context"
	"log"

	"github.com/floffah/catena/internal/app/api"
	"github.com/floffah/catena/internal/pkg/auth"
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/environment"
	"github.com/floffah/catena/internal/pkg/gitstore"
)

func main() {
	env, err := environment.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.Connect(context.Background(), env.Config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	authService := auth.NewAuthService(env.Config.ClerkSecretKey, conn)

	gitService := gitstore.NewStoreFromEnv(*env)

	server, err := api.NewServer(conn, authService, gitService, env.Config.CORSAllowedOrigins)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
