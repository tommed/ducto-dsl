package cli

import (
	"fmt"
	"io"
)

//goland:noinspection GoUnhandledErrorResult
func RunCLI(args []string, stdin io.Reader, stdout, stderr io.Writer) int {

	// Pre-Guards
	if len(args) < 1 {
		fmt.Fprintln(stderr, "usage: transformer-cli <program.json>")
		return 1
	}

	switch args[0] {
	case "lint":
		return LintCommand(args[1:])
	default:
		return TransformCommand(args, stdin, stdout, stderr)
	}
}
