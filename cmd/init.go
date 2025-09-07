package cmd

import (
	"github.com/urfave/cli/v2"
)

// App returns the CLI application with all commands
func App() *cli.App {
	return &cli.App{
		Name:        "gogen",
		Usage:       "Generate Golang project boilerplate",
		Description: `gogen is a CLI tool for quickly generating Go project boilerplates.`,
		Version:     "0.1.0",
		Commands: []*cli.Command{
			NewCommand(),
			InstallCommand(),
			RouterCommand(),
		},
	}
}
