# Stock Service

Stock Service is a microservice designed to facilitate the retrieval of stock prices using stock codes. It integrates with RabbitMQ to receive messages containing stock codes and makes HTTP requests to stooq.com to fetch real-time stock price data.

## Technologies
- Golang
- RabbitMQ
- testify

## Requirement
- Docker
- Golang `>=` 1.20
- golangci-lint

## Setup
You need to clone the repository to run images of RabbitMQ and PostgreSQL  <br />
<b>Repo : https://github.com/brcodingdev/chat-service.git </b>

Update .env file, env vars:

```
#RabbitMQ  

RABBIT_USERNAME=guest
RABBIT_PASSWORD=guest
RABBIT_HOST=localhost
```

## Run

### build docker image (optional)

```bash
# builds an image
$ make build-docker
```

### run outside docker (optional)

```bash
# run the service
$ make run
```

### run with docker

```bash
# run inside container
$ make run-docker
```

### Test

```bash
# run tests
$ make test
```

## TODO
- More Unit Tests
- Improve the code with the best practices in Golang