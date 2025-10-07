from rest_framework_simplejwt.serializers import TokenObtainPairSerializer

class CustomTokenObtainPairSerializer(TokenObtainPairSerializer):
    @classmethod
    def get_token(cls, user):
        token = super().get_token(user)

        token["sub"] = str(user.id)
        token["username"] = user.username
        token["email"] = user.email

        del token["user_id"]

        if hasattr(user, "org_id") and user.org_id:
            token["org_id"] = str(user.org_id)

        if hasattr(user, "roles"):
            token["roles"] = getattr(user, "roles", [])

        return token