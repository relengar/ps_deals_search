from __future__ import annotations

from typing import TYPE_CHECKING, TypedDict

from sentence_transformers import SentenceTransformer
from structlog import get_logger

if TYPE_CHECKING:
    from numpy import ndarray
    from torch import Tensor

class EmbedderConfig(TypedDict):
    cache_folder: str
    device: str
    model_name: str


class Embedder:
    def __init__(self, cfg: EmbedderConfig) -> None:
        log = get_logger()
        self.__log = log.bind(device=cfg.get('device'), model=cfg.get('model_name'))
        self.__log.info("creating embedder")
        self.__model = SentenceTransformer(cfg["model_name"], cache_folder=cfg.get('cache_folder'), device=cfg.get('device'))

    def embed(self, texts: list[str]) -> list[float]:
        resp = self.__model.encode(texts)

        return [self.__to_floats(embedding) for embedding in resp]

    def __to_floats(self, numbers: Tensor | ndarray):
        return [float(num) for num in numbers]
