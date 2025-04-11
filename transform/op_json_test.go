package transform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/transform"
)

func TestToJSONOperator_Name(t *testing.T) {
	op := &transform.ToJSONOperator{}
	assert.Equal(t, "to_json", op.Name())
}

func TestToJSONOperator_Validate(t *testing.T) {
	op := &transform.ToJSONOperator{}
	assert.Error(t, op.Validate(transform.Instruction{To: "x"}))
	assert.Error(t, op.Validate(transform.Instruction{From: "x"}))
	assert.NoError(t, op.Validate(transform.Instruction{From: "x", To: "y"}))
}

func TestToJSONOperator_Apply(t *testing.T) {
	op := &transform.ToJSONOperator{}

	t.Run("marshals value to JSON", func(t *testing.T) {
		input := map[string]interface{}{"a": map[string]interface{}{"b": 1}}
		instr := transform.Instruction{From: "a", To: "result"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"b":1}`, input["result"].(string))
	})

	t.Run("missing from path", func(t *testing.T) {
		input := map[string]interface{}{}
		instr := transform.Instruction{From: "x", To: "y"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})
}

func TestFromJSONOperator_Name(t *testing.T) {
	op := &transform.FromJSONOperator{}
	assert.Equal(t, "from_json", op.Name())
}

func TestFromJSONOperator_Validate(t *testing.T) {
	op := &transform.FromJSONOperator{}
	assert.Error(t, op.Validate(transform.Instruction{To: "x"}))
	assert.Error(t, op.Validate(transform.Instruction{From: "x"}))
	assert.NoError(t, op.Validate(transform.Instruction{From: "x", To: "y"}))
}

func TestFromJSONOperator_Apply(t *testing.T) {
	op := &transform.FromJSONOperator{}

	t.Run("parses JSON string", func(t *testing.T) {
		input := map[string]interface{}{"raw": `{"foo": 123}`}
		instr := transform.Instruction{From: "raw", To: "parsed"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		parsed := input["parsed"].(map[string]interface{})
		assert.Equal(t, float64(123), parsed["foo"])
	})

	t.Run("missing from path", func(t *testing.T) {
		input := map[string]interface{}{}
		instr := transform.Instruction{From: "nope", To: "x"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})

	t.Run("non-string input", func(t *testing.T) {
		input := map[string]interface{}{"x": 42}
		instr := transform.Instruction{From: "x", To: "y"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		input := map[string]interface{}{"raw": `not{json]`}
		instr := transform.Instruction{From: "raw", To: "bad"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})
}
