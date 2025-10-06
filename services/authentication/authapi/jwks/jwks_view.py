import json
from pathlib import Path
from jwcrypto import jwk
from django.http import JsonResponse, HttpResponseServerError
from decouple import config


def jwks_view(_request):
    try:
        keys_dir = Path(config("KEYS_DIR", default="/keys"))
        public_key_path = keys_dir / "public.pem"
        keymap_path = keys_dir / "keymap.json"

        if not public_key_path.exists():
            return HttpResponseServerError(
                "Public key not found. Verify KEYS_DIR/public.pem exists."
            )

        with open(public_key_path, "rb") as f:
            public_pem = f.read()

        key = jwk.JWK.from_pem(public_pem)

        kid = "v1"
        public_pem_path = public_key_path
        if keymap_path.exists():
            with open(keymap_path) as f:
                data = json.load(f)
                active_item = next(
                    ((k, v) for k, v in data.items() if v.get("active")), None
                )
                if active_item:
                    kid, meta = active_item
                    public_pem_path = keys_dir / meta.get("public", "public.pem")

        with open(public_pem_path, "rb") as f:
            public_pem = f.read()
        key = jwk.JWK.from_pem(public_pem)
        key_dict = json.loads(key.export_public())
        key_dict.update({
            "use": "sig",
            "alg": "RS256",
            "kid": kid
        })

        return JsonResponse({"keys": [key_dict]})

    except Exception as e:
        return HttpResponseServerError(f"JWKS generation error: {str(e)}")