package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/luigimorel/gogen/internal"
	"github.com/urfave/cli/v2"
)

type ProjectCreator struct {
	Name              string
	ModuleName        string
	Template          string
	Router            string
	FrontendFramework string
	DirName           string
	UseTypeScript     bool
	Runtime           string
}

func NewProjectCreator(name, moduleName, template, router, frontendFramework, projectDir string, useTypeScript bool, runtime string) *ProjectCreator {
	if projectDir == "" {
		projectDir = name
	}

	return &ProjectCreator{
		Name:              name,
		ModuleName:        moduleName,
		Template:          template,
		Router:            router,
		DirName:           projectDir,
		FrontendFramework: frontendFramework,
		UseTypeScript:     useTypeScript,
		Runtime:           runtime,
	}
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "new",
		Usage: "Create a new Go project",
		Description: `Create a new Go project with proper structure and initialization.
This command will create a new directory, initialize a Go module, and create a new api project`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Project name",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "module",
				Aliases: []string{"m"},
				Usage:   "Go module path (default: project name)",
			},
			&cli.StringFlag{
				Name:    "template",
				Aliases: []string{"t"},
				Usage:   "Project template (cli, web, api)",
				Value:   "api",
			},
			&cli.StringFlag{
				Name:    "router",
				Aliases: []string{"r"},
				Usage:   "Router type for API/web projects (stdlib, chi, gorilla, httprouter)",
				Value:   "stdlib",
			},
			&cli.StringFlag{
				Name:    "frontend",
				Aliases: []string{"fe"},
				Usage:   "Frontend framework for web projects (react, vue, svelte, solidjs, angular)",
			},
			&cli.StringFlag{
				Name:  "dir",
				Usage: "Directory name for the project (default: project name)",
			},
			&cli.StringFlag{
				Name:  "runtime",
				Usage: "JavaScript runtime to use (node, deno, bun)",
				Value: "node",
			},
			&cli.BoolFlag{
				Name:  "ts",
				Usage: "Use TypeScript for frontend projects (only applicable with --frontend)",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			projectName := c.String("name")
			moduleName := c.String("module")
			template := c.String("template")
			router := c.String("router")
			frontend := c.String("frontend")
			projectDir := c.String("dir")
			useTypeScript := c.Bool("ts")
			runtime := c.String("runtime")

			// Check if runtime was explicitly set by user
			runtimeExplicitlySet := c.IsSet("runtime")
			if runtimeExplicitlySet && template != "web" {
				return fmt.Errorf("runtime flag is only applicable when template is 'web'")
			}

			creator := NewProjectCreator(projectName, moduleName, template, router, frontend, projectDir, useTypeScript, runtime)
			return creator.execute()
		},
	}
}

func (pc *ProjectCreator) execute() error {
	if err := pc.validate(); err != nil {
		return err
	}

	fmt.Printf("Creating new project '%s'...\n", pc.DirName)

	if err := pc.createProjectDirectory(); err != nil {
		return err
	}

	_, cleanup, err := pc.ChangeToProjectDirectory()
	if err != nil {
		return err
	}
	defer cleanup()

	if err := pc.initializeGoModule(); err != nil {
		return err
	}

	if err := pc.createProjectFiles(); err != nil {
		return fmt.Errorf("failed to create project files: %w", err)
	}

	pc.printNextSteps()

	return nil
}

func (pc *ProjectCreator) validate() error {
	if pc.FrontendFramework != "" && pc.Template != "web" {
		return fmt.Errorf("frontend flag is only applicable when template is 'web'")
	}

	if pc.UseTypeScript && pc.FrontendFramework == "" {
		return fmt.Errorf("TypeScript flag is only applicable when frontend is specified")
	}

	return nil
}

func (pc *ProjectCreator) createProjectDirectory() error {
	if err := os.Mkdir(pc.DirName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}
	return nil
}

func (pc *ProjectCreator) ChangeToProjectDirectory() (string, func(), error) {
	originalDir, _ := os.Getwd()
	if err := os.Chdir(pc.DirName); err != nil {
		return "", nil, fmt.Errorf("failed to change to project directory: %w", err)
	}

	cleanup := func() {
		if err := os.Chdir(originalDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to change back to original directory: %v\n", err)
		}
	}

	return originalDir, cleanup, nil
}

func (pc *ProjectCreator) initializeGoModule() error {
	if pc.Template != "web" {
		moduleName := pc.ModuleName
		if moduleName == "" {
			moduleName = pc.Name
		}
		cmd := exec.Command("go", "mod", "init", moduleName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to initialize go module: %w", err)
		}
	}
	return nil
}

func (pc *ProjectCreator) createProjectFiles() error {
	pg := internal.NewProjectGenerator()

	switch pc.Template {
	case "cli":
		return pg.CreateCLIProject(pc.Name, pc.ModuleName)
	case "web":
		return pg.CreateWebProject(pc.Name, pc.ModuleName, pc.Router, pc.FrontendFramework, pc.UseTypeScript, pc.Runtime)
	case "api":
		return pg.CreateAPIProject(pc.Name, pc.ModuleName, pc.Router)
	default:
		return fmt.Errorf("unsupported template: %s", pc.Template)
	}
}

func (pc *ProjectCreator) printNextSteps() {
	fmt.Println("\nNext steps:")
	fmt.Printf("   cd %s\n", pc.Name)

	if pc.Template == "web" {
		fmt.Println("   cd api")
		fmt.Println("   go run main.go")
		if pc.FrontendFramework != "" {
			fmt.Println("\n   # In another terminal:")
			fmt.Println("   cd frontend")
			fmt.Println("   npm run dev")
		}
	} else {
		fmt.Println("   go run main.go")
	}
}
