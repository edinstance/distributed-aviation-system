import re
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from rest_framework.permissions import AllowAny
from django.contrib.auth import get_user_model
from django_tenants.utils import schema_context
from organizations.models import Organization

User = get_user_model()


class CreateOrganization(APIView):
    """
    POST /api/organizations/create/
    Creates a new tenant and optional admin user.
    """
    permission_classes = [AllowAny]

    def post(self, request):
        org_name = request.data.get("name")
        schema_name = (request.data.get("schema_name") or "").lower()
        admin_data = request.data.get("admin", {})

        # --- Validate input ---
        if not all([org_name, schema_name]):
            return Response(
                {"error": "Both 'name' and 'schema_name' are required."},
                status=status.HTTP_400_BAD_REQUEST,
            )

        if not re.match(r"^[a-z0-9_]+$", schema_name):
            return Response(
                {"error": "schema_name must be lowercase, alphanumeric, and underscores only."},
                status=status.HTTP_400_BAD_REQUEST,
            )

        try:
            # --- Create Organization ---
            org, created = Organization.objects.get_or_create(
                schema_name=schema_name,
                defaults={"name": org_name},
            )

            if not created:
                return Response(
                    {"error": f"Organization schema '{schema_name}' already exists"},
                    status=status.HTTP_400_BAD_REQUEST,
                )

            org.save()  # Important: creates the actual DB schema

            response_data = {
                "organization": {
                    "id": org.id,
                    "name": org.name,
                    "schema_name": org.schema_name,
                }
            }

            # --- Optional Admin User ---
            admin_username = admin_data.get("username")
            admin_email = admin_data.get("email")
            admin_password = admin_data.get("password")

            if admin_username and admin_email and admin_password:
                with schema_context(schema_name):
                    user = User.objects.create_superuser(
                        username=admin_username,
                        email=admin_email,
                        password=admin_password,
                        org_id=org.id,
                    )
                    user.org = org
                    user.save()
                    response_data["admin_user"] = {
                        "id": user.id,
                        "username": user.username,
                        "email": user.email,
                        "org_id": org.id,
                    }

            return Response(response_data, status=status.HTTP_201_CREATED)

        except Exception as e:
            print(f"[ERROR] Organization creation failed: {e}")
            return Response(
                {"error": f"Failed to create organization: {str(e)}"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR,
            )