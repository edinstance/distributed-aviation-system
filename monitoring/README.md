# Monitoring

This directory contains all the configurations for monitoring.

It is configured to work with all the services deployed using Docker Compose
and [this configuration](../docker-compose.yml).

## Prometheus

Prometheus is configured to scrape metrics from the services. The configuration is located
in [prometheus.yml](./prometheus/prometheus.yml).