package transform

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExecutionContext_DefaultErrorBehaviour(t *testing.T) {
	exec := NewExecutionContext("")
	assert.Equal(t, "ignore", exec.OnError)
}

func TestNewExecutionContext_HandleError(t *testing.T) {
	type args struct {
		onError string
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess bool
	}{
		{
			name: "empty means ignore",
			args: args{
				onError: "",
			},
			wantSuccess: true, // ignored error
		},
		{
			name: "ignore means success regardless",
			args: args{
				onError: "ignore",
			},
			wantSuccess: true, // ignored error
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecutionContext(tt.args.onError)
			exec.OnError = tt.args.onError // reset to avoid defaulting behaviour
			got := exec.Handle(errors.New("not relevant"))
			expected := tt.wantSuccess
			assert.Equal(t, expected, got)
		})
	}
}
