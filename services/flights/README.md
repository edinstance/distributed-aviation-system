# Flights Service

This is a [go](https://go.dev/) microservice for managing flight data in the larger aviation system.

## Prerequisites

- Go 1.25.1
- PostgreSQL 17+
- Docker (optional, for containerized deployment)
- Docker Compose (optional, for local development with a database)

### Required Environment Variables

The service requires the following env vars:

- `PORT`: The port the service will listen on (e.g. `8081`)
- `DATABASE_URL`: Postgres connection string in the form `postgres://user:pass@host:5432/dbname?sslmode=disable`
- `ENVIRONMENT`: `dev`, `staging`, or `prod` (controls logging/telemetry)

See [.env.example](.env.example) for a template.

## Development Setup

### 1. Install Dependencies

```bash

  go mod download
```

### 2. Database Setup

Start PostgreSQL locally or use Docker:

```bash

# Using Docker
docker run -d \
  --name postgres-flights \
  -e DATABASE_URL=postgres://postgres@localhost:5432/flights \
  -p 5432:5432 \
  postgres:15
```

### 3. Run Migrations

Migrations are stored in `migrations/` and are applied
using [golang-migrate](https://github.com/golang-migrate/migrate).
They must be applied in order — use `migrate` to handle that automatically.

```bash

migrate -path migrations -database "$DATABASE_URL" up
```

### 4. Run the Service

```bash

# Set environment variables or use the .env file
export PORT=8081
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/flights"

# Run the service
go run cmd/main.go
```

The service will start on `http://localhost:8081`

## Building

### Local Build

```bash
# Build binary
go build -o flights-service cmd/main.go

# Run binary
./flights-service
```

### Cross-platform Build

```bash
# Build for Linux (useful for Docker)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o flights-service cmd/main.go
```

## Docker

### Build Docker Image

```bash

docker build -t flights-service .
```

### Run with Docker

```bash

# Run container (assumes PostgreSQL is accessible)
docker run -d \
  --name flights-service \
  -p 8081:8081 \
  -e PORT=8081 \
  -e ENVIRONMENT=dev \
  -e DATABASE_URL=postgres://postgres:postgres@host.docker.internal:5432/flights?sslmode=disable
  flights-service
```

### Docker Compose

This uses the [docker-compose.yml](docker-compose.yml) file. Make sure the migrations have been run on the Docker database as well and that the environment variables are set in the [.env](.env) file.

```bash

docker compose up -d
```

## Database Migrations

### Manual Migration Management

The service includes SQL migration files in the `migrations/` directory:

- `000001_update_updated_date.up.sql` / `000001_update_updated_date.down.sql`
- `000002_create_flights.up.sql` / `000002_create_flights.down.sql`

### Using golang-migrate

Install the migration tool by following steps [here](https://github.com/golang-migrate/migrate).

Run migrations:

```bash

# Up migrations
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/flights?sslmode=disable" up

# Down migrations (rollback)
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/flights?sslmode=disable" down

# Migrate to specific version
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/flights?sslmode=disable" goto 1
```

Create new migrations:

```bash

migrate create -ext sql -dir migrations -seq add_new_table
```

## API Endpoints

The service exposes Connect RPC endpoints and HTTP routes:

- Health check: `GET /health`
- Flight operations: Connect RPC endpoints for flight management

## Development

### Code Formatting

```bash

go fmt ./...
```

### Linting

```bash

# Install golangci-lint if not already installed
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```
