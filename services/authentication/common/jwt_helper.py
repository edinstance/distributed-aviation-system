import json
from django.conf import settings
from decouple import config


def load_signing_key():
    """Load the RSA private key for JWT signing from the keys directory."""
    keys_dir = config("KEYS_DIR", default="/keys")
    keymap_path = f"{keys_dir}/keymap.json"
    private_key_path = f"{keys_dir}/private.pem"

    try:
        # Load keymap to find active key
        with open(keymap_path) as f:
            keymap = json.load(f)

        # Find active key (for future key rotation support)
        active_entry = next((v for v in keymap.values() if v.get("active")), None)
        if not active_entry:
            raise ValueError("No active signing key defined")

        private_key_path = f"{keys_dir}/{active_entry.get('private', 'private.pem')}"

        # Read the private key
        with open(private_key_path, "r") as f:
            key_content = f.read()

        # Ensure proper PEM format
        if key_content and not key_content.startswith("-----BEGIN"):
            raise ValueError("Invalid PEM format")

        return key_content

    except (FileNotFoundError, json.JSONDecodeError, ValueError):
        # Fallback to SECRET_KEY if keys not found or invalid
        return getattr(settings, "SECRET_KEY", "fallback-key")
