import time

import structlog
from rest_framework import status
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.exceptions import TokenError, InvalidToken
from rest_framework_simplejwt.tokens import AccessToken

from common.telemetry.helpers import get_tracer
from common.telemetry.telemetry import auth_metrics
from users.models import User as CustomUser


class VerifyToken(APIView):
    """
    POST /api/auth/verify-token/
    Validates an access token and returns user info if valid.
    """
    permission_classes = [AllowAny]

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.logger = structlog.get_logger(__name__)
        self.tracer = get_tracer()

    def post(self, request):
        start_time = time.time()
        tenant_schema = getattr(request, 'tenant_schema', 'public')

        with self.tracer.start_as_current_span("verify_token") as span:
            span.set_attribute("tenant", tenant_schema)

            token_str = request.data.get("token")
            if not token_str:
                self.logger.warning(
                    "Token verification failed - no token provided",
                    tenant=tenant_schema,
                )
                return Response(
                    {"error": "Token is required"},
                    status=status.HTTP_400_BAD_REQUEST,
                )

            try:
                token = AccessToken(token_str)
                user_id = token["user_id"]
                user = CustomUser.objects.get(id=user_id)
                duration = time.time() - start_time

                data = {
                    "valid": True,
                    "user_id": user.id,
                    "username": user.username,
                    "email": user.email,
                    "org_id": getattr(user, "organization_id", None),
                    "roles": getattr(user, "roles", []),
                }

                self.logger.info(
                    "Token verification successful",
                    user_id=user.id,
                    username=user.username,
                    duration_seconds=round(duration, 3),
                )

                # Record metrics
                auth_metrics.record_token_validation(tenant_schema, True, duration)

                span.set_attribute("token.valid", True)
                span.set_attribute("user.id", str(user.id))
                span.set_attribute("user.username", user.username)

                return Response(data, status=status.HTTP_200_OK)

            except (InvalidToken, TokenError) as e:
                duration = time.time() - start_time

                self.logger.warning(
                    "Token verification failed - invalid token",
                    error=str(e),
                    duration_seconds=round(duration, 3),
                )

                # Record metrics
                auth_metrics.record_token_validation(tenant_schema, False, duration)

                span.set_attribute("token.valid", False)
                span.set_attribute("error.message", str(e))

                return Response({"valid": False}, status=status.HTTP_401_UNAUTHORIZED)

            except CustomUser.DoesNotExist:
                duration = time.time() - start_time

                self.logger.error(
                    "Token verification failed - user not found",
                    user_id=user_id,
                    duration_seconds=round(duration, 3),
                )

                # Record metrics
                auth_metrics.record_token_validation(tenant_schema, False, duration)

                span.set_attribute("token.valid", False)
                span.set_attribute("error.message", "User not found")

                return Response({"valid": False}, status=status.HTTP_404_NOT_FOUND)
