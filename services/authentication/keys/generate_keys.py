"""
Generate RSA key pairs for JWT signing if they don't already exist.
"""
import json
import os
from pathlib import Path
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.primitives import serialization

KEYS_DIR = Path(os.getenv("KEYS_DIR", "./keys"))
KEYS_DIR.mkdir(parents=True, exist_ok=True)

PRIVATE_KEY = KEYS_DIR / "private.pem"
PUBLIC_KEY = KEYS_DIR / "public.pem"
KEYMAP = KEYS_DIR / "keymap.json"

def key_exists():
    return PRIVATE_KEY.exists() and PUBLIC_KEY.exists()

def generate_pair():
    print(f"[generate_keys] Generating RSA keypair in {KEYS_DIR} …")

    key = rsa.generate_private_key(public_exponent=65537, key_size=4096)

    private_pem = key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.TraditionalOpenSSL,
        encryption_algorithm=serialization.NoEncryption(),
    )
    PRIVATE_KEY.write_bytes(private_pem)

    public_pem = key.public_key().public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    )
    PUBLIC_KEY.write_bytes(public_pem)

    meta = {"v1": {"private": "private.pem", "public": "public.pem", "active": True}}
    KEYMAP.write_text(json.dumps(meta, indent=2))
    print("[generate_keys] Done.")

if not key_exists():
    generate_pair()
else:
    print("[generate_keys] Keys already exist.")