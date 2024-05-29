from __future__ import annotations

import structlog
from dotenv import load_dotenv

from psembedding.config import load_config
from psembedding.embedder import Embedder, EmbedderConfig
from psembedding.http import Dependeneies, create_http_server
from psembedding.nats import NatsConfig, NatsSubscriber

load_dotenv()

structlog.configure(processors=[structlog.processors.JSONRenderer()])

config = load_config()
embedder_cfg = EmbedderConfig(device="cpu", cache_folder=config.get('model_cache_dir'), model_name=config.get('model_name'))
embedder = Embedder(embedder_cfg)

nats_cfg = NatsConfig(url=config["nats_url"], token=config["nats_token"], subject=config["nats_subject"])
nats_client = NatsSubscriber(handler=embedder.embed, cfg=nats_cfg)

app = create_http_server(Dependeneies(nats=nats_client))
