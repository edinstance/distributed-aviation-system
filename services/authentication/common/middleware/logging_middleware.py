import time
import uuid

import structlog
from django_tenants.utils import get_public_schema_name
from opentelemetry import trace

from common.telemetry.helpers import get_tracer


class EnhancedLoggingMiddleware:
    """Enhanced logging with tenant context and correlation IDs."""

    async_capable = False
    sync_capable = True

    def __init__(self, get_response):
        self.get_response = get_response
        self.logger = structlog.get_logger(__name__)
        self.tracer = get_tracer()

    def __call__(self, request):
        if request.path.startswith(("/health")):
            return self.get_response(request)

        correlation_id = str(uuid.uuid4())
        request.correlation_id = correlation_id

        tenant_obj = getattr(request, "tenant", None)
        tenant_schema = getattr(tenant_obj, "schema_name", get_public_schema_name())

        request.start_time = time.time()
        request.tenant_schema = tenant_schema

        structlog.contextvars.clear_contextvars()
        structlog.contextvars.bind_contextvars(
            correlation_id=correlation_id,
            tenant=tenant_schema,
            method=request.method,
            path=request.path,
            user_agent=request.META.get("HTTP_USER_AGENT", ""),
            remote_addr=self._get_client_ip(request),
        )

        # Start tracing span
        with self.tracer.start_as_current_span(f"{request.method} {request.path}") as span:
            span.set_attribute("http.method", request.method)
            span.set_attribute("http.url", request.build_absolute_uri())
            span.set_attribute("tenant.schema", tenant_schema)
            span.set_attribute("correlation.id", correlation_id)
            request.span = span

            self.logger.info(
                "Request started",
                method=request.method,
                path=request.path,
                tenant=tenant_schema,
                correlation_id=correlation_id,
            )

            try:
                response = self.get_response(request)
            except Exception as exc:
                self.process_exception(request, exc)
                raise

        # process_response
        if hasattr(request, "start_time"):
            duration = time.time() - request.start_time
            self.logger.info(
                "Request completed",
                status_code=getattr(response, "status_code", "unknown"),
                duration_seconds=round(duration, 3),
                correlation_id=correlation_id,
            )

            if hasattr(request, "span"):
                request.span.set_attribute("http.status_code", response.status_code)
                if hasattr(response, "content"):
                    request.span.set_attribute("http.response_size", len(response.content))

        return response

    def process_exception(self, request, exception):
        """Log exceptions with full context."""
        self.logger.error(
            "Request failed with exception",
            exception=str(exception),
            exception_type=type(exception).__name__,
            correlation_id=getattr(request, "correlation_id", "unknown"),
            exc_info=True,
        )
        if hasattr(request, "span"):
            request.span.record_exception(exception)
            request.span.set_status(
                trace.Status(trace.StatusCode.ERROR, str(exception))
            )

    def _get_client_ip(self, request):
        x_forwarded_for = request.META.get("HTTP_X_FORWARDED_FOR")
        return x_forwarded_for.split(",")[0].strip() if x_forwarded_for else request.META.get("REMOTE_ADDR")
