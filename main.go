package main

import (
	"fmt"
	"os"

	"github.com/chaostreff-flensburg/moc/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run command: %v\n", err)
		os.Exit(1)
	}
}
