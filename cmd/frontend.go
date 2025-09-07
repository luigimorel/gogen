package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

// FrontendCommand creates the frontend framework command for the CLI
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

			fmt.Printf("Adding %s frontend framework...\n", getFrameworkDisplayName(frameworkType))

			if err := validateFrontendSetup(frameworkType); err != nil {
				return err
			}

			if err := createFrontendProject(frameworkType, dirName, useTypeScript); err != nil {
				return fmt.Errorf("failed to create frontend project: %w", err)
			}

			fmt.Printf("Successfully added %s\n", getFrameworkDisplayName(frameworkType))
			fmt.Printf("Frontend project created in: %s\n", dirName)
			printFrontendInstructions(dirName)

			return nil
		},
	}
}

func getFrameworkDisplayName(frameworkType string) string {
	switch frameworkType {
	case "react":
		return "React"
	case "vue":
		return "Vue.js"
	case "svelte":
		return "Svelte"
	case "solidjs":
		return "SolidJS"
	case "angular":
		return "Angular"
	default:
		return frameworkType
	}
}

func validateFrontendSetup(frameworkType string) error {
	if !commandExists("node") {
		return fmt.Errorf("node.js is required but not installed. Please install Node.js from https://nodejs.org/")
	}

	if !commandExists("npm") {
		return fmt.Errorf("npm is required but not installed. Please install npm")
	}

	switch frameworkType {
	case "angular":
		if !commandExists("ng") {
			fmt.Println("Angular CLI not found. Installing @angular/cli globally...")
			cmd := exec.Command("npm", "install", "-g", "@angular/cli")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to install Angular CLI: %w", err)
			}
		}
	case "react", "vue", "svelte", "solidjs":
	default:
		return fmt.Errorf("unsupported frontend framework: %s", frameworkType)
	}

	return nil
}

func createFrontendProject(frameworkType, dirName string, useTypeScript bool) error {
	//TODO: Remove directory if it exists?
	if _, err := os.Stat(dirName); err == nil {
		fmt.Printf("Directory %s already exists, removing...\n", dirName)
	}

	var cmd *exec.Cmd

	switch frameworkType {
	case "react":
		if useTypeScript {
			cmd = exec.Command("npm", "create", "vite@latest", dirName, "--", "--template", "react-ts")
		} else {
			cmd = exec.Command("npm", "create", "vite@latest", dirName, "--", "--template", "react")
		}

	case "vue":
		if useTypeScript {
			cmd = exec.Command("npm", "create", "vue@latest", dirName, "--", "--typescript")
		} else {
			cmd = exec.Command("npm", "create", "vue@latest", dirName)
		}

	case "svelte":
		if useTypeScript {
			cmd = exec.Command("npm", "create", "svelte@latest", dirName, "--", "--template", "skeleton", "--types", "typescript")
		} else {
			cmd = exec.Command("npm", "create", "svelte@latest", dirName, "--", "--template", "skeleton", "--types", "javascript")
		}

	case "solidjs":
		if useTypeScript {
			cmd = exec.Command("npm", "create", "solid@latest", dirName, "--", "--template", "ts")
		} else {
			cmd = exec.Command("npm", "create", "solid@latest", dirName, "--", "--template", "js")
		}

	case "angular":
		args := []string{"new", dirName, "--routing=true", "--style=css", "--skip-git=true"}
		if useTypeScript {
			args = append(args, "--strict=true")
		}
		cmd = exec.Command("ng", args...)

	default:
		return fmt.Errorf("unsupported frontend framework: %s", frameworkType)
	}

	fmt.Printf("Creating %s project...\n", getFrameworkDisplayName(frameworkType))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create %s project: %w", frameworkType, err)
	}

	// Install dependencies for non-Angular projects (Angular CLI handles this automatically)
	if frameworkType != "angular" {
		originalDir, _ := os.Getwd()
		if err := os.Chdir(dirName); err != nil {
			return fmt.Errorf("failed to change to frontend directory: %w", err)
		}
		defer func() {
			if err := os.Chdir(originalDir); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to change back to original directory: %v\n", err)
			}
		}()

		fmt.Println("Installing dependencies...")
		installCmd := exec.Command("npm", "install")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install dependencies: %w", err)
		}
	}

	if err := createEnvFile(dirName); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	return nil
}

func createEnvFile(dirName string) error {
	envContent := `
VITE_API_URL=http://localhost:8080
VITE_API_BASE_PATH=/api

# Development
VITE_NODE_ENV=development
`

	envPath := filepath.Join(dirName, ".env.example")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		fmt.Printf("Warning: failed to create .env.example: %v\n", err)
	}

	return nil
}

func printFrontendInstructions(dirName string) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("Frontend setup complete!")
	fmt.Println(strings.Repeat("=", 50))

	fmt.Printf("\nNext steps:\n")
	fmt.Printf("   cd %s\n", dirName)
	fmt.Printf("   npm run dev\n")
}
