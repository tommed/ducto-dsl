package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RunCLI_Simplest(t *testing.T) {
	var in bytes.Buffer
	var out bytes.Buffer
	var err bytes.Buffer

	// Prepare input
	in.WriteString(`{"foo":"bar"}`)

	code := RunCLI([]string{"../../examples/simplest.json"}, &in, &out, &err)

	assert.Equal(t, 0, code)
	assert.Contains(t, out.String(), `"foo": "bar"`)
}
