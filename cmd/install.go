package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli/v2"
)

// InstallCommand creates the install command for the CLI
func InstallCommand() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "Install gogen CLI to system PATH",
		Description: `Install gogen CLI to system PATH for easy access from anywhere.
Supports Windows, Linux, and Nix package manager.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "method",
				Aliases: []string{"m"},
				Usage:   "Installation method (auto, binary, nix, brew)",
				Value:   "auto",
			},
			&cli.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Force reinstall even if already installed",
				Value:   false,
			},
		},
		Action: func(c *cli.Context) error {
			method := c.String("method")
			force := c.Bool("force")

			fmt.Println("Installing gogen CLI...")

			switch method {
			case "auto":
				return autoInstall(force)
			case "binary":
				return binaryInstall(force)
			case "nix":
				return nixInstall(force)
			case "brew":
				return brewInstall(force)
			default:
				return fmt.Errorf("unsupported installation method: %s", method)
			}
		},
	}
}

func autoInstall(force bool) error {
	switch runtime.GOOS {
	case "darwin":
		if commandExists("brew") {
			fmt.Println("Detected Homebrew, using brew installation...")
			return brewInstall(force)
		}
		return binaryInstall(force)
	case "linux":
		if commandExists("nix-env") || commandExists("nix") {
			fmt.Println("Detected Nix, using nix installation...")
			return nixInstall(force)
		}
		return binaryInstall(force)
	case "windows":
		return binaryInstall(force)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func binaryInstall(force bool) error {
	binDir := getBinaryInstallDir()
	binPath := filepath.Join(binDir, getBinaryName())

	if !force && fileExists(binPath) {
		fmt.Printf("gogen is already installed at %s\n", binPath)
		fmt.Println("Use --force to reinstall")
		return nil
	}

	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	downloadURL := getDownloadURL()
	fmt.Printf("Downloading from %s...\n", downloadURL)

	if err := downloadFile(downloadURL, binPath); err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(binPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}
	}

	fmt.Printf("gogen installed successfully to %s\n", binPath)
	printPathInstructions(binDir)

	return nil
}

func nixInstall(force bool) error {
	if !commandExists("nix-env") && !commandExists("nix") {
		return fmt.Errorf("nix is not installed on this system")
	}

	fmt.Println("Nix package not available yet, using binary installation...")
	return binaryInstall(force)
}

func brewInstall(force bool) error {
	if !commandExists("brew") {
		return fmt.Errorf("homebrew is not installed on this system")
	}

	fmt.Println("Homebrew formula not available yet, using binary installation...")
	return binaryInstall(force)
}

// Helper functions
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func getBinaryInstallDir() string {
	switch runtime.GOOS {
	case "windows":
		// Use %USERPROFILE%\AppData\Local\gogen
		home, _ := os.UserHomeDir()
		return filepath.Join(home, "AppData", "Local", "gogen")
	default:
		// Use ~/.local/bin for Unix systems
		home, _ := os.UserHomeDir()
		return filepath.Join(home, ".local", "bin")
	}
}

func getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "gogen.exe"
	}
	return "gogen"
}

func getDownloadURL() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	if arch == "amd64" {
		arch = "x86_64"
	}

	version := "latest"
	filename := fmt.Sprintf("gogen_%s_%s", os, arch)
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}

	return fmt.Sprintf("https://github.com/luigimorel/gogen/releases/%s/download/%s", version, filename)
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func printPathInstructions(binDir string) {
	switch runtime.GOOS {
	case "windows":
		fmt.Println("\nTo add gogen to your PATH:")
		fmt.Println("1. Press Win + R, type 'sysdm.cpl', and press Enter")
		fmt.Println("2. Click 'Environment Variables'")
		fmt.Println("3. Under 'User variables', select 'Path' and click 'Edit'")
		fmt.Printf("4. Click 'New' and add: %s\n", binDir)
		fmt.Println("5. Click OK to save changes")
		fmt.Println("6. Restart your terminal and run 'gogen --help'")
	default:
		shell := os.Getenv("SHELL")
		if strings.Contains(shell, "fish") {
			fmt.Printf("\nTo add gogen to your PATH, run:\n")
			fmt.Printf("fish_add_path %s\n", binDir)
		} else if strings.Contains(shell, "zsh") {
			fmt.Printf("\nTo add gogen to your PATH, add this to your ~/.zshrc:\n")
			fmt.Printf("export PATH=\"%s:$PATH\"\n", binDir)
		} else {
			fmt.Printf("\nTo add gogen to your PATH, add this to your ~/.bashrc or ~/.profile:\n")
			fmt.Printf("export PATH=\"%s:$PATH\"\n", binDir)
		}
		fmt.Println("Then restart your terminal or run 'source ~/.bashrc' (or equivalent)")
		fmt.Println("After that, you can run 'gogen --help' from anywhere")
	}
}
