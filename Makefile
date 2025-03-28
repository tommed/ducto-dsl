.PHONY: cli example-simplest test test-all

cli-macos:
	go build -o cli-macos ./cmd/transformer-cli
cli:
	go run ./cmd/transformer-cli

example-simplest:
	echo '{"foo":"bar"}' | go run ./cmd/transformer-cli examples/simplest.json

test:
	go test -short -v ./...

test-all:
	go test -v ./...