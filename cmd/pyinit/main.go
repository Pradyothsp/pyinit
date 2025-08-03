package main

import (
	"fmt"
	"os"

	"github.com/Pradyothsp/pyinit/internal/commands"
)

func main() {
	cmd := commands.NewCommands()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
