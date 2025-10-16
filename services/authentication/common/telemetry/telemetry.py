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

logger = logging.getLogger("telemetry")
logger.setLevel(logging.DEBUG)


def _get_log_level_from_env() -> int:
    # Prefer LOG_LEVEL, then OTEL_LOG_LEVEL, then DJANGO_LOG_LEVEL; default INFO
    level_str = os.environ.get("LOG_LEVEL") or os.environ.get("OTEL_LOG_LEVEL") or os.environ.get("DJANGO_LOG_LEVEL") or "INFO"
    level_str = str(level_str).upper()
    return getattr(logging, level_str, logging.INFO)


def configure_telemetry(otel_endpoint: str = None):
    effective_level = _get_log_level_from_env()
    root_logger = logging.getLogger()
    root_logger.setLevel(effective_level)
    logger.debug(f"Starting OpenTelemetry setup (log level={logging.getLevelName(effective_level)})")

    if not otel_endpoint:
        otel_endpoint = os.environ.get("OTEL_EXPORTER_OTLP_ENDPOINT")

    resource = Resource(
        attributes={
            SERVICE_NAME: os.environ.get("OTEL_SERVICE_NAME", "authentication-service"),
            SERVICE_VERSION: "1.0.0",
            "service.instance.id": os.environ.get("HOSTNAME", "unknown"),
            "deployment.environment": os.environ.get("ENVIRONMENT", "dev"),
        }
    )

    logger.debug(f"Telemetry resource attributes: {resource.attributes}")

    # ---- Tracing ----
    if otel_endpoint:
        logger.debug(f"Initializing OTLPSpanExporter with endpoint={otel_endpoint}")
        span_exporter = OTLPSpanExporter(endpoint=otel_endpoint, insecure=True)
    else:
        logger.debug("No OTLP endpoint found — using ConsoleSpanExporter")
        span_exporter = ConsoleSpanExporter()

    tracer_provider = TracerProvider(resource=resource)
    tracer_provider.add_span_processor(BatchSpanProcessor(span_exporter))
    trace.set_tracer_provider(tracer_provider)
    logger.debug("Tracer provider configured")

    # ---- Metrics ----
    if otel_endpoint:
        export_interval = int(os.environ.get("OTEL_METRIC_EXPORT_INTERVAL", 30000))
        logger.debug(f"Setting OTLP metrics exporter (interval={export_interval}ms)")
        otlp_metric_exporter = OTLPMetricExporter(endpoint=otel_endpoint, insecure=True)
        otlp_reader = PeriodicExportingMetricReader(
            otlp_metric_exporter, export_interval_millis=export_interval
        )
        meter_provider = MeterProvider(resource=resource, metric_readers=[otlp_reader])
        metrics.set_meter_provider(meter_provider)
        logger.debug("Meter provider configured")

    # ---- Log export ----
    logs_exporter_cfg = (os.environ.get("OTEL_LOGS_EXPORTER", "otlp") or "").lower()
    if otel_endpoint and logs_exporter_cfg not in ("none", "disabled", "false", "off"):
        logger.debug(f"Setting up OTLP log exporter to {otel_endpoint}")
        try:
            log_exporter = OTLPLogExporter(endpoint=otel_endpoint, insecure=True)
            log_processor = BatchLogRecordProcessor(log_exporter)
            logger_provider = LoggerProvider(resource=resource)
            logger_provider.add_log_record_processor(log_processor)

            handler = LoggingHandler(level=logging.NOTSET, logger_provider=logger_provider)
            root_logger = logging.getLogger()
            # Avoid attaching duplicate OTLP logging handlers
            if not any(isinstance(h, LoggingHandler) for h in root_logger.handlers):
                root_logger.addHandler(handler)

            for name in ["django", "authentication", "users", "organizations", "authapi"]:
                named_logger = logging.getLogger(name)
                if not any(isinstance(h, LoggingHandler) for h in named_logger.handlers):
                    named_logger.addHandler(handler)

            logger.debug("OTLP log handler attached to configured loggers")
        except Exception as exc:
            logger.exception(f"Error configuring OTLP log exporter: {exc}")
    elif not otel_endpoint:
        logger.info("OTLP log exporter not configured: no OTEL_EXPORTER_OTLP_ENDPOINT provided")
    else:
        logger.info("OTLP log exporter disabled via OTEL_LOGS_EXPORTER")

    # ---- Auto-instrument common libraries ----
    logger.debug("Instrumenting Django, Requests, Psycopg2, Logging")
    DjangoInstrumentor().instrument()
    RequestsInstrumentor().instrument()
    Psycopg2Instrumentor().instrument()
    LoggingInstrumentor().instrument(set_logging_format=True)

    logger.info(f"OpenTelemetry instrumentation fully configured for {resource.attributes[SERVICE_NAME]}")


try:
    auth_metrics = AuthMetrics()
except Exception as exc:
    auth_metrics = None
    print(f"[Telemetry] AuthMetrics initialization failed: {exc}")
