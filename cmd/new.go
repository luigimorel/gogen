package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// NewCommand creates the new project command for the CLI
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
		},
		Action: func(c *cli.Context) error {
			projectName := c.String("name")
			moduleName := c.String("module")
			template := c.String("template")

			if moduleName == "" {
				moduleName = projectName
			}

			fmt.Printf("Creating new Go project '%s'...\n", projectName)

			if err := os.Mkdir(projectName, 0755); err != nil {
				return fmt.Errorf("failed to create project directory: %w", err)
			}

			originalDir, _ := os.Getwd()
			if err := os.Chdir(projectName); err != nil {
				return fmt.Errorf("failed to change to project directory: %w", err)
			}
			defer func() {
				if err := os.Chdir(originalDir); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to change back to original directory: %v\n", err)
				}
			}()

			fmt.Printf("Initializing Go module '%s'...\n", moduleName)
			cmd := exec.Command("go", "mod", "init", moduleName)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to initialize go module: %w", err)
			}

			if err := createProjectFiles(projectName, moduleName, template); err != nil {
				return fmt.Errorf("failed to create project files: %w", err)
			}

			fmt.Printf("Project '%s' created successfully!\n", projectName)
			fmt.Printf("Location: %s\n", filepath.Join(originalDir, projectName))
			fmt.Printf("Template: %s\n", template)
			fmt.Println("\nNext steps:")
			fmt.Printf("   cd %s\n", projectName)
			fmt.Println("   go run main.go")

			return nil
		},
	}
}

func createProjectFiles(projectName, moduleName, template string) error {
	switch template {
	case "cli":
		return createCLIProject(projectName, moduleName)
	case "web":
		return createWebProject(projectName, moduleName)
	case "api":
		return createAPIProject(projectName, moduleName)
	default:
		return fmt.Errorf("unsupported template: %s", template)
	}
}

func createCLIProject(projectName, moduleName string) error {
	mainContent := fmt.Sprintf(`package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "%s",
		Usage: "A CLI application built with gogen",
		Action: func(c *cli.Context) error {
			return cli.ShowAppHelp(c)
		},
		Commands: []*cli.Command{
			{
				Name:    "greet",
				Aliases: []string{"g"},
				Usage:   "Greet someone",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Value: "World",
						Usage: "Name to greet",
					},
				},
				Action: func(c *cli.Context) error {
					name := c.String("name")
					fmt.Printf("Hello %%s!\n", name)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
`, projectName)

	if err := os.WriteFile("main.go", []byte(mainContent), 0644); err != nil {
		return err
	}

	cmd := exec.Command("go", "get", "github.com/urfave/cli/v2")
	return cmd.Run()
}

func createWebProject(projectName, moduleName string) error {
	mainContent := `package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from ` + projectName + ` web server!")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	port := ":8080"
	fmt.Printf("Starting ` + projectName + ` web server on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
`

	return os.WriteFile("main.go", []byte(mainContent), 0644)
}

func createAPIProject(projectName, moduleName string) error {
	mainContent := fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string `+"`"+`json:"message"`+"`"+`
	Service string `+"`"+`json:"service"`+"`"+`
}

func main() {
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message: "Hello from %s API!",
			Service: "%s",
		}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := Response{
			Message: "API is healthy",
			Service: "%s",
		}
		json.NewEncoder(w).Encode(response)
	})

	port := ":8080"
	fmt.Printf("Starting %s API server on http://localhost%%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
`, projectName, projectName, projectName, projectName)

	return os.WriteFile("main.go", []byte(mainContent), 0644)
}
