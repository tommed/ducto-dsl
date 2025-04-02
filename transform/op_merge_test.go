package transform_test

import (
	"context"
	transform2 "github.com/tommed/ducto-dsl/transform"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeOperator_Apply(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		instr     transform2.Instruction
		want      map[string]interface{}
		wantError bool
	}{
		{
			name:  "basic merge",
			input: map[string]interface{}{"foo": 1},
			instr: transform2.Instruction{
				Op:    "merge",
				Value: map[string]interface{}{"bar": 2},
			},
			want: map[string]interface{}{"foo": 1, "bar": 2},
		},
		{
			name:  "overwrites existing",
			input: map[string]interface{}{"foo": 1},
			instr: transform2.Instruction{
				Op:    "merge",
				Value: map[string]interface{}{"foo": 999},
			},
			want: map[string]interface{}{"foo": 999},
		},
		{
			name:  "if_not_set prevents overwrite",
			input: map[string]interface{}{"foo": 1},
			instr: transform2.Instruction{
				Op:       "merge",
				IfNotSet: true,
				Value:    map[string]interface{}{"foo": 999, "bar": 2},
			},
			want: map[string]interface{}{"foo": 1, "bar": 2},
		},
		{
			name:  "if_not_set works on missing keys",
			input: map[string]interface{}{},
			instr: transform2.Instruction{
				Op:       "merge",
				IfNotSet: true,
				Value:    map[string]interface{}{"foo": 42},
			},
			want: map[string]interface{}{"foo": 42},
		},
		{
			name:      "error on missing value",
			input:     map[string]interface{}{},
			instr:     transform2.Instruction{Op: "merge"},
			wantError: true,
		},
		{
			name:      "error on non-map value",
			input:     map[string]interface{}{},
			instr:     transform2.Instruction{Op: "merge", Value: "not-a-map"},
			wantError: true,
		},
	}

	op := &transform2.MergeOperator{}
	reg := transform2.NewRegistry()
	reg.Register(&transform2.MergeOperator{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := transform2.NewExecutionContext(context.Background(), "fail")
			err := transform2.ValidateProgram(reg, &transform2.Program{
				Version:      1,
				OnError:      "fail",
				Instructions: []transform2.Instruction{tt.instr},
			})
			if err == nil {
				err = op.Apply(ctx, reg, tt.input, tt.instr)
			}
			assert.Equal(t, "merge", op.Name())
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.input)
			}
		})
	}
}
