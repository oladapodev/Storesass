# storefront-saas-go

A production-ready, multi-tenant storefront SaaS template built with Go (Gin) and React (Vite + TypeScript). Features clean architecture, type-safe API codegen pipeline, Redis caching, SQLite for local dev with one-line PostgreSQL switch, Docker deployment, and hot reload.

## Tech Stack

**Backend**
- Go 1.24, Gin v1.11, GORM v1.25
- SQLite (local dev) / PostgreSQL (production)
- Redis (optional caching layer)
- swaggo/swag for Swagger/OpenAPI generation
- Air for hot reload

**Frontend**
- React 19, TypeScript strict, Vite 6
- TanStack Query v5 for data fetching
- shadcn/ui components (Tailwind v4, dark mode ready)
- React Router v7
- Orval for TanStack Query hook generation
- Vitest for unit tests

**Codegen Pipeline**
```
swag init -> docs/swagger.json
  -> swagger2openapi -> docs/openapi.json
    -> openapi-typescript -> web/src/types/api.ts
    -> orval -> web/src/api/ (TanStack Query hooks)
```

## Quick Start

```bash
# 1. Clone and copy env
cp .env.example .env

# 2. Install Go tools and frontend deps
make setup

# 3. Run the codegen pipeline
make codegen

# 4. Seed demo data
make seed

# 5. Start backend (terminal 1)
make dev-backend

# 6. Start frontend (terminal 2)
make dev-frontend
```

Open:
- Frontend: http://localhost:5173
- API: http://localhost:8080/api/v1
- Swagger UI: http://localhost:8080/swagger/index.html

## Project Structure

```
storefront-saas-go/
├── cmd/
│   ├── api/         # HTTP server entry point
│   └── seed/        # Demo data seeder (DELETE for production)
├── internal/
│   ├── config/      # Environment config
│   ├── db/          # GORM + Redis init
│   ├── domain/      # Pure Go structs (User, Store, Product, Order)
│   ├── handler/     # Thin Gin HTTP handlers
│   ├── middleware/  # CORS, etc.
│   ├── repository/  # GORM queries only
│   ├── service/     # Business logic + Redis caching
│   └── util/        # Response helpers
├── docs/            # Swagger docs (generated, do not edit)
├── web/             # React 19 frontend
│   ├── src/
│   │   ├── api/         # Orval-generated hooks (do not edit)
│   │   ├── types/       # openapi-typescript types (do not edit)
│   │   ├── components/  # UI components
│   │   ├── pages/       # Route pages
│   │   └── lib/         # Utilities
│   └── orval.config.ts
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── .env.example
```

## Makefile Targets

```bash
make help           # Show all targets
make setup          # First-time setup (install tools + deps)
make dev-backend    # Start Go API with hot reload (Air)
make dev-frontend   # Start Vite dev server
make build          # Build Go binary
make test           # Run all tests (Go + Vitest)
make test-backend   # Go tests only
make test-frontend  # Vitest only
make lint           # Go vet
make seed           # Seed demo data
make codegen        # Full codegen pipeline
make swagger        # Regenerate swagger docs
make openapi-types  # Regenerate TypeScript types
make orval-hooks    # Regenerate TanStack Query hooks
make docker-up      # Start with Docker Compose
make docker-down    # Stop Docker services
make docker-logs    # Tail API logs
make clean          # Remove build artifacts
```

## Database Switch

**SQLite (local dev, default)**
```env
DB_DRIVER=sqlite
DB_DSN=file:./dev.db
```

**PostgreSQL (production / Docker)**
```env
DB_DRIVER=postgres
DB_DSN=host=localhost user=storefront password=storefront dbname=storefront port=5432 sslmode=disable
```

GORM auto-migrates all models on startup. No migrations to run manually.

## Docker Deploy

```bash
# Start PostgreSQL 17 + Redis 8 + API
make docker-up

# Check logs
make docker-logs

# Stop
make docker-down
```

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/api/v1/stores` | List active stores (paginated) |
| POST | `/api/v1/stores` | Create a store |
| GET | `/api/v1/stores/:slug` | Get store by slug |
| GET | `/api/v1/stores/:slug/products` | List products for a store |
| GET | `/api/v1/products` | List hot products (Redis cached) |
| GET | `/api/v1/products/search?q=` | Search products |
| GET | `/api/v1/products/:id` | Get product by ID |
| POST | `/api/v1/products` | Create a product |
| GET | `/swagger/*any` | Swagger UI |

## Removing Demo Data

All demo data is in `cmd/seed/main.go` and is never seeded automatically. To remove it:

```bash
rm cmd/seed/main.go
# Remove the `seed` target from Makefile
```

## Adding New Endpoints

See `AGENT.md` for the step-by-step guide to adding new endpoints while maintaining the full codegen pipeline.

## License

MIT