from django.contrib.auth import get_user_model
from rest_framework import status
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework.views import APIView
from rest_framework_simplejwt.exceptions import TokenError, InvalidToken
from rest_framework_simplejwt.settings import api_settings
from rest_framework_simplejwt.tokens import RefreshToken

User = get_user_model()


class Refresh(APIView):
    permission_classes = [AllowAny]

    def post(self, request):
        refresh_str = request.data.get("refresh")
        if not refresh_str:
            return Response(
                {"error": "Refresh token is required"},
                status=status.HTTP_400_BAD_REQUEST,
            )

        try:
            token = RefreshToken(refresh_str)
            payload = token.payload
            user_id = payload.get("user_id")

            if not user_id:
                return Response(
                    {"error": "Invalid token payload"},
                    status=status.HTTP_401_UNAUTHORIZED,
                )

            new_access = str(token.access_token)
            data = {"access": new_access}

            if api_settings.ROTATE_REFRESH_TOKENS:
                # Attempt to blacklist old refresh
                if hasattr(token, "blacklist"):
                    try:
                        token.blacklist()
                    except Exception:
                        # If blacklist app not installed, skip silently
                        print("Blacklist not enabled — skipping invalidation.")

                try:
                    user = User.objects.get(pk=user_id)
                except User.DoesNotExist:
                    return Response(
                        {"error": "User not found"}, status=status.HTTP_404_NOT_FOUND
                    )

                new_refresh = RefreshToken.for_user(user)
                data["refresh"] = str(new_refresh)
            else:
                data["refresh"] = refresh_str

            return Response(data, status=status.HTTP_200_OK)

        except (TokenError, InvalidToken):
            return Response(
                {"error": "Invalid or expired refresh token"},
                status=status.HTTP_401_UNAUTHORIZED,
            )
        except Exception as e:
            print(f"Unexpected refresh error: {e}")
            return Response(
                {"error": f"Token refresh failed: {e}"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR,
            )
