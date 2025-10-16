# Aviation Gateway Service

An authentication and authorization gateway service built with Rust and Axum that provides JWT verification and request forwarding for the distributed aviation system.

## Overview

The Gateway Service acts as a secure entry point for the aviation system, handling JWT token verification and forwarding authenticated requests to downstream services. It validates tokens against JWKS (JSON Web Key Set) and injects user context headers for downstream services.

## Features

- **JWT Authentication**: Validates Bearer tokens using JWKS
- **Request Forwarding**: Proxies authenticated requests to the GraphQL router
- **User Context Injection**: Adds user metadata headers for downstream services
- **Observability**: Comprehensive logging, tracing, and metrics with OpenTelemetry
- **JWKS Caching**: Efficient token validation with cached public keys
- **Request ID Tracking**: Distributed tracing support

## Architecture

The gateway validates incoming requests with JWT tokens, extracts user claims, and forwards authenticated requests to the GraphQL router with additional context headers:

```
Client → Gateway (JWT Validation) → GraphQL Router → Microservices
```

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# JWT Configuration
JWKS_URL=http://localhost:8000/.well-known/jwks.json
ROUTER_URL=http://localhost:4000

# Server Configuration
PORT=4001

# Observability Configuration
SERVICE_NAME=aviation-gateway
SERVICE_VERSION=0.1.0
LOG_LEVEL=info
JSON_LOGS=false

# OpenTelemetry Configuration
ENABLE_OTLP=false
OTLP_ENDPOINT=http://localhost:4317
```

### Required Environment Variables

- `JWKS_URL`: URL to fetch JSON Web Key Set for token validation
- `ROUTER_URL`: GraphQL router endpoint to forward requests to
- `PORT`: Port for the gateway service (default: 4001)

## Development

### Prerequisites

- Rust 2024 edition
- Docker (optional)

### Running Locally

```bash
# Install dependencies
cargo build

# Run the service
cargo run

# Run with specific environment
RUST_LOG=debug cargo run
```

### Docker

```bash
# Build image
docker build -t aviation-gateway .

# Run container
docker-compose up gateway-service
```

## API

The gateway accepts all HTTP methods and paths, acting as a transparent proxy after authentication:

### Authentication

All requests must include a valid Bearer token:

```http
Authorization: Bearer <jwt-token>
```

The Bearer tokens should be retrieved from the authentication microservice which can be queried directly.

### Request Flow

1. Extract JWT from `Authorization` header
2. Validate token against JWKS endpoint
3. Add user context headers:
   - `x-user-sub`: User ID from token
   - `x-org-id`: Organization ID from token
   - `x-user-roles`: Comma-separated user roles
   - `x-request-id`: Unique request identifier
4. Forward request to GraphQL router

### Response Codes

- `200-299`: Successful request forwarded
- `401`: Missing, invalid, or expired JWT token
- `5xx`: Gateway or downstream service errors

## Dependencies

Key dependencies from `Cargo.toml`:

- **axum**: Web framework for HTTP handling
- **jsonwebtoken**: JWT token validation
- **reqwest**: HTTP client for forwarding requests
- **tokio**: Async runtime
- **tracing**: Structured logging and observability
- **opentelemetry**: Distributed tracing and metrics

## Monitoring

The service provides comprehensive observability:

- **Structured Logging**: JSON logs with request correlation
- **Metrics**: Request counts, durations, and status codes
- **Tracing**: OpenTelemetry spans for distributed tracing
- **Health**: Service readiness and liveness indicators

## Security

- Validates all requests with JWT tokens
- Caches JWKS keys with automatic refresh
- Secure token parsing and validation
- Request timeout protection (30s)
- No credential logging or exposure

## Development Notes

The service is designed for high availability and security:

- Async request handling with Tokio
- Connection pooling and timeouts
- Graceful shutdown handling
- Comprehensive error handling and logging
- Stateless design for horizontal scaling
