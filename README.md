## Playstation store crawler
crawls trough currently discounted games in playstation store and dispatches them to some message queue.

### Run
Create a .env file with playstation domain and deals page as start url.

- `go run main.go`
or
- `docker build -t ps_crawler .`
- `docker run --env-file=.env ps_crawler`