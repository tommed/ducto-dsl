package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tommed/dsl-transformer/internal/dsl"
	"os"
)

func main() {
	// Pre-guards
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintln(os.Stderr, "usage: transformer-cli <program.json>")
		os.Exit(1)
	}

	// First argument is the instructions file
	progFile := os.Args[1]
	progData, err := os.ReadFile(progFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read program file: %v\n", err)
		os.Exit(1)
	}

	// Deserialise the instruction set
	var prog dsl.Program
	if err := json.Unmarshal(progData, &prog); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse program: %v\n", err)
		os.Exit(1)
	}

	// Read input piped in
	var input map[string]interface{}
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to parse input json: %v\n", err)
		os.Exit(1)
	}

	tr := dsl.New()
	out, err := tr.Apply(context.Background(), input, &prog)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}
