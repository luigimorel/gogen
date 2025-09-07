package main

import (
	"log/slog"
	"os"

	"github.com/luigimorel/gogen/cmd"
)

func main() {
	app := cmd.App()

	if err := app.Run(os.Args); err != nil {
		slog.Error("failed to start application", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
