from uuid import UUID

import structlog
from django.db import connection
from django.http import JsonResponse
from django_tenants.utils import get_tenant_model
from rest_framework_simplejwt.exceptions import TokenError, InvalidToken
from rest_framework_simplejwt.tokens import AccessToken


class HeaderTenantMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response
        self.logger = structlog.get_logger(__name__)

    def __call__(self, request):
        tenant_model = get_tenant_model()
        connection.set_schema('public')
        request.tenant = None
        request.schema_name = 'public'

        # ----- classify request type -----
        path = request.path

        public_endpoints = [
            '/api/auth/jwks.json',
            '/api/organizations/create/',
            '/health',
        ]

        xorg_allowed_endpoints = [
            '/api/auth/login/',
            '/api/auth/refresh/',
        ]

        normalized_path = path.rstrip('/')
        is_public = any(normalized_path == ep.rstrip('/') for ep in public_endpoints)
        allows_xorg = any(path.startswith(ep) for ep in xorg_allowed_endpoints)

        # Public endpoints bypass tenant resolution entirely
        if is_public:
            self.logger.debug("Public endpoint bypassed", path=path)
            return self.get_response(request)

        org_id = None
        jwt_invalid = False

        # Try extracting from JWT first
        auth_header = request.headers.get('Authorization')
        if auth_header and auth_header.startswith('Bearer '):
            token_str = auth_header[7:]
            try:
                token = AccessToken(token_str)
                org_id = token.get('org_id')
                self.logger.debug("JWT parsed successfully", org_id=org_id)

            except (TokenError, InvalidToken) as e:
                jwt_invalid = True
                self.logger.warning("Invalid JWT token", error=str(e), path=path)

        # If no org_id found in jwt check the headers
        if not org_id and allows_xorg:
            self.logger.debug(
                "Using X-Org-Id fallback",
                org_id=org_id,
                path=path,
                allows_xorg=allows_xorg,
            )
            org_id = request.headers.get('X-Org-Id')

        if not org_id:
            if jwt_invalid:
                self.logger.info(
                    "JWT invalid and no X-Org-Id provided", path=path, status=401
                )
                return JsonResponse({'error': 'Invalid JWT token'}, status=401)

            self.logger.info(
                "No org_id found in request", path=path, status=400, allows_xorg=allows_xorg
            )
            return JsonResponse(
                {
                    'error': (
                        'Missing org identifier. Provide a valid JWT with org_id '
                        'or X-Org-Id header for allowed endpoints.'
                    )
                },
                status=400,
            )

        try:
            org_uuid = UUID(str(org_id))
        except ValueError:
            self.logger.warning("Invalid org_id", org_id=org_id)
            return JsonResponse({'error': f'Invalid org_id: {org_id}'}, status=400)

        try:
            tenant = tenant_model.objects.get(id=org_uuid)
            self.logger.debug(
                "Tenant resolved successfully", org_id=str(org_uuid), schema=tenant.schema_name
            )
        except tenant_model.DoesNotExist:
            self.logger.warning("Tenant not found", org_id=str(org_uuid))
            return JsonResponse(
                {'error': f'Organization with id "{org_uuid}" not found'}, status=404
            )

        # Set tenant context
        connection.set_tenant(tenant)
        request.tenant = tenant
        request.schema_name = tenant.schema_name
        request.org_id = str(org_uuid)

        self.logger.debug("Tenant context set", schema=request.schema_name)
        return self.get_response(request)
