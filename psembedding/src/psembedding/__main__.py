from __future__ import annotations

import structlog
from dotenv import load_dotenv

from psembedding.config import load_config
from psembedding.embedder import Embedder, EmbedderConfig
from psembedding.http import Dependeneies, create_http_server
from psembedding.nats import NatsConfig, NatsSubscriber

import uvicorn

load_dotenv()

structlog.configure(processors=[structlog.processors.JSONRenderer()])
log = structlog.get_logger()

config = load_config()
embedder_cfg = EmbedderConfig(device="cpu", cache_folder=config.get('model_cache_dir'), model_name=config.get('model_name'))
embedder = Embedder(embedder_cfg)

nats_cfg = NatsConfig(url=config["nats_url"], token=config["nats_token"], subject=config["nats_subject"])
nats_client = NatsSubscriber(handler=embedder.embed, cfg=nats_cfg)

app = create_http_server(Dependeneies(nats=nats_client))


if __name__ == "__main__":
    try:
        uvicorn.run(app, host="0.0.0.0", port=config["port"])
    except KeyboardInterrupt as e:
        log = log.bind(reason=e)
        log.info("Interrupted")
    except Exception as e:
        log = log.bind(error=e)
        log.error("Server exception")
