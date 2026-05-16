package main

import (
	"context"
	"log"

	"github.com/floffah/catena/internal/app/api"
	"github.com/floffah/catena/internal/pkg/auth"
	"github.com/floffah/catena/internal/pkg/db"
	"github.com/floffah/catena/internal/pkg/environment"
	"github.com/floffah/catena/internal/pkg/gitstore"
	"github.com/gin-gonic/gin"
)

func main() {
	env, err := environment.LoadEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	ginMode := "debug"
	if env.Config.Mode == "production" {
		ginMode = "release"
	}
	gin.SetMode(ginMode)

	conn, err := db.Connect(context.Background(), env.Config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	authService := auth.NewAuthService(env.Config.ClerkSecretKey, conn)

	gitService := gitstore.NewStoreFromEnv(env)

	server, err := api.NewServer(conn, authService, gitService, env.Config.CORSAllowedOrigins, env.Config.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.ListenAndServe())
}
