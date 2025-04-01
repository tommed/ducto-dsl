package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RunCLI(t *testing.T) {
	const goodInput = `{"foo":"bar"}`
	const goodArrayInput = `{"items": [{}, {}]}`
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
			wantErrContains: "failed to read program",
		},
		{
			name: "invalid json instructions file",
			args: args{
				input: goodInput,
				args:  []string{"../../test/data/invalid-json.json"},
			},
			want:            1,
			wantErrContains: "failed to parse program",
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
		{
			name: "map",
			args: args{
				input: goodArrayInput,
				args:  []string{"../../examples/map.json"},
			},
			want:         0,
			wantContains: `"status": "processed"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			in, out, err := buf()

			// Assemble
			in.WriteString(tt.args.input)

			// Act
			code := RunCLI(tt.args.args, in, out, err)

			// Assert
			assert.Equal(t, tt.want, code)
			assert.Contains(t, out.String(), tt.wantContains)
			assert.Contains(t, err.String(), tt.wantErrContains)
		})
	}
}

func TestRunCLI_Lint(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "simplest",
			args: args{
				filename: "../../examples/simplest.json",
			},
			want: 0,
		},
		{
			name: "invalid_op",
			args: args{
				filename: "../../test/data/invalid-op.json",
			},
			want: 1,
		},
		{
			name: "invalid_file_path",
			args: args{
				filename: "../../test/data/invalid-file-path.json",
			},
			want: 1,
		},
		{
			name: "invalid_file_path",
			args: args{
				filename: "../../test/data/invalid-json.json",
			},
			want: 1,
		},
		{
			name: "no_file",
			args: args{},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			in, out, err := buf()
			inputArgs := []string{"lint"}
			if tt.args.filename != "" {
				inputArgs = append(inputArgs, tt.args.filename)
			}
			exitCode := RunCLI(inputArgs, in, out, err)
			assert.Equal(t, tt.want, exitCode)
		})
	}
}

func buf() (*bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	return &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
}
