package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tommed/ducto-dsl/transform"
	"io"
)

//goland:noinspection GoUnhandledErrorResult
func TransformCommand(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	prog, err := transform.LoadProgram(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "%s", err.Error())
		return 1
	}

	var input map[string]interface{}
	if err := json.NewDecoder(stdin).Decode(&input); err != nil {
		fmt.Fprintf(stderr, "failed to parse input json: %v\n", err)
		return 1
	}

	tr := transform.New()
	out, err := tr.Apply(context.Background(), input, prog)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
	return 0
}
