package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go tool sqlc generate

var conn *pgxpool.Pool

func GetConn(ctx context.Context) (*pgxpool.Pool, error) {
	if conn == nil {
		var err error
		conn, err = pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}
