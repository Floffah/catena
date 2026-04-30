set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

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
dev: dev-api dev-web

dev-api:
	go run cmd/api/api.go

dev-web:
	cd web && bun run dev