import json
import os
from pathlib import Path

import structlog
from cryptography.hazmat.primitives import serialization
from decouple import config
from django.conf import settings


def load_signing_key():
    logger = structlog.get_logger(__name__)

    keys_dir = Path(config("KEYS_DIR", default="/keys"))
    keymap_path = keys_dir / "keymap.json"

    try:
        keymap = json.loads(keymap_path.read_text())
        active_entry = next((v for v in keymap.values() if v.get("active")), None)
        if not active_entry:
            raise ValueError("No active signing key defined")

        private_path = keys_dir / active_entry.get("private", "private.pem")
        private_data = private_path.read_bytes()

        password = os.getenv("KEY_PASSWORD")
        key = serialization.load_pem_private_key(
            private_data,
            password=password.encode() if password else None,
        )

        return key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.TraditionalOpenSSL,
            encryption_algorithm=serialization.NoEncryption(),
        ).decode("utf-8")

    except Exception as e:
        logger.error(f"[load_signing_key] Fallback due to error: {e}")
        return getattr(settings, "SECRET_KEY", "fallback-key")
