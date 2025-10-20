from rest_framework import status
from rest_framework.decorators import api_view, permission_classes
from rest_framework.permissions import AllowAny
from rest_framework.response import Response
from rest_framework_simplejwt.tokens import RefreshToken

from organizations.models import Organization
from users.models import User as CustomUser


@api_view(['POST'])
@permission_classes([AllowAny])
def create(request):
    username = request.data.get('username')
    password = request.data.get('password')
    email = request.data.get('email')

    organization_id = getattr(request, "org_id", None)

    if not all([username, password, email]):
        return Response(
            {'error': 'Username, password and email are required'},
            status=status.HTTP_400_BAD_REQUEST
        )

    if not organization_id:
        return Response({'error': 'Organization ID is required'},
                        status=status.HTTP_400_BAD_REQUEST)

    if CustomUser.objects.filter(username=username, org_id=organization_id).exists():
        return Response({'error': 'Username already exists'},
                        status=status.HTTP_400_BAD_REQUEST)

    if CustomUser.objects.filter(email=email, org_id=organization_id).exists():
        return Response({'error': 'Email already exists'},
                        status=status.HTTP_400_BAD_REQUEST)

    # Get organization from tenant middleware (required by middleware)
    if not hasattr(request, 'tenant') or not request.tenant:
        return Response(
            {'error': 'Organization context is required for registration'},
            status=status.HTTP_400_BAD_REQUEST
        )

    try:
        organization = Organization.objects.get(id=request.tenant.id)
    except Organization.DoesNotExist:
        return Response({'error': 'Invalid organization'},
                        status=status.HTTP_400_BAD_REQUEST)

    user = CustomUser.objects.create_user(
        username=username,
        password=password,
        email=email,
        org_id=organization.id
    )

    refresh = RefreshToken.for_user(user)

    return Response({
        'user_id': user.id,
        'username': user.username,
        'email': user.email,
        'org_id': user.org_id,
        'access': str(refresh.access_token),
        'refresh': str(refresh)
    }, status=status.HTTP_201_CREATED)
