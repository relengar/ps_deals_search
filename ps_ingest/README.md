# PS Ingest Service
Consumes NATS messages from crawler service and serializes them to postgres DB. It calls embedder service to receive text embeddings for vector search

## Setup
- Run dependencies as described in root README.md. This service requires NATS and postgres containers running.
- Populate the `.env` file. The values such as token must be in sync with root dependencies env. See `.env.sample`
    - `cp .env.sample .env`
- Start the service
    - Directly from go `go run main.go`
    - Trough docker first build `docker build -t ps_ingest .` and then start `docker run --env-file=.env --net=host ps_ingest`
