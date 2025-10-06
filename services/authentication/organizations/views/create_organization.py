import re
import structlog
from rest_framework.views import APIView
from rest_framework.response import Response
from rest_framework import status
from rest_framework.permissions import AllowAny
from django.contrib.auth import get_user_model
from django_tenants.utils import schema_context
from organizations.models import Organization

logger = structlog.get_logger(__name__)
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

        logger.info(
            "Received organization creation request",
            org_name=org_name,
            schema_name=schema_name,
        )

        # --- Validate input ---
        if not all([org_name, schema_name]):
            logger.warning(
                "Missing required fields for organization creation",
                org_name=org_name,
                schema_name=schema_name,
            )
            return Response(
                {"error": "Both 'name' and 'schema_name' are required."},
                status=status.HTTP_400_BAD_REQUEST,
            )

        if not re.match(r"^[a-z0-9_]+$", schema_name):
            logger.warning("Invalid schema_name format", schema_name=schema_name)
            return Response(
                {
                    "error": "schema_name must be lowercase, alphanumeric, and underscores only."
                },
                status=status.HTTP_400_BAD_REQUEST,
            )

        try:
            # --- Create Organization ---
            org, created = Organization.objects.get_or_create(
                schema_name=schema_name,
                defaults={"name": org_name},
            )

            if not created:
                logger.info("Organization already exists", schema_name=schema_name)
                return Response(
                    {"error": f"Organization schema '{schema_name}' already exists"},
                    status=status.HTTP_400_BAD_REQUEST,
                )

            org.save()
            logger.info(
                "Organization created successfully",
                org_id=org.id,
                org_name=org.name,
                schema_name=org.schema_name,
            )

            response_data = {
                "organization": {
                    "id": org.id,
                    "name": org.name,
                    "schema_name": org.schema_name,
                }
            }

            # Check if admin credentials are nested under 'admin' key
            admin_username = admin_data.get("username")
            admin_email = admin_data.get("email")
            admin_password = admin_data.get("password")

            logger.info("Extracted admin credentials:",
                       username=admin_username,
                       email=admin_email,
                       password_present=bool(admin_password))

            if admin_username and admin_email and admin_password:
                logger.info(
                    "Creating admin user",
                    username=admin_username,
                    email=admin_email,
                    org_id=org.id,
                    schema_name=schema_name,
                )

                with schema_context(schema_name):
                    logger.info("Attempting to create user in schema context", schema=schema_name)
                    try:
                        user = User.objects.create_superuser(
                            username=admin_username,
                            email=admin_email,
                            password=admin_password,
                        )
                        logger.info("User created successfully, setting org relationships")
                        user.org_id = org.id
                        user.save()
                        logger.info("User saved with org_id", user_id=user.id, org_id=user.org_id)
                    except Exception as user_creation_error:
                        logger.exception("Failed to create user", error=str(user_creation_error))
                        raise

                response_data["admin_user"] = {
                    "id": user.id,
                    "username": user.username,
                    "email": user.email,
                    "org_id": org.id,
                }

                logger.info(
                    "Admin user created successfully",
                    admin_user_id=user.id,
                    username=user.username,
                    schema_name=schema_name,
                )
            else:
                logger.debug(
                    "Admin credentials not provided — skipping admin user creation",
                    schema_name=schema_name,
                )

            return Response(response_data, status=status.HTTP_201_CREATED)

        except Exception as e:
            logger.exception(
                "Organization creation failed",
                org_name=org_name,
                schema_name=schema_name,
                error=str(e),
            )
            return Response(
                {"error": f"Failed to create organization: {str(e)}"},
                status=status.HTTP_500_INTERNAL_SERVER_ERROR,
            )