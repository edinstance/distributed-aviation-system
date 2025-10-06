import time
import structlog
from rest_framework import status
from rest_framework.response import Response
from rest_framework_simplejwt.views import TokenObtainPairView

from authapi.serializers.jwt import CustomTokenObtainPairSerializer
from common.telemetry.helpers import get_tracer
from common.telemetry.telemetry import auth_metrics


class Login(TokenObtainPairView):
    serializer_class = CustomTokenObtainPairSerializer

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.logger = structlog.get_logger(__name__)
        self.tracer = get_tracer()

    def post(self, request, *args, **kwargs):
        start_time = time.time()
        tenant_schema = getattr(request, 'tenant_schema', 'public')
        username = request.data.get('username', 'unknown')

        with self.tracer.start_as_current_span("user_login") as span:
            span.set_attribute("tenant", tenant_schema)
            span.set_attribute("username", username)

            self.logger.info(
                "Login attempt started",
                username=username,
                tenant=tenant_schema,
                remote_addr=request.META.get('REMOTE_ADDR', 'unknown'),
            )

            serializer = self.get_serializer(data=request.data)
            try:
                serializer.is_valid(raise_exception=True)
            except Exception as e:
                duration = time.time() - start_time

                self.logger.warning(
                    "Login failed - invalid credentials",
                    username=username,
                    error=str(e),
                    duration_seconds=round(duration, 3),
                )

                # Record failed login metrics
                auth_metrics.record_login_attempt(tenant_schema, False, duration)

                span.set_attribute("login.success", False)
                span.set_attribute("error.message", str(e))

                return Response({"error": str(e)}, status=status.HTTP_400_BAD_REQUEST)

            tokens = serializer.validated_data
            user = serializer.user
            duration = time.time() - start_time

            response_data = {
                "refresh": tokens["refresh"],
                "access": tokens["access"],
                "user": {
                    "id": user.id,
                    "username": user.username,
                    "email": user.email,
                    "org_id": str(user.org_id),
                },
            }

            # Record successful login
            self.logger.info(
                "Login successful",
                user_id=user.id,
                username=user.username,
                org_id=str(user.org_id),
                duration_seconds=round(duration, 3),
            )

            # Record metrics
            auth_metrics.record_login_attempt(tenant_schema, True, duration)
            auth_metrics.record_token_generation(tenant_schema, "access")
            auth_metrics.record_token_generation(tenant_schema, "refresh")

            span.set_attribute("login.success", True)
            span.set_attribute("user.id", str(user.id))
            span.set_attribute("user.username", user.username)

            return Response(response_data, status=status.HTTP_200_OK)