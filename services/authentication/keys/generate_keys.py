"""
Generate RSA key pairs for JWT signing if they don't already exist.
"""
import os
import subprocess
import json
from pathlib import Path

KEYS_DIR = Path(os.getenv("KEYS_DIR", "./keys"))
KEYS_DIR.mkdir(parents=True, exist_ok=True)

PRIVATE_KEY = KEYS_DIR / "private.pem"
PUBLIC_KEY = KEYS_DIR / "public.pem"
KEYMAP = KEYS_DIR / "keymap.json"

def key_exists():
    return PRIVATE_KEY.exists() and PUBLIC_KEY.exists()

def generate_pair():
    print(f"[generate_keys] Generating RSA keypair in {KEYS_DIR} …")
    subprocess.run(["openssl", "genrsa", "-out", str(PRIVATE_KEY), "4096"], check=True)
    subprocess.run(
        ["openssl", "rsa", "-in", str(PRIVATE_KEY), "-pubout", "-out", str(PUBLIC_KEY)],
        check=True,
    )
    # basic map file (for kid/versioning)
    meta = {"v1": {"private": "private.pem", "public": "public.pem", "active": True}}
    with open(KEYMAP, "w") as f:
        json.dump(meta, f, indent=2)
    print("[generate_keys] Done.")

if not key_exists():
    generate_pair()
else:
    print("[generate_keys] Keys already exist.")