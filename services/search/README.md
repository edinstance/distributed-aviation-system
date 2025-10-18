# Aircraft Service

This is a [Java](https://www.java.com/) microservice for searching for data in the larger aviation system.

## Prerequisites

- Java 21
- Docker (optional, for containerized deployment)
- Docker Compose (optional, for local development with a database)

### Required Environment Variables

See [.env.example](.env.example) for a template.

## Development Setup

### 1. Install Dependencies

```bash

  mvn clean install
```

### 2. Run the Service

```bash

# Set environment variables or use the .env file
export PORT=8082

# Run the service
mvn spring-boot:run
The service will start on `http://localhost:8082`

## Building

### Local Build

```bash
# Build binary
mvn clean package

# Run binary
java -javaagent:target/opentelemetry-javaagent.jar \
     -Dotel.service.name=search-service \
     -Dotel.exporter.otlp.endpoint=http://localhost:4317 \
     -jar target/aircraft-0.0.1-SNAPSHOT.jar
```

## Docker

### Build Docker Image

```bash

docker build -t search-service .
```

### Run with Docker

```bash

# Run container (assumes PostgreSQL is accessible)
docker run -d \
  --name search-service \
  -p 8082:8080 \
  -e PORT=8082 \
  search-service
```

### Docker Compose

This uses the [docker-compose.yml](docker-compose.yml) file.

```bash

docker compose up -d
```

## Development

### Code Quality

#### Formatting

#### Linting

This project uses [checkstyle](https://checkstyle.sourceforge.io) with the Google Checkstyle configuration which is in [this folder](config).:

```bash

# Run linter
mvn checkstyle:checkstyle
```

### Testing

```bash
# Run tests
mvn test
```

## Telemetry Summary

- Agent download is automated via Maven plugin → target/opentelemetry-javaagent.jar
- Included in Dockerfile and always attached at runtime
- Configured via environment variables (OTEL_*)
- Supports OTLP → collector (Tempo, Grafana, Jaeger, Prometheus, etc.)