set shell := ["bash", "-eu", "-o", "pipefail", "-c"]
set dotenv-load := true

alias gen := generate
alias da := dev-api
alias dw := dev-web
alias d := dev

default:
    @just --list

generate:
	go generate ./...
	bun run --cwd web generate

lint:
	golangci-lint run

test:
	go test -v -coverprofile=coverage.out ./...

format:
	go fmt ./...
	golangci-lint fmt

check: lint test

# -- dev --

[parallel]
dev: dev-api dev-web dev-db

dev-api:
	go run cmd/api/api.go

dev-web:
	cd web && bun run dev

dev-db:
	docker compose -f deployments/dev.docker-compose.yml up -d db

# -- db --

db-new-migration NAME:
    TERN_MIGRATIONS=./data/migrations/ go tool tern new {{NAME}}

db-migrate:
    TERN_MIGRATIONS=./data/migrations/ go tool tern migrate