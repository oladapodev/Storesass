BINARY_NAME=storefront-saas-api
GOPATH=$(shell go env GOPATH)
AIR=$(GOPATH)/bin/air
SWAG=$(GOPATH)/bin/swag

.PHONY: help dev dev-backend dev-frontend build clean test lint seed \
        codegen swagger openapi-types orval-hooks \
        docker-up docker-down docker-build docker-logs \
        install-tools install-frontend

## help: Show this help message
help:
	@echo "Storefront SaaS - Available targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'

## dev-backend: Start Go backend with Air hot reload
dev-backend:
	$(AIR) -c .air.toml

## dev-frontend: Start Vite dev server
dev-frontend:
	cd web && pnpm dev

## dev: Start both backend and frontend (requires two terminals)
dev:
	@echo "Run 'make dev-backend' and 'make dev-frontend' in separate terminals."
	@echo "Or use: make dev-backend & make dev-frontend"

## build: Build the Go binary
build:
	go build -o bin/$(BINARY_NAME) ./cmd/api

## clean: Remove build artifacts
clean:
	rm -rf bin/ docs/swagger.json docs/swagger.yaml docs/docs.go web/dist web/src/api web/src/types/api.ts

## test: Run all tests (Go + Vitest)
test: test-backend test-frontend

## test-backend: Run Go tests
test-backend:
	go test -v -count=1 ./...

## test-frontend: Run Vitest tests
test-frontend:
	cd web && pnpm test

## lint: Lint Go code
lint:
	go vet ./...

## seed: Seed database with demo data
seed:
	go run cmd/seed/main.go

## codegen: Full codegen pipeline (swagger -> openapi -> types -> orval hooks)
codegen: swagger openapi-types orval-hooks

## swagger: Generate Swagger docs from Go annotations
swagger:
	$(SWAG) init -g cmd/api/main.go -o docs

## openapi-types: Convert swagger.json to openapi.json and generate TS types
openapi-types:
	cd web && npx swagger2openapi ../docs/swagger.json -o ../docs/openapi.json
	cd web && pnpm run codegen:types

## orval-hooks: Generate TanStack Query hooks with Orval
orval-hooks:
	cd web && pnpm run codegen:api

## install-tools: Install Go dev tools (air, swag)
install-tools:
	go install github.com/air-verse/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest

## install-frontend: Install frontend dependencies
install-frontend:
	cd web && pnpm install

## docker-up: Start all services with Docker Compose
docker-up:
	docker compose up --build -d

## docker-down: Stop all Docker services
docker-down:
	docker compose down

## docker-build: Build the Docker image only
docker-build:
	docker build -t $(BINARY_NAME) .

## docker-logs: Tail Docker Compose logs
docker-logs:
	docker compose logs -f api

## setup: First-time project setup
setup: install-tools install-frontend
	cp -n .env.example .env || true
	@echo "Setup complete. Edit .env then run: make seed && make dev"
