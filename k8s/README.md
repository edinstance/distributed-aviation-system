# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the distributed aviation system components.



## Components

### Namespace

- **File**: `00-namespace.yml`
- **Purpose**: Creates the `aviation` namespace for all resources

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
