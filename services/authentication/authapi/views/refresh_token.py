import structlog
from django.contrib.auth import get_user_model
from rest_framework import status
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.exceptions import TokenError, InvalidToken
from rest_framework_simplejwt.settings import api_settings

from common.jwt.tokens import CustomRefreshToken

User = get_user_model()


class Refresh(APIView):
    permission_classes = [AllowAny]

    def __init__(self, **kwargs):
        super().__init__(**kwargs)
        self.logger = structlog.get_logger(type(self).__name__)

    def post(self, request):
        refresh_str = request.data.get("refresh")
        if not refresh_str:
            return Response(
                {"error": "Missing 'refresh' token."},
                status=status.HTTP_400_BAD_REQUEST,
            )

        try:
            token = CustomRefreshToken(refresh_str)
            payload = token.payload
            user_id = payload.get("sub")

            if not user_id:
                self.logger.warning("Refresh token missing 'sub' claim")
                return Response(
                    {"error": "Invalid refresh token payload."},
                    status=status.HTTP_401_UNAUTHORIZED,
                )

            new_access = str(token.access_token)

            response_data = {"access": new_access}

            if api_settings.ROTATE_REFRESH_TOKENS:
                if hasattr(token, "blacklist"):
                    try:
                        token.blacklist()
                        self.logger.info("Old refresh token blacklisted", user_id=user_id)
                    except Exception as e:
                        self.logger.warning(
                            "Refresh blacklist failed or not enabled", error=str(e)
                        )

                try:
                    user = User.objects.get(pk=user_id)
                except User.DoesNotExist:
                    return Response(
                        {"error": "User not found."},
                        status=status.HTTP_404_NOT_FOUND,
                    )

                new_refresh = CustomRefreshToken.for_user(user)
                response_data["refresh"] = str(new_refresh)
            else:
                response_data["refresh"] = refresh_str

            self.logger.info("Token refresh successful", user_id=user_id)
            return Response(response_data, status=status.HTTP_200_OK)

        except (InvalidToken, TokenError) as e:
            self.logger.warning("Invalid or expired refresh token", error=str(e))
            return Response(
                {"error": "Invalid or expired refresh token."},
                status=status.HTTP_401_UNAUTHORIZED,
            )

        except Exception as e:
            self.logger.exception("Unexpected error in token refresh", error=str(e))
            return Response(
                {"error": "Token refresh failed due to server error."},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR,
            )