import json

from pathlib import Path

import structlog
from decouple import config


def _get_active_kid() -> str:
    logger = structlog.get_logger(__name__)

    try:
        keys_dir = Path(config("KEYS_DIR", default="/keys"))
        keymap_path = keys_dir / "keymap.json"

        if not keymap_path.exists():
            logger.warning("keymap.json not found, using default kid")
            return "default"

        with open(keymap_path) as f:
            keymap = json.load(f)

        active_entry = next(
            (kid for kid, data in keymap.items() if data.get("active")), None
        )

        if not active_entry:
            logger.warning("No active key found in keymap, using default kid")
            return "default"

        return active_entry

    except Exception as e:
        logger.error(f"Error loading active kid: {e}")
        return "default"