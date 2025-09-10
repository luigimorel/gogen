package main

import (
	"fmt"
	"os"

	"github.com/luigimorel/gogen/cmd"
)

func main() {
	app := cmd.App()

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start application: %v\n", err)
		os.Exit(1)
	}
}
