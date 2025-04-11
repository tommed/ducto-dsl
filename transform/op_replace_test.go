package transform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/transform"
)

func TestReplaceOperator_Name(t *testing.T) {
	op := &transform.ReplaceOperator{}
	assert.Equal(t, "replace", op.Name())
}

func TestReplaceOperator_Validate(t *testing.T) {
	op := &transform.ReplaceOperator{}
	assert.Error(t, op.Validate(transform.Instruction{Key: "x"}))
	assert.Error(t, op.Validate(transform.Instruction{To: "x"}))
	assert.Error(t, op.Validate(transform.Instruction{Key: "x", To: "y", Match: "foo"}))
	assert.NoError(t, op.Validate(transform.Instruction{Key: "x", To: "y", Match: "foo", With: "bar"}))
	assert.NoError(t, op.Validate(transform.Instruction{Key: "x", Match: "foo", With: "bar"}))
}

func TestReplaceOperator_Apply(t *testing.T) {
	op := &transform.ReplaceOperator{}

	t.Run("basic replace", func(t *testing.T) {
		input := map[string]interface{}{"text": "abc123abc"}
		instr := transform.Instruction{
			Key: "text", To: "out", Match: "abc", With: "X",
		}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, "X123X", input["out"])
	})

	t.Run("non-string input", func(t *testing.T) {
		input := map[string]interface{}{"text": 42}
		instr := transform.Instruction{From: "text", To: "out", Key: "abc", Value: "x"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})

	t.Run("invalid match type", func(t *testing.T) {
		input := map[string]interface{}{"text": 123}
		instr := transform.Instruction{Key: "text", To: "out", Match: "e", Value: 123}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})
}

func TestRegexReplaceOperator_Name(t *testing.T) {
	op := &transform.RegexReplaceOperator{}
	assert.Equal(t, "regex_replace", op.Name())
}

func TestRegexReplaceOperator_Validate(t *testing.T) {
	op := &transform.RegexReplaceOperator{}
	assert.Error(t, op.Validate(transform.Instruction{Key: "x"}))
	assert.Error(t, op.Validate(transform.Instruction{Match: "x"}))
	assert.Error(t, op.Validate(transform.Instruction{Key: "x", Match: "y", To: "["}))
	assert.Error(t, op.Validate(transform.Instruction{Key: "x", Match: "(\\.{", With: ""}))
	assert.NoError(t, op.Validate(transform.Instruction{Key: "x", To: "y", Match: "abc", With: "z"}))
	assert.NoError(t, op.Validate(transform.Instruction{Key: "x", Match: "abc", With: "z"}))
}

func TestRegexReplaceOperator_Apply(t *testing.T) {
	op := &transform.RegexReplaceOperator{}

	t.Run("basic regex replace", func(t *testing.T) {
		input := map[string]interface{}{"msg": "error: failed"}
		instr := transform.Instruction{Key: "msg", To: "cleaned", Match: `error: `, With: ""}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, "failed", input["cleaned"])
	})

	t.Run("with capture groups", func(t *testing.T) {
		input := map[string]interface{}{"val": "123-456"}
		instr := transform.Instruction{
			Key: "val", To: "swapped",
			Match: `^(\d+)-(\d+)$`, With: `$2/$1`,
		}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, "456/123", input["swapped"])
	})

	t.Run("invalid pattern", func(t *testing.T) {
		bad := &transform.RegexReplaceOperator{}
		instr := transform.Instruction{From: "x", To: "y", Key: `[`, Value: "z"}
		err := bad.Validate(instr)
		assert.Error(t, err)
	})

	t.Run("non-string input", func(t *testing.T) {
		input := map[string]interface{}{"msg": 101}
		instr := transform.Instruction{From: "msg", To: "y", Key: "e", Value: "x"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})
}
