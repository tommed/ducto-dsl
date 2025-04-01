package main

import (
	"github.com/tommed/ducto-dsl/internal/cli"
	"os"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
