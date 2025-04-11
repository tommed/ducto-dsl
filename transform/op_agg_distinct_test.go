package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggDistinctOperator_Validate(t *testing.T) {
	op := &AggDistinctOperator{}
	assert.Equal(t, "agg_distinct_value", op.Name())

	t.Run("missing from", func(t *testing.T) {
		instr := Instruction{To: "x", Key: "foo"}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'from'")
	})
	t.Run("missing to", func(t *testing.T) {
		instr := Instruction{From: "x", Key: "foo"}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'to'")
	})
	t.Run("missing key", func(t *testing.T) {
		instr := Instruction{From: "x", To: "y"}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'key'")
	})
	t.Run("valid", func(t *testing.T) {
		instr := Instruction{From: "x", To: "y", Key: "foo"}
		err := op.Validate(instr)
		assert.NoError(t, err)
	})
}

func TestAggDistinctOperator_Apply(t *testing.T) {
	op := &AggDistinctOperator{}

	t.Run("extract distinct values", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"sku": "a"},
				map[string]interface{}{"sku": "b"},
				map[string]interface{}{"sku": "a"},
			},
		}
		instr := Instruction{From: "items", To: "distinct", Key: "sku"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []interface{}{"a", "b"}, input["distinct"])
	})

	t.Run("handles nested fields", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"meta": map[string]interface{}{"id": "x"}},
				map[string]interface{}{"meta": map[string]interface{}{"id": "y"}},
				map[string]interface{}{"meta": map[string]interface{}{"id": "x"}},
			},
		}
		instr := Instruction{From: "items", To: "distinct", Key: "meta.id"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []interface{}{"x", "y"}, input["distinct"])
	})

	t.Run("missing from path", func(t *testing.T) {
		input := map[string]interface{}{}
		instr := Instruction{From: "items", To: "distinct", Key: "sku"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'from' path not found")
	})

	t.Run("non-array from path", func(t *testing.T) {
		input := map[string]interface{}{"items": "not-array"}
		instr := Instruction{From: "items", To: "distinct", Key: "sku"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not an array")
	})

	t.Run("skips invalid or missing values", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				"invalid",
				map[string]interface{}{"sku": "z"},
				map[string]interface{}{},
			},
		}
		instr := Instruction{From: "items", To: "distinct", Key: "sku"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, []interface{}{"z"}, input["distinct"])
	})
}

func TestNormaliseKey(t *testing.T) {
	assert.Equal(t, "abc", normaliseKey("abc"))
	assert.Equal(t, true, normaliseKey(true))
	assert.Equal(t, float64(42), normaliseKey(42.0))
	assert.Equal(t, int(42), normaliseKey(int(42)))
	assert.Equal(t, int64(42), normaliseKey(int64(42)))
	assert.Equal(t, "[1 2 3]", normaliseKey([]int{1, 2, 3}))
	assert.Equal(t, "<invalid>", normaliseKey(nil))

	var invalid interface{} = (*int)(nil)
	assert.Equal(t, "<nil>", normaliseKey(invalid))
}
