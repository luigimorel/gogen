package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func (pg *ProjectGenerator) CreateFrontendProject(frameworkType, dirName string, useTypeScript bool) error {
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
			cmd = exec.Command("npm", "--yes", "create", "vue@latest", dirName, "--", "--ts", "--jsx", "--router", "--pinia", "--vitest", "--playwright", "--eslint", "--prettier")
		} else {
			cmd = exec.Command("npm", "--yes", "create", "vue@latest", dirName, "--", "--jsx", "--router", "--pinia", "--vitest", "--playwright", "--eslint", "--prettier")
		}

	case "svelte":
		if useTypeScript {
			cmd = exec.Command("npx", "sv", "create", dirName,
				"--template", "minimal",
				"--types", "ts",
				"--no-add-ons",
				"--install", "npm")
		} else {
			cmd = exec.Command("npx", "sv", "create", dirName,
				"--template", "minimal",
				"--types", "jsdoc",
				"--no-add-ons",
				"--install", "npm")
		}

	case "solidjs":
		if useTypeScript {
			cmd = exec.Command("npx", "--yes", "degit", "solidjs/templates/ts", dirName, "--force")
		} else {
			cmd = exec.Command("npx", "--yes", "degit", "solidjs/templates/js", dirName, "--force")
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

	return nil
}
