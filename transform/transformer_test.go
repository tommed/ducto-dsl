package transform

import (
	"context"
	"errors"
	"github.com/tommed/ducto-dsl/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformer_Apply_NoOp(t *testing.T) {
	// Assemble
	tr := New()

	input := map[string]interface{}{"foo": "bar"}

	prog := &model.Program{
		Version: 1,
		Instructions: []model.Instruction{
			{Op: "set", Key: "greeting", Value: "hello world"},
		},
	}

	// Act
	out, err := tr.Apply(context.Background(), input, prog)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "bar", out["foo"])
	assert.Equal(t, "hello world", out["greeting"])
}

func TestTransformer_Apply_InvalidOp(t *testing.T) {
	// Assemble
	tr := New()
	input := map[string]interface{}{"foo": "bar"}

	prog := &model.Program{
		Version: 1,
		OnError: "fail",
		Instructions: []model.Instruction{
			{Op: "invalid-op", Key: "foo", Value: 2},
		},
	}

	// Act
	_, err := tr.Apply(context.Background(), input, prog)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "instruction #0: unknown operator 'invalid-op'", err.Error())
}

func TestTransformer_Apply_WrongVersion(t *testing.T) {
	type args struct {
		version int
	}
	var tests = []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				version: 1,
			},
		},
		{
			name: "negative version",
			args: args{
				version: -1,
			},
			wantErr: errors.New("program version -1 is unsupported"),
		},
		{
			name: "version too high",
			args: args{
				version: 2,
			},
			wantErr: errors.New("program version 2 is unsupported"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program := &model.Program{
				OnError: "fail",
				Version: tt.args.version,
			}
			tr := New()
			_, err := tr.Apply(context.Background(), map[string]interface{}{}, program)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestTransformer_NoProgram(t *testing.T) {
	tr := New()
	_, err := tr.Apply(context.Background(), map[string]interface{}{"foo": "bar"}, nil)
	assert.Error(t, err)
}

func TestTransformer_Apply_ErrorsReturned(t *testing.T) {
	// Assemble
	tr := New()
	input := map[string]interface{}{}

	prog := &model.Program{
		Version:      1,
		OnError:      "capture",
		Instructions: []model.Instruction{{Op: "fail", Value: "Failed on purpose"}},
	}

	// Act
	ctx := context.WithValue(context.Background(), ContextKeyDebug, true)
	output, err := tr.Apply(ctx, input, prog)

	// Assert
	assert.NoError(t, err) // should have been ignored due to OnError value
	errorList, ok := output["@dsl_errors"].([]string)
	assert.True(t, ok)
	assert.Len(t, errorList, 1)

	debugInfo, ok := output["@dsl_debug"].(map[string]interface{})
	assert.True(t, ok)
	assert.Len(t, debugInfo, 2)
}

func TestTransformer_Apply_FailOnError(t *testing.T) {
	// Assemble
	tr := New()
	input := map[string]interface{}{}

	prog := &model.Program{
		Version:      1,
		OnError:      "fail",
		Instructions: []model.Instruction{{Op: "fail", Value: "Failed on purpose"}},
	}

	// Act
	_, err := tr.Apply(context.Background(), input, prog)

	// Assert
	assert.Error(t, err)
}
