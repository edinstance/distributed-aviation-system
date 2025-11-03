### Aircraft Service Stack

#### Database

- **Type**: PostgreSQL 17.6 StatefulSet
- **Storage**: 1Gi persistent volume
- **Port**: 5432
- **Credentials**: postgres/postgres
- **Database**: aircraft

#### Cache

- **Type**: Redis 8.2 Deployment
- **Port**: 6379
- **Purpose**: Caching layer for aircraft service

#### Migrations

- **Type**: Job
- **Image**: `liquibase/liquibase:4.31`
- **Dependencies**: Waits for database to be ready

#### Service

- **Type**: Deployment with Service
- **Port**: 8080
- **Image**: `ghcr.io/edinstance/distributed-aviation-system-services/aircraft:latest`
- **Dependencies**: Waits for database, migrations and cache to be ready
- **Health Check**: `/health` endpoint
