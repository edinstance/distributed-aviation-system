from uuid import UUID

from django.db import connection
from django.http import JsonResponse
from django_tenants.utils import get_tenant_model


class HeaderTenantMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        tenant_model = get_tenant_model()

        # Default to public
        connection.set_schema('public')
        request.tenant = None
        request.schema_name = 'public'

        org_id = (
                request.headers.get("X-Org-Id")
                or request.headers.get("X-Org")
                or request.META.get("HTTP_X_ORG_ID")
        )

        # Add safe, public endpoints
        public_endpoints = [
            '/api/auth/jwks.json',
            '/api/organizations/create/',
            '/api/organizations/',
            '/health'
        ]

        is_public_endpoint = any(request.path.startswith(ep) for ep in public_endpoints)

        # 🔹 Make tenant required for all /api/auth/* including /login/
        if request.path.startswith('/api/') and not is_public_endpoint:
            if not org_id:
                return JsonResponse(
                    {'error': 'X-Org-Id header is required for tenant-specific requests'},
                    status=400,
                )

            try:
                org_uuid = UUID(org_id)
            except ValueError:
                return JsonResponse(
                    {'error': f'Invalid X-Org-Id header value: {org_id}'},
                    status=400,
                )

            try:
                tenant = tenant_model.objects.get(id=org_uuid)
                connection.set_tenant(tenant)
                request.tenant = tenant
                request.schema_name = tenant.schema_name
            except tenant_model.DoesNotExist:
                return JsonResponse(
                    {'error': f'Organization with id "{org_uuid}" not found'},
                    status=404,
                )

        response = self.get_response(request)
        return response
