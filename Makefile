# ----------------------
# Configuration
# ----------------------

COVERAGE_OUT=coverage.out
GO=go
LINTER=golangci-lint
LINTER_REMOTE=github.com/golangci/golangci-lint/cmd/golangci-lint@latest
LINTER_OPTS=--timeout=2m

# ----------------------
# General Targets
# ----------------------

.PHONY: all check ci lint test test-full coverage example-simplest clean

all: check

check: lint test coverage

ci: check example-simplest

clean:
	@rm -f $(COVERAGE_OUT) coverage.html

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
	$(GO) tool cover -html=$(COVERAGE_OUT) -o coverage.html

# ----------------------
# CLI
# ----------------------

example-simplest:
	@echo "==> Running simplest example"
	@echo '{"foo":"bar"}' | $(GO) run ./cmd/transformer-cli examples/simplest.json
