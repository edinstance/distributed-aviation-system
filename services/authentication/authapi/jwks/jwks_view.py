import json
from pathlib import Path

import structlog
from decouple import config
from django.http import JsonResponse, HttpResponseServerError
from jwcrypto import jwk


def jwks_view(_request):
    logger = structlog.get_logger(__name__)
    try:
        keys_dir = Path(config("KEYS_DIR", default="/keys"))
        keymap_path = keys_dir / "keymap.json"
        kid = "v1"

        public_pem_path = keys_dir / "public.pem"
        if keymap_path.exists():
            with open(keymap_path) as f:
                data = json.load(f)
                active_item = next(
                    ((k, v) for k, v in data.items() if v.get("active")), None
                )
                if active_item:
                    kid, meta = active_item
                    public_pem_path = keys_dir / meta.get("public", "public.pem")

        if not public_pem_path.exists():
            logger.error(
                "Public key file missing", path=str(public_pem_path)
            )
            return HttpResponseServerError(
                f"Public key file not found: {public_pem_path}"
            )

        with open(public_pem_path, "rb") as f:
            public_pem = f.read()

        key = jwk.JWK.from_pem(public_pem)

        key_dict = json.loads(key.export_public())
        key_dict.update({
            "use": "sig",
            "alg": "RS256",
            "kid": kid
        })

        logger.info("JWKS generated successfully", kid=kid)
        return JsonResponse({"keys": [key_dict]})

    except Exception as e:
        logger.exception("JWKS generation failed", error=str(e))
        return HttpResponseServerError(f"JWKS generation error: {str(e)}")
