from rest_framework_simplejwt.serializers import TokenObtainPairSerializer

from common.jwt.tokens import CustomRefreshToken


class CustomTokenObtainPairSerializer(TokenObtainPairSerializer):
    token_class = CustomRefreshToken  # ensures kid header is included

    @classmethod
    def get_token(cls, user):
        token = super().get_token(user)

        # Replace default claim
        del token["user_id"]

        # Standard + custom claims
        token["sub"] = str(user.id)
        token["username"] = user.username
        token["email"] = user.email

        if getattr(user, "org_id", None):
            token["org_id"] = str(user.org_id)

        if hasattr(user, "roles"):
            token["roles"] = getattr(user, "roles", [])

        return token
