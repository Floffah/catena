set shell := ["bash", "-eu", "-o", "pipefail", "-c"]
set dotenv-load := true

alias gen := generate
alias da := dev-api
alias dw := dev-web
alias d := dev

# These are embedded so we can distribute the CLI without telling users our clerk publishable key (secret should never be embedded)
common_ldflags := "-X github.com/floffah/catena/internal/pkg/auth.ClerkPublishableKey=$CLERK_PUBLISHABLE_KEY -X github.com/floffah/catena/internal/pkg/auth.ClerkFrontendApiUrl=$CLERK_FRONTEND_API_URL"
cli_ldflags := common_ldflags
api_ldflags := common_ldflags

default:
    @just --list

build:
	go build -ldflags="{{cli_ldflags}}" -o build/catena cmd/catena/catena.go
	go build -ldflags="{{api_ldflags}}" -o build/api cmd/api/api.go

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
	go tool wgo run cmd/api/api.go

dev-web:
	cd web && bun run dev

dev-db:
	docker compose -f deployments/dev.docker-compose.yml up -d db

[confirm]
dev-reset-gitstore:
	rm -rf ${CATENA_GIT_ROOT}

dev-reset:
	@just db-reset
	@just dev-reset-gitstore

# -- db --

db-new-migration NAME:
    TERN_MIGRATIONS=./data/migrations/ go tool tern new {{NAME}}

db-migrate:
    TERN_MIGRATIONS=./data/migrations/ go tool tern migrate

db-rollback:
	TERN_MIGRATIONS=./data/migrations/ go tool tern migrate -d -1

[confirm]
db-reset:
    TERN_MIGRATIONS=./data/migrations/ go tool tern migrate -d 0
    @just db-migrate