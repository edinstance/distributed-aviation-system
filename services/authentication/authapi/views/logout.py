import time

import structlog
from rest_framework import status
from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.exceptions import InvalidToken, TokenError
from rest_framework_simplejwt.tokens import RefreshToken

from common.telemetry.helpers import get_tracer


class Logout(APIView):
    """
    POST /api/auth/logout/
    Blacklists the provided refresh token (if blacklist app enabled).
    """
    permission_classes = [IsAuthenticated]

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.logger = structlog.get_logger(__name__)
        self.tracer = get_tracer()

    def post(self, request):
        start_time = time.time()
        tenant_schema = getattr(request, 'tenant_schema', 'public')
        user = request.user

        with self.tracer.start_as_current_span("user_logout") as span:
            span.set_attribute("tenant", tenant_schema)
            span.set_attribute("user.id", str(user.id))
            span.set_attribute("user.username", user.username)

            self.logger.info(
                "Logout attempt started",
                user_id=user.id,
                username=user.username,
                tenant=tenant_schema,
            )

            refresh_token = request.data.get("refresh")
            if not refresh_token:
                self.logger.warning(
                    "Logout failed - no refresh token provided",
                    user_id=user.id,
                )
                return Response(
                    {"error": "Refresh token is required"},
                    status=status.HTTP_400_BAD_REQUEST,
                )

            try:
                token = RefreshToken(refresh_token)
                token.blacklist()
                duration = time.time() - start_time

                self.logger.info(
                    "Logout successful",
                    user_id=user.id,
                    username=user.username,
                    duration_seconds=round(duration, 3),
                )

                span.set_attribute("logout.success", True)

                return Response(
                    {"message": "Successfully logged out"},
                    status=status.HTTP_200_OK,
                )
            except (InvalidToken, TokenError) as e:
                duration = time.time() - start_time

                self.logger.error(
                    "Logout failed",
                    user_id=user.id,
                    username=user.username,
                    error=str(e),
                    duration_seconds=round(duration, 3),
                    exc_info=True,
                )

                span.set_attribute("logout.success", False)
                span.set_attribute("error.message", str(e))

                return Response(
                    {"error": "Invalid or expired token"},
                    status=status.HTTP_401_UNAUTHORIZED,
                )

            except AttributeError:
                duration = time.time() - start_time
                self.logger.info(
                    "Logout completed without blacklist support",
                    user_id=user.id,
                    username=user.username,
                    duration_seconds=round(duration, 3),
                )
                span.set_attribute("logout.success", True)
                span.set_attribute("logout.blacklist_skipped", True)

            except Exception as e:
                duration = time.time() - start_time
                self.logger.error(
                    "Logout failed",
                    user_id=user.id,
                    username=user.username,
                    duration_seconds=round(duration, 3),
                    error=str(e),
                )
                span.set_attribute("logout.success", False)
                span.set_attribute("logout.blacklist_skipped", False)


        return Response(
            {"message": "Successfully logged out"},
            status=status.HTTP_200_OK,
        )
