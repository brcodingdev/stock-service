# Stock Service

Stock Service is a microservice designed to facilitate the retrieval of stock prices using stock codes. It integrates with RabbitMQ to receive messages containing stock codes and makes HTTP requests to stooq.com to fetch real-time stock price data.

## Technologies
- RabbitMQ

## Requirement
- Docker
- Golang `>=` 1.20
- golangci-lint

## Setup and run
- Update .env file
## Run

```bash
# run the service
$ make run
```

### Test

```bash
# run tests
$ make test
```

## TODO
- Unit Tests
- Run inside docker container
- Improve the code with the best practices in Golang