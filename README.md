# Playstation deals search

Collects deals from playstation store and provides a semantic search uppon them. There are multiple services that facilitates this.

## Servces

### Crawler
Is a web crawler that will collect the deals from playstation store and dispatch each game to nats. This is supposed to be executed as a periodic job, not to be deplyed as a permanently running service.
 - in `./ps_crawler` directory.
 - Dependencies
    - [NATS](https://docs.nats.io/)
 - Consumed by [ps_ingest](#ingest)

### Ingest
Consumes games from crawler service and serilaizes them to database. Is dispatching synchronous requests to [embedder service](#embedder) to receive text embeddings
 - in `./ps_ingest` directory
 - Dependencies
    - [NATS](https://docs.nats.io/)
    - [embedder service](#embedder)

### Embedder
Is used to provide embedding inference from a trained model, which it downloads from [huggingface](https://huggingface.co/) on startup.
 - in `./psembedding`
 - Dependencies
    - [NATS](https://docs.nats.io/)
 - Used by
    - [ingest service](#ingest)

## Setup

To establish environment variables you can create a `resources/.env` file. See `resources/.env.sample`.

The dependencies can be initialized on local with docker.
```shell
docker compose -f resources/docker-compose.dev.yaml up -d
```

Each service than has it's own setup with it's own `.env` file. It is important to keep the env variables such as tokens in sync. See respective `README.md` files for more info.

