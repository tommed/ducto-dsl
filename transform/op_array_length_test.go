package transform_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/transform"
	"testing"
)

func TestArrayLengthOperator(t *testing.T) {
	txf := transform.New()
	prog := &transform.Program{
		Version: 1,
		OnError: "fail",
		Instructions: []transform.Instruction{
			{
				Op:   "array_length",
				From: "items",
				To:   "summary.count",
			},
		},
	}

	t.Run("basic array count", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{1, 2, 3},
		}
		out, err := txf.Apply(context.Background(), input, prog)
		assert.NoError(t, err)
		assert.Equal(t, 3, out["summary"].(map[string]interface{})["count"])
	})
}

func TestArrayLengthOperator_Apply_Unit(t *testing.T) {
	op := &transform.ArrayLengthOperator{}

	t.Run("from path not found", func(t *testing.T) {
		input := map[string]interface{}{
			"other": []interface{}{1},
		}
		instr := transform.Instruction{
			Op:   "array_length",
			From: "items",
			To:   "summary.count",
		}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "'from' path not found")
	})

	t.Run("non-array from path", func(t *testing.T) {
		input := map[string]interface{}{
			"items": "not an array",
		}
		instr := transform.Instruction{
			Op:   "array_length",
			From: "items",
			To:   "summary.count",
		}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not an array")
	})
}

func TestArrayLengthOperator_Validate(t *testing.T) {
	op := &transform.ArrayLengthOperator{}

	t.Run("missing from", func(t *testing.T) {
		instr := transform.Instruction{
			To: "summary.count",
		}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'from'")
	})

	t.Run("missing to", func(t *testing.T) {
		instr := transform.Instruction{
			From: "items",
		}
		err := op.Validate(instr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing 'to'")
	})

	t.Run("valid instruction", func(t *testing.T) {
		instr := transform.Instruction{
			From: "items",
			To:   "summary.count",
		}
		err := op.Validate(instr)
		assert.NoError(t, err)
	})
}
