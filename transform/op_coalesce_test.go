package transform

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoalesceOperator_Validate(t *testing.T) {
	op := &CoalesceOperator{}

	tests := []struct {
		name    string
		instr   Instruction
		wantErr bool
	}{
		{
			name:    "missing key",
			instr:   Instruction{Op: "coalesce", Value: "default"},
			wantErr: true,
		},
		{
			name:    "missing value",
			instr:   Instruction{Op: "coalesce", Key: "foo"},
			wantErr: true,
		},
		{
			name:    "valid coalesce",
			instr:   Instruction{Op: "coalesce", Key: "foo", Value: "default"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := op.Validate(tt.instr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCoalesceOperator_Apply(t *testing.T) {
	op := &CoalesceOperator{}

	tests := []struct {
		name     string
		input    map[string]interface{}
		instr    Instruction
		expected map[string]interface{}
	}{
		{
			name:     "value missing - default applied",
			input:    map[string]interface{}{},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: "bar"},
			expected: map[string]interface{}{"foo": "bar"},
		},
		{
			name:     "value exists - default ignored",
			input:    map[string]interface{}{"foo": "baz"},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: "bar"},
			expected: map[string]interface{}{"foo": "baz"},
		},
		{
			name:     "string empty coalesced",
			input:    map[string]interface{}{"foo": ""},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: "bar"},
			expected: map[string]interface{}{"foo": "bar"},
		},
		{
			name:     "zero coalesced",
			input:    map[string]interface{}{"foo": 0},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: 5},
			expected: map[string]interface{}{"foo": 5},
		},
		{
			name:     "bool false is not empty",
			input:    map[string]interface{}{"foo": false},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: true},
			expected: map[string]interface{}{"foo": false},
		},
		{
			name:     "empty array",
			input:    map[string]interface{}{"foo": []interface{}{}},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: []interface{}{"1"}},
			expected: map[string]interface{}{"foo": []interface{}{"1"}},
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{"foo": map[string]interface{}{}},
			instr:    Instruction{Op: "coalesce", Key: "foo", Value: map[string]interface{}{"1": "1"}},
			expected: map[string]interface{}{"foo": map[string]interface{}{"1": "1"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewExecutionContext(context.Background(), "fail")
			require.NoError(t, op.Validate(tt.instr))
			err := op.Apply(ctx, nil, tt.input, tt.instr)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, tt.input)
		})
	}
}

func Test_isZero(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil",
			args: args{
				v: nil,
			},
			want: true,
		},
		{
			name: "empty string",
			args: args{
				v: "",
			},
			want: true,
		},
		{
			name: "non-empty string",
			args: args{
				v: "t",
			},
			want: false,
		},
		{
			name: "zero int",
			args: args{
				v: 0,
			},
			want: true,
		},
		{
			name: "zero float64",
			args: args{
				v: float64(0),
			},
			want: true,
		},
		{
			name: "zero uint",
			args: args{
				v: uint(0),
			},
			want: true,
		},
		{
			name: "false is non-zero",
			args: args{
				v: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, isZero(tt.args.v), "isZero(%v)", tt.args.v)
		})
	}
}
