import jwt

from common.jwt.helpers import _get_active_kid


class HeaderKidMixin:
    def _encode_with_header(self, backend, payload: dict) -> str:
        signing_key = getattr(backend, "prepared_signing_key", backend.signing_key)
        token = jwt.encode(
            payload,
            signing_key,
            algorithm=backend.algorithm,
            headers={"kid": _get_active_kid()},
            json_encoder=backend.json_encoder,
        )
        return token.decode("utf-8") if isinstance(token, bytes) else token
