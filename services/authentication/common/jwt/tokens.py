from rest_framework_simplejwt.tokens import RefreshToken, AccessToken

from common.jwt.header_mixin import HeaderKidMixin


class CustomAccessToken(HeaderKidMixin, AccessToken):

    def __str__(self) -> str:
        backend = self.get_token_backend()

        payload = self.payload.copy()
        if backend.audience:
            payload["aud"] = backend.audience
        if backend.issuer:
            payload["iss"] = backend.issuer

        return self._encode_with_header(backend, payload)


class CustomRefreshToken(HeaderKidMixin, RefreshToken):

    def __str__(self) -> str:
        backend = self.get_token_backend()

        payload = self.payload.copy()
        if backend.audience:
            payload["aud"] = backend.audience
        if backend.issuer:
            payload["iss"] = backend.issuer

        return self._encode_with_header(backend, payload)

    @property
    def access_token(self):
        """
        Return a CustomAccessToken linked to this refresh token,
        enriched with user claims if available.
        """
        access = super().access_token

        custom_access = CustomAccessToken()
        custom_access.payload = access.payload.copy()
        custom_access.current_time = access.current_time

        # Copy custom user data, if present
        if hasattr(self, "_user") and self._user:
            user = self._user
            custom_access["username"] = user.username
            custom_access["email"] = user.email

            if getattr(user, "org_id", None):
                custom_access["org_id"] = str(user.org_id)
                if getattr(user, "organization", None):
                    custom_access["org_name"] = user.organization.name


            if hasattr(user, "roles"):
                custom_access["roles"] = getattr(user, "roles", [])

        return custom_access
