from rest_framework_simplejwt.serializers import TokenObtainPairSerializer

from common.jwt.tokens import CustomRefreshToken


class CustomTokenObtainPairSerializer(TokenObtainPairSerializer):
    token_class = CustomRefreshToken  # ensures kid header is included

    @classmethod
    def get_token(cls, user):
        # Get refresh token with minimal payload (just user ID)
        token = super().get_token(user)

        # Store user reference for access token generation
        token._user = user

        return token
