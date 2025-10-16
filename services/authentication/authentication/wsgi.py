"""
WSGI config for authentication project.

It exposes the WSGI callable as a module-level variable named ``application``.

For more information on this file, see
https://docs.djangoproject.com/en/5.2/howto/deployment/wsgi/
"""

import os

from django.core.wsgi import get_wsgi_application

from common.telemetry.telemetry import configure_telemetry

os.environ.setdefault("DJANGO_SETTINGS_MODULE", "authentication.settings")

configure_telemetry(
    otel_endpoint=os.environ.get("OTEL_EXPORTER_OTLP_ENDPOINT"),
)

application = get_wsgi_application()
