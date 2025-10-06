import uuid
from django_tenants.models import TenantMixin, DomainMixin
from django.db import models

class Organization(TenantMixin):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    name = models.CharField(max_length=200)
    auto_create_schema = True

    def __str__(self):
        return self.schema_name