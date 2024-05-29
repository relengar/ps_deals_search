import json
from typing import Callable, NewType, TypedDict

import nats
import structlog
from nats.aio.msg import Msg

log = structlog.get_logger()

Handler = NewType("Handler", Callable[[list[str]], list[list[float]]])

class NatsClient:
    async def start(self):
        pass
    async def close(self):
        pass

class NatsConfig(TypedDict):
    url: str
    token: str
    subject: str


class NatsSubscriber(NatsClient):
    def __init__(self, handler: Handler, cfg: NatsConfig) -> None:
        self.__cfg = cfg
        self.__handler = handler
        self.__log = structlog.get_logger().bind(url=self.__cfg.get('url'), subject=self.__cfg.get('subject'))

    async def start(self):
        self.__log.info("connecting to nats")
        self.conn = await nats.connect(self.__cfg.get('url'), token=self.__cfg.get('token'))

        log.msg("subscribing")
        self.sub = await self.conn.subscribe(subject=self.__cfg.get('subject'), cb=self.__on_message)

    async def close(self):
        await self.sub.unsubscribe()
        await self.conn.close()

    async def __on_message(self, msg: Msg):
        log = self.__log.bind(raw=msg)
        log.info("received message")
        # TODO: proper serialiation - pydantic
        payload = json.loads(msg.data.decode())

        # TODO: Error handlign - try except, send ok: False
        log = self.__log.bind(msg=payload)
        log.info("processing message")
        res = self.__handler(payload["texts"])

        res_raw = json.dumps({"embeddings": res, "ok": True }).encode()
        await msg.respond(res_raw)
