"""
OpenTelemetry instrumentation setup for the authentication service.
This module configures distributed tracing, metrics, and logging collection.
"""

import logging
import os

from opentelemetry import trace, metrics
from opentelemetry.exporter.otlp.proto.grpc._log_exporter import OTLPLogExporter
from opentelemetry.exporter.otlp.proto.grpc.metric_exporter import OTLPMetricExporter
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter
from opentelemetry.exporter.prometheus import PrometheusMetricReader
from opentelemetry.instrumentation.django import DjangoInstrumentor
from opentelemetry.instrumentation.logging import LoggingInstrumentor
from opentelemetry.instrumentation.psycopg2 import Psycopg2Instrumentor
from opentelemetry.instrumentation.requests import RequestsInstrumentor
from opentelemetry.sdk._logs import LoggerProvider, LoggingHandler
from opentelemetry.sdk._logs.export import BatchLogRecordProcessor
from opentelemetry.sdk.metrics import MeterProvider
from opentelemetry.sdk.metrics.export import PeriodicExportingMetricReader
from opentelemetry.sdk.resources import SERVICE_NAME, SERVICE_VERSION, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor, ConsoleSpanExporter

from common.telemetry.auth_metrics import AuthMetrics

SERVICE_NAME_VALUE = "authentication-service"
SERVICE_VERSION_VALUE = "1.0.0"


def configure_telemetry(otel_endpoint: str = None):
    resource = Resource(
        attributes={
            SERVICE_NAME: SERVICE_NAME_VALUE,
            SERVICE_VERSION: SERVICE_VERSION_VALUE,
            "service.instance.id": os.environ.get("HOSTNAME", "unknown"),
            "deployment.environment": os.environ.get("ENVIRONMENT", "dev"),
        }
    )

    DjangoInstrumentor().instrument()
    Psycopg2Instrumentor().instrument()

    # ---- Tracing ----
    if otel_endpoint:
        span_exporter = OTLPSpanExporter(endpoint=otel_endpoint, insecure=True)
    else:
        span_exporter = ConsoleSpanExporter()

    tracer_provider = TracerProvider(resource=resource)
    tracer_provider.add_span_processor(BatchSpanProcessor(span_exporter))
    trace.set_tracer_provider(tracer_provider)

    # ---- Metrics ----
    metric_readers = [PrometheusMetricReader()]
    if otel_endpoint:
        export_interval = int(os.environ.get("OTEL_METRIC_EXPORT_INTERVAL", 30000))
        otlp_metric_exporter = OTLPMetricExporter(endpoint=otel_endpoint, insecure=True)
        otlp_reader = PeriodicExportingMetricReader(
            otlp_metric_exporter, export_interval_millis=export_interval
        )
        metric_readers.append(otlp_reader)

    meter_provider = MeterProvider(resource=resource, metric_readers=metric_readers)
    metrics.set_meter_provider(meter_provider)

    # ---- Log export ----
    if otel_endpoint:
        log_exporter = OTLPLogExporter(endpoint=otel_endpoint, insecure=True)
        log_processor = BatchLogRecordProcessor(log_exporter)
        logger_provider = LoggerProvider(resource=resource)
        logger_provider.add_log_record_processor(log_processor)

        handler = LoggingHandler(level=logging.NOTSET, logger_provider=logger_provider)

        # Attach handler to global + specific loggers
        root_logger = logging.getLogger()
        root_logger.addHandler(handler)

        for name in ["django", "authentication", "users", "organizations", "authapi"]:
            logging.getLogger(name).addHandler(handler)

    # ---- Auto-instrument common libraries ----
    DjangoInstrumentor().instrument()
    RequestsInstrumentor().instrument()
    Psycopg2Instrumentor().instrument()
    LoggingInstrumentor().instrument(set_logging_format=True)

    print(f"OpenTelemetry instrumentation configured for {SERVICE_NAME_VALUE}")


try:
    auth_metrics = AuthMetrics()
except Exception as exc:
    auth_metrics = None
    print(f"[Telemetry] AuthMetrics initialization failed: {exc}")
