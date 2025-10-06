from rest_framework.views import APIView
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework import status
from rest_framework_simplejwt.tokens import AccessToken
from rest_framework_simplejwt.exceptions import TokenError, InvalidToken
from users.models import User as CustomUser


class VerifyToken(APIView):
    """
    POST /api/auth/verify-token/
    Validates an access token and returns user info if valid.
    """
    permission_classes = [AllowAny]

    def post(self, request):
        token_str = request.data.get("token")
        if not token_str:
            return Response(
                {"error": "Token is required"},
                status=status.HTTP_400_BAD_REQUEST,
            )

        try:
            token = AccessToken(token_str)
            user_id = token["user_id"]
            user = CustomUser.objects.get(id=user_id)

            data = {
                "valid": True,
                "user_id": user.id,
                "username": user.username,
                "email": user.email,
                "org_id": getattr(user, "organization_id", None),
                "roles": getattr(user, "roles", []),
            }
            return Response(data, status=status.HTTP_200_OK)

        except (InvalidToken, TokenError):
            return Response({"valid": False}, status=status.HTTP_401_UNAUTHORIZED)
        except CustomUser.DoesNotExist:
            return Response({"valid": False}, status=status.HTTP_404_NOT_FOUND)