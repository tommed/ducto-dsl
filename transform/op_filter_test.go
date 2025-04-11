// File: ducto-dsl/transform/op_filter_test.go
package transform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/transform"
)

func TestFilterOperator_Name(t *testing.T) {
	op := &transform.FilterOperator{}
	assert.Equal(t, "filter", op.Name())
}

func TestFilterOperator_Validate(t *testing.T) {
	op := &transform.FilterOperator{}

	t.Run("missing from", func(t *testing.T) {
		instr := transform.Instruction{
			Condition: map[string]interface{}{"exists": "x"},
			Then:      []transform.Instruction{{Op: "noop"}},
		}
		err := op.Validate(instr)
		assert.Error(t, err)
	})

	t.Run("invalid condition", func(t *testing.T) {
		instr := transform.Instruction{
			From:      "items",
			Condition: map[string]interface{}{"bad": true},
			Then:      []transform.Instruction{{Op: "noop"}},
		}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid condition")
	})

	t.Run("missing then", func(t *testing.T) {
		instr := transform.Instruction{
			From:      "items",
			Condition: map[string]interface{}{"exists": "x"},
		}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "then")
	})

	t.Run("valid", func(t *testing.T) {
		instr := transform.Instruction{
			From:      "items",
			Condition: map[string]interface{}{"exists": "x"},
			Then:      []transform.Instruction{{Op: "noop"}},
		}
		err := op.Validate(instr)
		assert.NoError(t, err)
	})
}

func TestFilterOperator_Apply(t *testing.T) {
	op := &transform.FilterOperator{}
	reg := transform.NewDefaultRegistry()

	t.Run("filters and applies", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"x": true},
				map[string]interface{}{"y": false},
			},
		}
		instr := transform.Instruction{
			From:      "items",
			Condition: map[string]interface{}{"exists": "x"},
			Then: []transform.Instruction{{
				Op:    "set",
				Key:   "matched",
				Value: true,
			}},
		}
		err := op.Apply(nil, reg, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, true, input["matched"])
	})

	t.Run("no matches, no action", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"y": false},
			},
		}
		instr := transform.Instruction{
			From:      "items",
			Condition: map[string]interface{}{"exists": "x"},
			Then: []transform.Instruction{{
				Op:    "set",
				Key:   "should_not_exist",
				Value: true,
			}},
		}
		err := op.Apply(nil, reg, input, instr)
		assert.NoError(t, err)
		_, exists := input["should_not_exist"]
		assert.False(t, exists)
	})

	t.Run("invalid from path", func(t *testing.T) {
		input := map[string]interface{}{}
		instr := transform.Instruction{
			From:      "missing",
			Condition: map[string]interface{}{"exists": "x"},
			Then:      []transform.Instruction{{Op: "noop"}},
		}
		err := op.Apply(nil, reg, input, instr)
		assert.Error(t, err)
	})

	t.Run("non-array from path", func(t *testing.T) {
		input := map[string]interface{}{
			"items": 42,
		}
		instr := transform.Instruction{
			From:      "items",
			Condition: map[string]interface{}{"exists": "x"},
			Then:      []transform.Instruction{{Op: "noop"}},
		}
		err := op.Apply(nil, reg, input, instr)
		assert.Error(t, err)
	})
}
