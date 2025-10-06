from rest_framework import status
from rest_framework.response import Response
from rest_framework_simplejwt.views import TokenObtainPairView

from authapi.serializers.jwt import CustomTokenObtainPairSerializer


class Login(TokenObtainPairView):
    serializer_class = CustomTokenObtainPairSerializer

    def post(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        try:
            serializer.is_valid(raise_exception=True)
        except Exception as e:
            return Response({"error": str(e)}, status=status.HTTP_400_BAD_REQUEST)

        tokens = serializer.validated_data
        user = serializer.user

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

        return Response(response_data, status=status.HTTP_200_OK)