### Flights Service Stack

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
