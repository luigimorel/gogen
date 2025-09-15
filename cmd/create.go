package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	gentstypes "github.com/luigimorel/gogen/internal/gen-ts-types"
	"github.com/urfave/cli/v2"
)

func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:        "generate",
		Usage:       "Utilities to generate code, go or ts",
		ArgsUsage:   "<command>",
		Description: `Utilities to generate code, go or ts`,
		Subcommands: []*cli.Command{
			typeGenCommand(),
		},
	}
}

func typeGenCommand() *cli.Command {
	return &cli.Command{
		Name:  "types",
		Usage: "Generate TypeScript types from Go code",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Go package path, local dir/file, or package.Type (like go doc)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output .d.ts file or directory",
				Value:   "./go-ts-types-out",
			},
		},
		Action: func(c *cli.Context) error {
			input := c.String("input")
			output := c.String("output")

			if input == "" {
				return errors.New("missing -input: provide a go package path, local dir/file, or package.Type")
			}

			absOut, err := filepath.Abs(output)
			if err != nil {
				return fmt.Errorf("failed to resolve absolute path for output: %w", err)
			}

			err = gentstypes.Generate(input, absOut)
			if err != nil {
				return fmt.Errorf("failed to generate typescript types: %w", err)
			}

			fmt.Println("Generation complete:", absOut)
			return nil
		},
	}
}
