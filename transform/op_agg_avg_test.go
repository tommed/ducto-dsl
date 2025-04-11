package transform_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tommed/ducto-dsl/transform"
)

func TestAggAvgOperator_Validate(t *testing.T) {
	op := &transform.AggAvgOperator{}
	assert.Equal(t, "agg_avg", op.Name())

	t.Run("missing from", func(t *testing.T) {
		instr := transform.Instruction{To: "x", Key: "a"}
		err := op.Validate(instr)
		assert.Error(t, err)
	})
	t.Run("missing to", func(t *testing.T) {
		instr := transform.Instruction{From: "x", Key: "a"}
		err := op.Validate(instr)
		assert.Error(t, err)
	})
	t.Run("missing key", func(t *testing.T) {
		instr := transform.Instruction{From: "x", To: "y"}
		err := op.Validate(instr)
		assert.Error(t, err)
	})
	t.Run("valid", func(t *testing.T) {
		instr := transform.Instruction{From: "x", To: "y", Key: "a"}
		err := op.Validate(instr)
		assert.NoError(t, err)
	})
}

func TestAggAvgOperator_Apply(t *testing.T) {
	op := &transform.AggAvgOperator{}

	items := []interface{}{
		map[string]interface{}{"v": 1},
		map[string]interface{}{"v": 2},
		map[string]interface{}{"v": 3},
		map[string]interface{}{"v": 2},
	}

	t.Run("mean is default", func(t *testing.T) {
		input := map[string]interface{}{"items": items}
		instr := transform.Instruction{From: "items", To: "summary.mean", Key: "v"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 2.0, input["summary"].(map[string]interface{})["mean"])
	})

	t.Run("median", func(t *testing.T) {
		input := map[string]interface{}{"items": items}
		instr := transform.Instruction{From: "items", To: "summary.median", Key: "v", Variant: "median"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 2.0, input["summary"].(map[string]interface{})["median"])
	})

	t.Run("mode", func(t *testing.T) {
		input := map[string]interface{}{"items": items}
		instr := transform.Instruction{From: "items", To: "summary.mode", Key: "v", Variant: "mode"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 2.0, input["summary"].(map[string]interface{})["mode"])
	})

	t.Run("non-array from path", func(t *testing.T) {
		input := map[string]interface{}{"items": 123}
		instr := transform.Instruction{From: "items", To: "x", Key: "v"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})

	t.Run("missing from path", func(t *testing.T) {
		input := map[string]interface{}{}
		instr := transform.Instruction{From: "items", To: "x", Key: "v"}
		err := op.Apply(nil, nil, input, instr)
		assert.Error(t, err)
	})

	t.Run("non-numeric fields ignored", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"v": "bad"},
				map[string]interface{}{"v": 10},
				map[string]interface{}{},
				"junk",
			},
		}
		instr := transform.Instruction{From: "items", To: "avg", Key: "v"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Equal(t, 10.0, input["avg"])
	})

	t.Run("handles float32 and int64", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"v": float32(1.5)},
				map[string]interface{}{"v": int64(2)},
			},
		}
		instr := transform.Instruction{From: "items", To: "summary.vals", Key: "v"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.InDelta(t, 1.75, input["summary"].(map[string]interface{})["vals"], 0.001)
	})

	t.Run("empty input returns nil", func(t *testing.T) {
		input := map[string]interface{}{
			"items": []interface{}{
				map[string]interface{}{"v": "invalid"},
			},
		}
		instr := transform.Instruction{From: "items", To: "none", Key: "v"}
		err := op.Apply(nil, nil, input, instr)
		assert.NoError(t, err)
		assert.Nil(t, input["none"])
	})
}
