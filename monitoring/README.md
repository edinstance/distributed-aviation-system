# Monitoring

This directory contains all the configurations for monitoring.

It is configured to work with all the services deployed using Docker Compose
and [this configuration](../docker-compose.yml). There are some environment variables
that need to be set, these currently configure grafana's admin user and password. To set them you need to take the
values from [.env.example](.env.example) and put them in a [.env](.env) file.

## Prometheus

Prometheus is configured to scrape metrics from the services. The configuration is located
in [prometheus.yml](./prometheus/prometheus.yml).

## Grafana

Grafana is configured to import dashboards from [dashboards](./grafana/provisioning/dashboards). Grafana gets the data
for these dashboards from the [datasources](grafana/provisioning/datasources) directory.

