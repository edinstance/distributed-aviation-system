# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the distributed aviation system components.



## Components

### Namespace

- **File**: `00-namespace.yml`
- **Purpose**: Creates the `aviation` namespace for all resources

### Flights Service Stack

Located in the `flights/` directory:

#### Database

- **Type**: PostgreSQL 17.6 StatefulSet
- **Storage**: 1Gi persistent volume
- **Port**: 5432
- **Credentials**: postgres/postgres
- **Database**: flights

#### Cache

- **Type**: Redis 8.2 Deployment
- **Port**: 6379
- **Purpose**: Caching layer for flights service

#### Migrations

- **Type**: Job
- **Image**: `migrate:migrate`
- **Dependencies**: Waits for database to be ready

#### Service

- **Type**: Deployment with Service
- **Port**: 8081
- **Image**: `ghcr.io/edinstance/distributed-aviation-system-services/flights:latest`
- **Dependencies**: Waits for database, migrations and cache to be ready
- **Health Check**: `/health` endpoint

## Deployment

1. Apply the namespace first:
   ```bash
   kubectl apply -f 00-namespace.yml
   ```

2. Deploy the flights stack:
   ```bash
   kubectl apply -f flights/
   ```

3. Verify deployment:
   ```bash
   kubectl get pods -n aviation
   ```

Or you can deploy everything with:

```bash
kubectl apply -R -f k8s/                 
```

## Environment Variables

The flights service is configured with:

- Database connection to internal PostgreSQL
- Redis cache connection
- Kafka broker and schema registry endpoints
- OTLP observability endpoint
- gRPC aircraft service endpoint

## Notes

- All services run in the `aviation` namespace
- Database uses persistent storage via StatefulSet
