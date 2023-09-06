
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

run:
	@echo "Starting stock service"
	@go run ./cmd