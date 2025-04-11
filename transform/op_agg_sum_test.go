package transform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/transform"
)

func TestAggSumOperator_Validate(t *testing.T) {
	op := &transform.AggSumOperator{}
	assert.Equal(t, "agg_sum", op.Name())

	t.Run("missing from", func(t *testing.T) {
		instr := transform.Instruction{To: "x", Key: "amount"}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'from'")
	})
	t.Run("missing to", func(t *testing.T) {
		instr := transform.Instruction{From: "a", Key: "amount"}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'to'")
	})
	t.Run("missing key", func(t *testing.T) {
		instr := transform.Instruction{From: "a", To: "b"}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'key'")
	})
	t.Run("valid", func(t *testing.T) {
		instr := transform.Instruction{From: "items", To: "total", Key: "amount"}
		err := op.Validate(instr)
		assert.NoError(t, err)
	})
}

func TestAggSumOperator_Apply(t *testing.T) {
	op := &transform.AggSumOperator{}

	t.Run("basic sum", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"amount": 10.5},
				map[string]interface{}{"amount": 4},
				map[string]interface{}{"amount": 2.25},
			},
		}
		instr := transform.Instruction{
			From: "items",
			To:   "summary.total",
			Key:  "amount",
		}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 16.75, input["summary"].(map[string]interface{})["total"])
	})

	t.Run("nested path sum", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"tender": map[string]interface{}{"amount": 3.5}},
				map[string]interface{}{"tender": map[string]interface{}{"amount": 6.5}},
			},
		}
		instr := transform.Instruction{
			From: "items",
			To:   "summary.amount",
			Key:  "tender.amount",
		}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 10.0, input["summary"].(map[string]interface{})["amount"])
	})

	t.Run("missing from path", func(t *testing.T) {
		input := map[string]interface{}{
			"data": []interface{}{},
		}
		instr := transform.Instruction{From: "missing", To: "x", Key: "a"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'from' path not found")
	})

	t.Run("non-array from path", func(t *testing.T) {
		input := map[string]interface{}{"items": "not-an-array"}
		instr := transform.Instruction{From: "items", To: "x", Key: "amount"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not an array")
	})

	t.Run("ignores non-numeric fields", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"amount": "invalid"},
				map[string]interface{}{"amount": 5},
			},
		}
		instr := transform.Instruction{From: "items", To: "total", Key: "amount"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 5.0, input["total"])
	})

	t.Run("handles int64 and float32", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"amount": int64(2)},
				map[string]interface{}{"amount": float32(3.5)},
			},
		}
		instr := transform.Instruction{From: "items", To: "total", Key: "amount"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.InDelta(t, 5.5, input["total"], 0.001)
	})

	t.Run("skips items that are not maps", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				"invalid",
				map[string]interface{}{"amount": 2.5},
			},
		}
		instr := transform.Instruction{From: "items", To: "total", Key: "amount"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 2.5, input["total"])
	})

	t.Run("field missing in some items", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"value": 1},
				map[string]interface{}{"amount": 3},
			},
		}
		instr := transform.Instruction{From: "items", To: "total", Key: "amount"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 3.0, input["total"])
	})
}
