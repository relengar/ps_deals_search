## Playstation store crawler
crawls trough currently discounted games in playstation store and dispatches them to some message queue.

### Run
You will need to have docker to create dependencies or use .env to adjust connection.

#### Setup
- create a `.env` file with playstation domain and deals page as start url. For local you can just copy the `.env.sample`
- provide a token value into the `QUEUE_TOKEN` variable inside the `.env` file
- start nats server by running `docker compose -f ./resources/docker-compose.yaml up -d`
#### Start
- Run crawler with `go run main.go`
<br> OR <br>
- `docker build -t ps_crawler .`
- `docker run --env-file=.env ps_crawler`