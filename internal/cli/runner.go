package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tommed/dsl-transformer/internal/model"
	"io"
	"os"

	"github.com/tommed/dsl-transformer/internal/dsl"
)

//goland:noinspection GoUnhandledErrorResult
func RunCLI(args []string, stdin io.Reader, stdout, stderr io.Writer) int {

	// Pre-Guards
	if len(args) < 1 {
		fmt.Fprintln(stderr, "usage: transformer-cli <program.json>")
		return 1
	}

	progFile := args[0]

	progData, err := os.ReadFile(progFile)
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
