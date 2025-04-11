# Aggregation Operators

This document outlines the design and specification for aggregation-related operators in `ducto-dsl`. These operators allow array traversal, summarisation, and statistical analysis over nested data structures.

They are useful for use cases such as:
- Calculating number of items in an order
- Summing values across line items
- Computing average discounts or item prices
- Extracting distinct sets of product IDs or SKUs

## Key Concepts

### Input Shape
All aggregation operators expect the **input value to be an array**, typically an array of JSON objects (map-like entries).

### Output
Each aggregation operator returns a single value, stored at the specified `to` path.

### Coercion
`ducto-dsl` uses `CoerceToMap(...)` to support iteration over any JSON-decoded map shape:
- `map[string]interface{}`
- `map[string]float64`
- `map[string]string`

All must be handled.

## Operators

### `array_length`
Counts the number of elements in an array.

```json
{
  "op": "array_length",
  "from": "order.items",
  "to": "summary.item_count"
}
```

### `agg_sum`
Sums the values of a numeric field in each object in an array.

```json
{
  "op": "agg_sum",
  "from": "order.items",
  "to": "summary.total_discount",
  "key": "discount"
}
```

### `agg_avg`
Computes average statistics over a field:
- `mean` (default)
- `median`
- `mode`

```json
{
  "op": "agg_avg",
  "variant": "mean",
  "from": "order.items",
  "to": "summary.avg_discount",
  "key": "discount"
}
```

```json
{
  "op": "agg_avg",
  "variant": "mode",
  "from": "order.items",
  "to": "summary.discount_mode",
  "key": "discount"
}
```

### `agg_distinct_value`
Returns an array of all distinct values for a given field.

```json
{
  "op": "agg_distinct_value",
  "from": "order.items",
  "to": "summary.unique_skus",
  "key": "sku"
}
```

## Filtering Support

Aggregation operations can be combined with a `filter` operator that narrows the input array based on conditions:

```json5
{
  "op": "filter",
  "from": "order.items",
  "condition": {
    "equals": { "key": "type", "value": "food" }
  },
  
  // `as` is optional, defaults to '_ctx' when omitted. Useful for 
  // filters in filters, so you can refer to all levels of arrays
  "as": "_ctx", 
  
  "then": [
    {
      "op": "agg_sum",
      "from": "_ctx",
      "to": "summary.food_discount",
      "key": "amount"
    }
  ]
}
```

This design would allow chained or grouped transformations within a conditional array subset.

## Nested objects in arrays

The implementation will also support nested objects inside the arrays like so:

```json
{
  "order": {
    "items": [
      {
        "tender": {
          "amount": 4.14
        }
      },
      {
        "tender": {
          "amount": 8.11
        }
      }
    ]
  }
}
```

You can `agg_sum` for instance like so:

```json
{
  "op": "agg_sum",
  "from": "order.items",
  "to": "summary.total_amount",
  "key": "tender.amount"
}
```

## Error Handling
- If the `from` path does not resolve to an array, the operator should fail with a helpful error
- If the `key` is missing in any item, it should be skipped (or optionally logged)
- If a field is not numeric when expected (for `agg_sum`, `agg_avg`), it should be skipped

## Naming Considerations
We may rename `agg_` to `reduce_` or use terms like `stats_` or `metrics_` based on final naming conventions. For now, we prefer `agg_` as it's short, intuitive, and maps to common analytical terminology.

## Future Extensions
- Support for aggregation over nested arrays
- Time-series aggregation (min/max per interval)
- Grouped aggregations (e.g., count by category)
- Inline expressions for value mapping before aggregation
