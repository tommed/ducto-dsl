# ----------------------
# Configuration
# ----------------------

COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html
GO=go
LINTER=golangci-lint
LINTER_REMOTE=github.com/golangci/golangci-lint/cmd/golangci-lint@latest
LINTER_OPTS=--timeout=2m

# ----------------------
# General Targets
# ----------------------

.PHONY: all check ci lint test test-full coverage example-simplest clean example-map build-all ducto-dsl-macos ducto-dsl-windows

all: check

check: lint test-full coverage

build-all: ducto-dsl-macos ducto-dsl-windows

ci: check example-simplest example-map build-all

clean:
	@rm -f $(COVERAGE_OUT) $(COVERAGE_HTML) ducto-dsl*

# ----------------------
# Linting
# ----------------------

lint:
	@echo "==> Running linter"
	$(LINTER) run $(LINTER_OPTS)

lint-install:
	go install $(LINTER_REMOTE)

# ----------------------
# Testing
# ----------------------

test:
	@echo "==> Running short tests"
	$(GO) test -short -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

test-full:
	@echo "==> Running full tests"
	$(GO) test -coverprofile=$(COVERAGE_OUT) -covermode=atomic -v ./...
	$(GO) tool cover -func=$(COVERAGE_OUT)

coverage:
	@echo "==> Generating coverage HTML report"
	$(GO) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

# ----------------------
# CLI
# ----------------------

example-simplest:
	@echo "==> Running simplest example"
	@echo '{"foo":"bar"}' | $(GO) run ./cmd/ducto-dsl examples/simplest.json

example-map:
	@echo "==> Running map example"
	@cat test/data/input.json | $(GO) run ./cmd/ducto-dsl examples/map.json

ducto-dsl-macos:
	@echo "==> Building macOS CLI"
	$(GO) build -o ducto-dsl ./cmd/ducto-dsl

ducto-dsl-windows:
	@echo "==> Building Windows CLI"
	GOOS=windows GOARCH=amd64 $(GO) build -o ducto-dsl.exe ./cmd/ducto-dsl