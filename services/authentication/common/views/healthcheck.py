import socket
import time

from rest_framework import status
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework.views import APIView


class HealthCheckView(APIView):
    permission_classes = [AllowAny]

    def get(self, request, *args, **kwargs):
        result = {
            "status": "UP",
        }

        return Response(result, status=status.HTTP_200_OK)