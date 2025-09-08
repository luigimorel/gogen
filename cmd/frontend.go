package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/luigimorel/gogen/internal"
	"github.com/urfave/cli/v2"
)

type FrontendManager struct {
	FrameworkType string
	DirName       string
	UseTypeScript bool
}

func NewFrontendManager(frameworkType, dirName string, useTypeScript bool) *FrontendManager {
	return &FrontendManager{
		FrameworkType: frameworkType,
		DirName:       dirName,
		UseTypeScript: useTypeScript,
	}
}

func FrontendCommand() *cli.Command {
	return &cli.Command{
		Name:  "frontend",
		Usage: "Add a frontend framework to your project",
		Description: `Add a frontend framework to your existing project.
This command will create a frontend directory with the selected framework setup.

Supported frameworks:
- react: React with Vite
- vue: Vue.js with Vite
- svelte: Svelte with Vite
- solidjs: SolidJS with Vite
- angular: Angular CLI`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "Frontend framework type (react, vue, svelte, solidjs, angular)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "dir",
				Aliases: []string{"d"},
				Usage:   "Directory name for the frontend project",
				Value:   "frontend",
			},
			&cli.BoolFlag{
				Name:    "typescript",
				Aliases: []string{"ts"},
				Usage:   "Use TypeScript (where supported)",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			frameworkType := c.String("type")
			dirName := c.String("dir")
			useTypeScript := c.Bool("typescript")

			manager := NewFrontendManager(frameworkType, dirName, useTypeScript)
			return manager.execute()
		},
	}
}

func (fm *FrontendManager) execute() error {
	if err := fm.validateSetup(); err != nil {
		return err
	}

	pg := internal.NewProjectGenerator()
	if err := pg.CreateFrontendProject(fm.FrameworkType, fm.DirName, fm.UseTypeScript); err != nil {
		return fmt.Errorf("failed to create frontend project: %w", err)
	}

	fmt.Printf("Frontend project created in: %s\n", fm.DirName)
	fm.printInstructions()

	return nil
}

func (fm *FrontendManager) validateSetup() error {
	if !fm.commandExists("node") {
		return fmt.Errorf("node.js is required but not installed. Please install Node.js from https://nodejs.org/")
	}

	if !fm.commandExists("npm") {
		return fmt.Errorf("npm is required but not installed. Please install npm")
	}

	switch fm.FrameworkType {
	case "angular":
		if !fm.commandExists("ng") {
			fmt.Println("Angular CLI not found. Installing @angular/cli globally...")
			cmd := exec.Command("npm", "install", "-g", "@angular/cli")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install Angular CLI: %w", err)
			}
		}
	case "react", "vue", "svelte", "solidjs":
		// These frameworks are supported
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
	fmt.Println("Frontend setup complete!")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("\nNext steps:\n")
	fmt.Printf("   cd %s\n", fm.DirName)
	fmt.Printf("   npm run dev\n")
}
