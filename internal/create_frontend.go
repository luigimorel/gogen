package internal

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	react   = "react"
	vue     = "vue"
	svelte  = "svelte"
	solidjs = "solidjs"
	angular = "angular"
)

const (
	bun  = "bun"
	node = "node"
)

func (pg *ProjectGenerator) CreateFrontendProject(framework, dirName string, useTypeScript bool, runtime string, useTailwind bool) error {
	fmt.Printf("DEBUG: Creating frontend project with runtime: %s, framework: %s, dir: %s, typescript: %v\n", runtime, framework, dirName, useTypeScript)
	allowPrompts := "--"

	if runtime == bun {
		allowPrompts = ""
	}

	packageManager := map[string]string{
		bun:  bun,
		node: "npm",
	}[runtime]

	//TODO: Remove directory if it exists?
	if _, err := os.Stat(dirName); err == nil {
		fmt.Printf("Directory %s already exists, removing...\n", dirName)
	}

	var cmd *exec.Cmd

	switch framework {
	case react:
		template := react
		if useTypeScript {
			template = "react-ts"
		}
		cmd = pg.getCreateCommand(packageManager, "create", "vite@latest", dirName, allowPrompts, "--template", template)

	case vue:
		args := []string{"create", "vue@latest", dirName, allowPrompts, "--jsx", "--router", "--pinia", "--vitest", "--playwright", "--eslint", "--prettier"}
		if useTypeScript {
			args = append(args, "--ts")
		}
		cmd = pg.getCreateCommand(packageManager, args...)

	case svelte:
		mode := "jsdoc"
		if useTypeScript {
			mode = "ts"
		}
		cmd = pg.getSvelteCommand(packageManager, dirName, mode)

	case solidjs:
		mode := "js"
		if useTypeScript {
			mode = "ts"
		}
		cmd = pg.getSolidCommand(packageManager, dirName, mode)

	case angular:
		args := []string{"new", dirName, "--routing=true", "--style=css", "--skip-git=true", "--package-manager=" + packageManager}
		if useTypeScript {
			args = append(args, "--strict=true")
		}
		cmd = exec.Command("ng", args...)

	default:
		return fmt.Errorf("unsupported frontend framework: %s", framework)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create %s project: %w", framework, err)
	}

	// Install dependencies for non-Angular projects (Angular CLI handles this automatically)
	if framework != angular {
		if err := pg.installDependencies(runtime, dirName); err != nil {
			return err
		}
	}

	if useTailwind {
		tailwindConfig := NewTailwindConfig(framework, runtime, dirName)
		if err := tailwindConfig.InstallTailwindCSS(); err != nil {
			return fmt.Errorf("failed to install Tailwind CSS: %w", err)
		}
	}

	return nil
}

func (pg *ProjectGenerator) installDependencies(runtime, dirName string) error {
	originalDir, _ := os.Getwd()
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to return to original directory: %v\n", err)
		}
	}()

	if err := os.Chdir(dirName); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	fmt.Println("Installing dependencies...")
	installCmd := pg.getInstallCommand(runtime)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

func (pg *ProjectGenerator) getCreateCommand(runtime string, args ...string) *exec.Cmd {
	switch runtime {
	case bun:
		if args[0] == "create" {
			bunArgs := append([]string{"create"}, args[1:]...)
			return exec.Command(bun, bunArgs...)
		}
		return exec.Command(bun, args...)
	default:
		return exec.Command("npm", args...)
	}
}

func (pg *ProjectGenerator) getSvelteCommand(runtime, dirName, typeOption string) *exec.Cmd {
	switch runtime {
	case bun:
		return exec.Command("bunx", "sv", "create", dirName,
			"--template", "minimal",
			"--types", typeOption,
			"--no-add-ons",
			"--install", bun)
	default:
		return exec.Command("npx", "sv", "create", dirName,
			"--template", "minimal",
			"--types", typeOption,
			"--no-add-ons",
			"--install", "npm")
	}
}

func (pg *ProjectGenerator) getSolidCommand(runtime, dirName, template string) *exec.Cmd {
	// Validate template to prevent command injection
	if template != "js" && template != "ts" {
		template = "js" // default to safe value
	}

	templatePath := "solidjs/templates/" + template

	switch runtime {
	case bun:
		return exec.Command("bunx", "--yes", "degit", templatePath, dirName, "--force")
	default:
		return exec.Command("npx", "--yes", "degit", templatePath, dirName, "--force")
	}
}

func (pg *ProjectGenerator) getInstallCommand(runtime string) *exec.Cmd {
	switch runtime {
	case bun:
		return exec.Command(bun, "install")
	default:
		return exec.Command("npm", "install")
	}
}
