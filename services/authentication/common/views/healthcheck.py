import socket
import time

from rest_framework import status
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework.views import APIView


def get():
    result = {
        "status": "UP",
    }

    status_code = status.HTTP_200_OK

    return Response(result, status=status_code)


class HealthCheckView(APIView):
    permission_classes = [AllowAny]
