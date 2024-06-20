import asyncio
from typing import TypedDict

from fastapi import FastAPI
from structlog import get_logger

from psembedding.nats import NatsClient

log = get_logger()

class Dependeneies(TypedDict):
    nats: NatsClient


def __create_lifespan(deps: Dependeneies):
    async def lifespan(_: FastAPI):
        log.info("Initializing nats")
        task = asyncio.create_task(deps["nats"].start())
        yield
        log.info("Shutting down")
        await deps["nats"].close()
        task.done()

    return lifespan


def create_http_server(deps: Dependeneies):
    app = FastAPI(lifespan=__create_lifespan(deps=deps))

    @app.get("/healthz")
    def healthz():
        return "OK"

    return app
