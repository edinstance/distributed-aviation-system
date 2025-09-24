# Aircraft Service

This is a [Java](https://www.java.com/) microservice for managing aircraft data in the larger aviation system.

## Prerequisites

- Java 21
- PostgreSQL 17+
- Docker (optional, for containerized deployment)
- Docker Compose (optional, for local development with a database)

### Required Environment Variables

The service requires the following env vars:

- `ENVIRONMENT`: `dev`, `staging`, or `prod` (controls logging/telemetry)
- `PORT`: The port the service will listen on (e.g. `8080`)
- `DATABASE_URL`: Postgres connection string in the form `postgres://host:5432/dbname?sslmode=disable`
- `DATABASE_USERNAME`: The username for accessing the database.
- `DATABASE_PASSWORD`: The password for accessing the database.
- `LIQUIBASE_ENABLED`: A toggle for if the service should run liquibase and the database updates.
- `CACHE_URL`: Redis connection string in the form of `redis://host:6379`

See [.env.example](.env.example) for a template.

## Development Setup

### 1. Install Dependencies

```bash

  mvn clean install
```

### 2. Database Setup

Start PostgreSQL locally or use Docker:

```bash

# Using Docker
docker run -d \
  --name postgres-aircraft \
  -e DATABASE_URL=postgres://postgres@localhost:5432/aircraft \
  -p 5432:5432 \
  postgres:17
```

### 3. Run Migrations

Migrations are stored in [this folder](src/main/resources/db/changelog/db.changelog-master.xml) and are applied
using [liquibase](https://www.liquibase.com/). They can be applied using the liquibase cli, the liquibase docker image
or the service can be configured to run the migrations by using the LIQUIBASE_ENABLED env variable mentioned above.

The CLI command, to install follow [these instructions](https://www.liquibase.com/download-oss).

```bash
liquibase \
  --url="${DATABASE_URL:-jdbc:postgresql://localhost:5432/aircraft?sslmode=disable}" \
  --username="${DATABASE_USERNAME:-postgres}" \
  --password="${DATABASE_PASSWORD:-postgres}" \
  --changelog-file="src/main/resources/db/changelog/db.changelog-master.xml" \
  update
```

### 4. Run the Service

```bash

# Set environment variables or use the .env file
export PORT=8080
export ENVIRONMENT=dev
export DATABASE_URL="jdbc:postgresql://localhost:5432/aircraft"
export DATABASE_USERNAME=postgres
export DATABASE_PASSWORD=postgres
export LIQUIBASE_ENABLED-true
export CACHE_URL=redis://localhost:6380

# Run the service
mvn spring-boot:run
```

The service will start on `http://localhost:8080`

## Building

### Local Build

```bash
# Build binary
mvn clean package

# Run binary
java -jar target/aircraft-0.0.1-SNAPSHOT.jar
```

## Docker

### Build Docker Image

```bash

docker build -t aircraft-service .
```

### Run with Docker

```bash

# Run container (assumes PostgreSQL is accessible)
docker run -d \
  --name aircraft-service \
  -p 8080:8080 \
  -e PORT=8080 \
  -e ENVIRONMENT=dev \
  -e DATABASE_URL="jdbc:postgresql://host.docker.internal:5432/aircraft" \
  -e DATABASE_USERNAME=postgres \
  -e DATABASE_PASSWORD=postgres \
  -e LIQUIBASE_ENABLED=true \
  -e CACHE_URL="redis://host.docker.internal:6380" \
  aircraft-service
```

### Docker Compose

This uses the [docker-compose.yml](docker-compose.yml) file. Make sure the migrations have been run on the Docker
database as well and that the environment variables are set in the [.env](.env) file.

```bash

docker compose up -d
```

## API Endpoints

The service exposes a graphql endpoint and HTTP routes:

- Health check: `GET details/health`
- Aircraft operations: Graphql endpoint at `/graphql`

## Development

### Testing

```bash
# Run tests
mvn test
```

