
APP_PKG = $(shell go list github.com/brcodingdev/stock-service/internal/...)

lint:
	@echo "Linting"
	@golint -set_exit_status $(APP_PKG)
	@golangci-lint run --timeout 3m0s

test:
	@echo "Testing "
	@go test ./... -v -count=1 -race

build:
	@echo "Building"
	@go build -o chatservice ./cmd

build-docker:
	@echo "Building docker image"
	@docker build -t stock-service:latest .

run:
	@echo "Starting stock service"
	@go run ./cmd

run-docker: build-docker
	@echo "Starting stock service with docker"
	@docker run stock-service:latest