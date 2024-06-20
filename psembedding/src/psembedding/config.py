import os
from typing import TypedDict


class Config(TypedDict):
    nats_url: str
    nats_token: str
    nats_subject: str
    model_cache_dir: str
    model_name: str
    port: int


def load_config():
    return Config(
        nats_url=os.environ.get("NATS_URL", "nats://0.0.0.0:4222"),
        nats_subject=os.environ.get("NATS_SUBJECT", "psgames.embedding"),
        nats_token=os.environ.get("NATS_TOKEN"),
        model_cache_dir=os.environ.get("MODEL_CACHE_DIR", "./.cache"),
        model_name=os.environ.get("MODEL_NAME", "sentence-transformers/all-MiniLM-L6-v2"),
        port=int(os.environ.get('PORT', 8000))
    )
