package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RunCLI(t *testing.T) {
	const goodInput = `{"foo":"bar"}`
	type args struct {
		input string
		args  []string
	}
	tests := []struct {
		name            string
		args            args
		want            int
		wantContains    string
		wantErrContains string
	}{
		{
			name: "simplest",
			args: args{
				input: goodInput,
				args:  []string{"../../examples/simplest.json"},
			},
			want:         0,
			wantContains: `"foo": "bar"`,
		},
		{
			name: "no instructions file",
			args: args{
				input: goodInput,
				args:  []string{},
			},
			want:            1,
			wantErrContains: "usage:",
		},
		{
			name: "invalid instructions file",
			args: args{
				input: goodInput,
				args:  []string{"./does-not-exist.json"},
			},
			want:            1,
			wantErrContains: "failed to read program file",
		},
		{
			name: "invalid json instructions file",
			args: args{
				input: goodInput,
				args:  []string{"../../invalid-json.json"},
			},
			want:            1,
			wantErrContains: "failed to read program file",
		},
		{
			name: "invalid json input piped in",
			args: args{
				input: `{`,
				args:  []string{"../../examples/simplest.json"},
			},
			want:            1,
			wantErrContains: "failed to parse input json",
		},
		{
			name: "invalid op",
			args: args{
				input: goodInput,
				args:  []string{"../../test/data/invalid-op.json"},
			},
			want:            1,
			wantErrContains: "error:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var in bytes.Buffer
			var out bytes.Buffer
			var err bytes.Buffer

			// Assemble
			in.WriteString(tt.args.input)

			// Act
			code := RunCLI(tt.args.args, &in, &out, &err)

			// Assert
			assert.Equal(t, tt.want, code)
			assert.Contains(t, out.String(), tt.wantContains)
			assert.Contains(t, err.String(), tt.wantErrContains)
		})
	}
}
