from uuid import UUID

from django.db import connection
from django.http import JsonResponse
from django_tenants.utils import get_tenant_model
from rest_framework_simplejwt.exceptions import TokenError, InvalidToken
from rest_framework_simplejwt.tokens import AccessToken


class HeaderTenantMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        tenant_model = get_tenant_model()
        connection.set_schema('public')
        request.tenant = None
        request.schema_name = 'public'

        # ----- classify request type -----
        path = request.path
        JWT_REQUIRED_PREFIX = '/api/'

        public_endpoints = [
            '/api/auth/jwks.json',
            '/api/organizations/create/',
            '/api/organizations/',
            '/health',
        ]
        xorg_allowed_endpoints = [
            '/api/auth/login/',
            '/api/auth/refresh/',
        ]

        is_public = any(path.startswith(ep) for ep in public_endpoints)
        allows_xorg = any(path.startswith(ep) for ep in xorg_allowed_endpoints)

        if is_public:
            return self.get_response(request)

        org_id = None

        # Try extracting from JWT
        auth_header = request.headers.get('Authorization')
        if auth_header and auth_header.startswith('Bearer '):
            token_str = auth_header[7:]
            try:
                token = AccessToken(token_str)
                org_id = token.get('org_id')
            except (TokenError, InvalidToken):
                return JsonResponse({'error': 'Invalid JWT token'}, status=401)

        # If no JWT, allow X-Org-Id for specific endpoints
        if not org_id and allows_xorg:
            org_id = request.headers.get('X-Org-Id')

        if not org_id:
            return JsonResponse(
                {'error': 'Missing org identifier. Provide JWT with org_id or X-Org-Id header.'},
                status=400,
            )

        try:
            org_uuid = UUID(str(org_id))
        except ValueError:
            return JsonResponse({'error': f'Invalid org_id: {org_id}'}, status=400)

        try:
            tenant = tenant_model.objects.get(id=org_uuid)
            connection.set_tenant(tenant)
            request.tenant = tenant
            request.schema_name = tenant.schema_name

            request.org_id = str(org_uuid)
        except tenant_model.DoesNotExist:
            return JsonResponse(
                {'error': f'Organization with id "{org_uuid}" not found'},
                status=404,
            )

        return self.get_response(request)
