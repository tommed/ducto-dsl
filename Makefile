COVERAGE_OUT=coverage.out

.PHONY: cli example-simplest test test-all coverage lint lint-install

cli-macos:
	go build -o cli-macos ./cmd/transformer-cli
cli:
	go run ./cmd/transformer-cli

example-simplest:
	echo '{"foo":"bar"}' | go run ./cmd/transformer-cli examples/simplest.json

# Run short unit tests quickly
test:
	go test -short -v ./...

# Run full test suite and measure coverage
test-all:
	go test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	go tool cover -func=$(COVERAGE_OUT)

coverage:
	go tool cover -html=$(COVERAGE_OUT) -o coverage.html

lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run --timeout=5m