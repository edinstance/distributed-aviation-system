import uuid
from django.contrib.auth.models import AbstractUser
from django.db import models

class User(AbstractUser):
    id = models.UUIDField(primary_key=True, default=uuid.uuid4, editable=False)
    org_id = models.UUIDField(null=True, blank=True)
    roles = models.JSONField(default=list)

    @property
    def organization(self):
        """Get organization instance if org_id is set"""
        if self.org_id:
            from organizations.models import Organization
            try:
                return Organization.objects.get(id=self.org_id)
            except Organization.DoesNotExist:
                return None
        return None
