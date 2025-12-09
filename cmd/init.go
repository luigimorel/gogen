package cmd

import (
	"github.com/urfave/cli/v2"
)

var Tag = "dev"

func App() *cli.App {
	return &cli.App{
		Name:        "gogen",
		Usage:       "Generate Golang project boilerplate",
		Description: `gogen is a CLI tool for quickly generating Go project boilerplates.`,
		Version:     Tag,
		Commands: []*cli.Command{
			NewCommand(),
			InstallCommand(),
			FrontendCommand(),
		},
	}
}
