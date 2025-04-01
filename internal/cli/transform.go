package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tommed/ducto-dsl/internal/dsl"
	"github.com/tommed/ducto-dsl/model"
	"io"
	"os"
)

//goland:noinspection GoUnhandledErrorResult
func TransformCommand(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	progData, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Fprintf(stderr, "failed to read program file: %v\n", err)
		return 1
	}

	var prog model.Program
	if err := json.Unmarshal(progData, &prog); err != nil {
		fmt.Fprintf(stderr, "failed to parse program: %v\n", err)
		return 1
	}

	var input map[string]interface{}
	if err := json.NewDecoder(stdin).Decode(&input); err != nil {
		fmt.Fprintf(stderr, "failed to parse input json: %v\n", err)
		return 1
	}

	tr := dsl.New()
	out, err := tr.Apply(context.Background(), input, &prog)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
	return 0
}
