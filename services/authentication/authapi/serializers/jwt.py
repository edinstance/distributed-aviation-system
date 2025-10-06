from rest_framework_simplejwt.serializers import TokenObtainPairSerializer

class CustomTokenObtainPairSerializer(TokenObtainPairSerializer):
    @classmethod
    def get_token(cls, user):
        token = super().get_token(user)

        token["username"] = user.username
        token["email"] = user.email

        if hasattr(user, "organization") and user.organization:
            token["org_id"] = str(user.org_id)

        if hasattr(user, "roles"):
            token["roles"] = getattr(user, "roles", [])

        return token