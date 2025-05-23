package cli

import (
	"encoding/json"
	"fmt"
	"github.com/tommed/ducto-dsl/transform"
	"os"
)

//goland:noinspection GoUnhandledErrorResult
func LintCommand(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: ducto-dsl lint <program.json>")
		return 1
	}

	// Load file
	programFile := args[0]
	data, err := os.ReadFile(programFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", programFile, err)
		return 1
	}

	// Parse JSON
	var prog transform.Program
	if err := json.Unmarshal(data, &prog); err != nil {
		fmt.Fprintf(os.Stderr, "invalid JSON: %v\n", err)
		return 1
	}

	// Validate program
	reg := transform.NewDefaultRegistry()
	if err := transform.ValidateProgram(reg, &prog); err != nil {
		fmt.Fprintf(os.Stderr, "program validation failed: %v\n", err)
		return 1
	}

	fmt.Println("Program is valid ✅")
	return 0
}
