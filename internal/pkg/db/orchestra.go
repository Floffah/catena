package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go tool sqlc generate

func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, url)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
