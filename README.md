<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="./assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# Ducto DSL

[![CI](https://github.com/tommed/ducto-dsl/actions/workflows/ci.yml/badge.svg)](https://github.com/tommed/ducto-dsl/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/tommed/ducto-dsl/branch/main/graph/badge.svg)](https://codecov.io/gh/tommed/ducto-dsl)

---

## About

<p align="center">
  <img alt="Graphical representation of Ducto manipulate streaming data in a system of pipes" 
       src="./assets/ducto-representation-small.png"/>
</p>

`dsl-transformer` is a lightweight, embeddable data transformation engine designed for structured data (JSON, maps, structs). Transformations are defined using a DSL (JSON or HCL-based), making it suitable for use cases like:

- Event stream manipulation
- API response mutation
- ETL pipelines
- Testing tools and CLI automation

It is:
- ⚡ Minimal
- 🟣 Composable
- 🟩 Fully testable

`dsl-transformer` is a part of the larger [Ducto project](https://github.com/tommed), combining many interesting practices together including Feature Flagging.

![Topology Diagram](./assets/topology-medium.png)

## DSL Specification

See [doc/specs.md](docs/specs.md) for a specification of the DSL.

---

## Features (WIP)
- [x] CLI for local testing
- [x] Serverless runtime compatibility
- [x] Declarative `set` and `copy` operations
- [x] [Support for basics `map`, `delete`, `merge`](./docs/spec-v1.md)
- [x] Linter included for instruction validation
- [x] Conditionals
- [x] [Aggregations and Filtering](./docs/spec-aggs)
- [x] Embeddable Go SDK
- [ ] Input can be JSON or YAML
- [ ] HCL-powered syntax option
- [ ] OpenTelemetry instrumentation

Also, see our [OSS Release Checklist](./OSS_RELEASE_CHECKLIST.md).

---

## Example

### `examples/01-simplest.json`

**Purpose:** Show the simplest possible working example of a Ducto 'program'.

```json
{
  "version": 1,
  "on_error": "fail",
  "instructions": [
    {"op": "set", "key": "greeting", "value": "hello world"}
  ]
}
```

### `examples/02-enrich-log.json`

**Purpose:** Enrich incoming telemetry events with environment defaults, severity mapping, and drop test/debug data.

#### Input:
```json5
{
  "message": "Disk space low",
  "level": "warn"
}
```

The goal here is to:
- Default missing `env` to `"default-env"`. 
- Convert `"level"` into a `"severity"` field using a mapping. 
- Remove `debug_info`. 
- Set a `"processed": true` flag.

#### Program:
```json5
{
  "version": 1,
  "instructions": [
    
    // Defaults
    { "op": "coalesce", "key": "env", "value": "default-env" },
    { "op": "coalesce", "key": "level", "value": "low" },

    // Alter severity by known log levels 
    { "op": "if", "condition": { "equals": { "key": "level", "value": "warn" } }, "then": [
      { "op": "set", "key": "severity", "value": "medium" }
    ]},
    {
      "op": "if",
      "condition": {
        "or": [
          { "equals": { "key": "level", "value": "error" } },
          { "equals": { "key": "level", "value": "critical" } }
        ]
      },
      "then": [
        { "op": "set", "key": "severity", "value": "high" }
      ]
    },

    // Remove unnecessary keys
    { "op": "delete", "key": "debug_info" },

    // Not needed, but by set this value last we show we've finished processing this program
    { "op": "set", "key": "processed", "value": true }
  ]
}
```

#### Output:
```json5
{
  "message": "Disk space low",
  "level": "warn",
  "env": "default-env",
  "severity": "medium",
  "processed": true
}
```

## Install

```bash
go install github.com/tommed/ducto-dsl/cmd/ducto-dsl@latest

# Run (From a file)
ducto-dsl program.json < input.json

# Run (From piped text)
echo '{"foo": "bar"}' | ducto-dsl program.json

# Lint
ducto-dsl lint program.json
```

## From Source

### Lint

```bash
go run ./cmd/ducto-dsl lint examples/01-simplest.json
```

### Run

```bash
# Simple Example
echo '{"foo":"bar"}' | go run ./cmd/ducto-dsl examples/01-simplest.json

# Telemetry Example
go run ./cmd/ducto-dsl examples/02-enrich-log.json < test/data/input-telemetry-log.json
```

## Development

Please make sure you read our [Code of Conduct](./CODE_OF_CONDUCT.md) before engaging with this project. 

### Testing

```bash
make test         # Short tests
make test-full    # Full tests
make coverage     # Coverage report (HTML)
make lint-install # Install lint prerequisites
make lint         # Run static analysis
make clean        # Remove binaries and generated artifacts
```

### CLI

There are [Makefile](./Makefile) targets for a macOS binary and Windows binary. Or simply build all:

```bash
# Build all binaries
make build-all

# Build (macOS)
make ducto-dsl-macos

# Build (Microsoft Windows)
make ducto-dsl-windows
```

## Status

See [status.md](docs/status.md) for up-to-date CI, coverage, and project health.
Our [OSS Release Checklist](./OSS_RELEASE_CHECKLIST.md) also provides an overview of where we are with this project.

## License

- Code is all licensed under [MIT](./LICENSE)
- The Ducto name, logos and robot illustrations (and likeness) are (C) Copyright 2025 Tom Medhurst, all rights reserved.