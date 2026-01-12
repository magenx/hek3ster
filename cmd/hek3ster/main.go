package main

import (
	"fmt"
	"os"

	"github.com/magenx/hek3ster/cmd/hek3ster/commands"
	"github.com/magenx/hek3ster/pkg/version"
)

func main() {
	if err := commands.Execute(version.Get()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
