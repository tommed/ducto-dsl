<!--suppress HtmlDeprecatedAttribute -->
<p align="right">
    <a href="https://github.com/tommed" title="See Project Ducto">
        <img src="../assets/ducto-logo-small.png" alt="A part of Project Ducto"/>
    </a>
</p>

# DSL Transformer — Specification (Draft v1)

> A declarative, embeddable DSL for transforming structured data (JSON-like objects) through composable instructions.

## 🟣 Core Concepts

### ✅ Input:
A simple map-based structure (e.g., JSON object, `map[string]interface{}`)

### ✅ Program:
A sequence of instructions executed in order.

```json5
{
  "version": 1,
  "on_error": "ignore",
  "instructions": [
    { "op": "set", "key": "foo", "value": "bar" },
    { "op": "delete", "key": "temp" }
  ]
}
```

## 🟣 Version

To use this spec, you must set `version` to `1` in your program.

## 🟣 On Error

By default, if an operation is not successful the program will continue anyway. You can change this behaviour using `on_error`.

| Value     | Behaviour                                                |
|-----------|----------------------------------------------------------|
| `ignore`  | Default behaviour, just keep going                       |
| `fail`    | Will stop on the first failure                           |
| `capture` | Will collect all errors in the array field `@dsl_errors` |

## 🟣 Operators

| Op                   | Description                                    | Parameters                                   | Notes                                                                                              |
|----------------------|------------------------------------------------|----------------------------------------------|----------------------------------------------------------------------------------------------------|
| `set`                | Assigns a value to a key                       | `key`, `value`                               | Supports nested keys via dot notation (optional)                                                   |
| `copy`               | Copies value from one key to another           | `from`, `to`                                 | Useful for reorganizing input                                                                      |
| `delete`             | Removes a key from the map                     | `key`, `regex` (bool)                        | Supports nested keys                                                                               |
| `noop`               | No-Operation (no-op) does nothing              | _(none)_                                     | Mostly for unit testing, always passes, does nothing to the input                                  |
| `fail`               | Always raises an error                         | `value` (string)                             | Useful inside `if` statement for validation, but mostly for unit tests                             |
| `map`                | Applies a sub-program to each item in an array | `key`, `then`                                | Recursively processes arrays                                                                       |
| `merge`              | Patches values in `value` to the root map      | `value` (object), `if_not_set` (bool)        | Could merge defaults or config                                                                     |
| `if`                 | Conditionally applies sub-instructions         | `condition`, `not` (bool), `then`            | See Special Case below                                                                             |
| `coalesce`           | Defaults a field value if unset                | `key`, `value`                               | Convenience for a set, if empty                                                                    |
| `to_json`            | Marshals value to a JSON string                | `from`, `to`                                 | Useful for objects, can also be used to stringify primitives                                       |
| `from_json`          | Unmarshal a JSON string to an Object           | `from`, `to`                                 | Useful for converting JSON strings into complex objects to be manipulated by other operations here |
| `filter`             | Ability to conditionally select items          | `from`, `condition`, `as` (optional), `then` | Great mixed with aggregation operators to count/sum a sub-set of values in an array                |
| `array_length`       | The number of items in an array                | `from`, `to`                                 | See [the Aggregations Spec](./spec-aggs.md) for more info                                          |
| `agg_sum`            | A sum of numeric values in an array            | `from`, `to`, `key`                          | See [the Aggregations Spec](./spec-aggs.md) for more info                                          |
| `agg_avg`            | An average of numeric values in an array       | `variant`, `from`, `to`, `key`               | Supports mean, mode, and medium. See [the Aggregations Spec](./spec-aggs.md) for more info         |
| `agg_distinct_value` | Get a list of distinct values from an array    | `from`, `to`, `key`                          | See [the Aggregations Spec](./spec-aggs.md) for more info                                          |
| `replace`            | Performs a string replace                      | `key`, `match`, `with`, `to` (optional)      | Useful for cleaning values                                                                         |
| `regex_replace`      | Performs a regular expression replacement      | `key`, `match`, `with`, `to` (optional)      | Useful for pattern matching, supports Golang's capture groups via `$1`, `$2`, etc.                 |

### 🟡 Special Cases

#### If Statement

The following conditions currently exist:

| Condition | Purpose                                                 |
|-----------|---------------------------------------------------------|
| `exists`  | Runs if the field is present and its value is non-null  |
| `equals`  | True if the field is present and equals a given value   |
| `and`     | An array of conditions which all must match             |
| `or`      | An array of conditions where at least one should match  |

## 🟣 Instruction Schema (JSON Draft)

```json
{
  "op": "set",
  "key": "string",
  "value": "any",
  
  "from": "string",
  "to": "string",
  
  "regex": "bool",
  
  "condition": "map[string]any",
  "not": "bool",
  "then": "[]Instruction",
  "if_not_set": "bool"
}
```

## 🟣 Example Program

```json5
{
  "version": 1,
  "instructions": [
    
    // Basics
    { "op": "set", "key": "status", "value": "active" },
    { "op": "copy", "from": "user.name", "to": "username" },
    
    // JSON conversion
    { "op": "to_json", "from": "user", "to": "user_json" },
    
    // Delete key and with regex
    { "op": "delete", "key": "debug" },
    { "op": "delete", "key": "test.*", "regex": true },
    
    // Conditional set
    { "op": "if", "condition": { "exists": "isAdmin" }, "then": [
      { "op": "set", "key": "role", "value": "admin" }
    ]},

    // remove domain and slash from Windows user
    { "op": "regex_replace", "match": "\\([A-Z0-9\\-_.]+)$", "with": "$1" },
  ]
}
```

## 🟣 Notes

- Version number MUST be provided or will error out
- Instructions MUST execute in order 
- Conditionals can later allow mini control flow, but currently limited to `exists`
- Nested programs are permitted for:
  - `map` (array processing)
  - `if` (conditional branching)
  - `filter` (filtering arrays by item values)
