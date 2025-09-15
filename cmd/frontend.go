package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/luigimorel/gogen/internal"
)

// Runtime constants
const (
	node = "node"
	bun  = "bun"
)

type FrontendManager struct {
	FrameworkType string
	DirName       string
	UseTypeScript bool
	Runtime       string
	UseTailwind   bool
}

func NewFrontendManager(frameworkType, dirName, runtime string, useTypeScript bool, useTailwind bool) *FrontendManager {
	return &FrontendManager{
		FrameworkType: frameworkType,
		DirName:       dirName,
		UseTypeScript: useTypeScript,
		Runtime:       runtime,
		UseTailwind:   useTailwind,
	}
}

func FrontendCommand() *cli.Command {
	return &cli.Command{
		Name:      "frontend",
		Usage:     "Add a frontend framework to your project",
		ArgsUsage: "<framework-type>",
		Description: `Add a frontend framework to your existing project.
This command will create a frontend directory with the selected framework setup.

Supported frameworks:
- react: React with Vite
- vue: Vue.js with Vite
- svelte: Svelte with Vite
- solidjs: SolidJS with Vite
- angular: Angular CLI

Supported runtimes:
- node: Node.js (default)
- bun: Bun

Usage:
  gogen frontend react
  gogen frontend vue --typescript
  gogen frontend svelte --dir client --runtime bun`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "Directory name for the frontend project",
				// Value:   "frontend",
			},
			&cli.BoolFlag{
				Name:    "typescript",
				Aliases: []string{"ts"},
				Usage:   "Use TypeScript (where supported)",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "runtime",
				Aliases: []string{"r"},
				Usage:   "JavaScript runtime to use (node, bun)",
				Value:   "node",
			},
			&cli.BoolFlag{
				Name:    "tailwind",
				Aliases: []string{"tw"},
				Usage:   "Add Tailwind CSS to the project",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() == 0 {
				return errors.New("framework type is required. Usage: gogen frontend <framework-type>")
			}

			frameworkType := c.Args().Get(0)
			dirName := c.String("dir")
			useTypeScript := c.Bool("typescript")
			runtime := c.String("runtime")
			useTailwind := c.Bool("tailwind")

			fmt.Printf("DEBUG CLI: framework=%s, dir=%s, typescript=%v, runtime=%s, tailwind=%v\n", frameworkType, dirName, useTypeScript, runtime, useTailwind)

			manager := NewFrontendManager(frameworkType, dirName, runtime, useTypeScript, useTailwind)
			return manager.execute()
		},
	}
}

func (fm *FrontendManager) execute() error {
	if err := fm.validateSetup(); err != nil {
		return err
	}

	pg := internal.NewProjectGenerator()
	if err := pg.CreateFrontendProject(fm.FrameworkType, fm.DirName, fm.UseTypeScript, fm.Runtime, fm.UseTailwind); err != nil {
		return fmt.Errorf("failed to create frontend project: %w", err)
	}

	fmt.Printf("Frontend project created in: %s\n", fm.DirName)
	fm.printInstructions()

	return nil
}

func (fm *FrontendManager) validateSetup() error {
	switch fm.Runtime {
	case node:
		if !fm.commandExists("node") {
			return errors.New("node.js is required but not installed. Please install Node.js from https://nodejs.org/")
		}
		if !fm.commandExists("npm") {
			return errors.New("npm is required but not installed. Please install npm")
		}
	case bun:
		if !fm.commandExists("bun") {
			return errors.New("bun is required but not installed. Please install Bun from https://bun.sh/")
		}
	default:
		return fmt.Errorf("unsupported runtime: %s. Supported runtimes: node, bun", fm.Runtime)
	}

	switch fm.FrameworkType {
	case "angular":
		if fm.Runtime == node && !fm.commandExists("ng") {
			fmt.Println("Angular CLI not found. Installing @angular/cli globally...")
			var cmd *exec.Cmd
			switch fm.Runtime {
			case node:
				cmd = exec.Command("npm", "install", "-g", "@angular/cli")
			case bun:
				cmd = exec.Command("bun", "add", "-g", "@angular/cli")
			}

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install Angular CLI: %w", err)
			}
		}
	case "react", "vue", "svelte", "solidjs":
	default:
		return fmt.Errorf("unsupported frontend framework: %s", fm.FrameworkType)
	}

	return nil
}

func (fm *FrontendManager) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (fm *FrontendManager) printInstructions() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Printf("Frontend setup complete! (Runtime: %s)\n", fm.Runtime)
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("\nNext steps:\n")
	fmt.Printf("   cd %s\n", fm.DirName)

	var devCommand string
	switch fm.Runtime {
	case node:
		devCommand = "npm run dev"
	case bun:
		devCommand = "bun run dev"
	default:
		devCommand = "npm run dev"
	}

	fmt.Printf("   %s\n", devCommand)
}
