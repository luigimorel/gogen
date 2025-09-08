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

type Installer struct {
	Method string
	Force  bool
}

func NewInstaller(method string, force bool) *Installer {
	return &Installer{
		Method: method,
		Force:  force,
	}
}

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

			installer := NewInstaller(method, force)
			return installer.execute()
		},
	}
}

func (i *Installer) execute() error {
	fmt.Println("Installing gogen CLI...")

	switch i.Method {
	case "auto":
		return i.autoInstall()
	case "binary":
		return i.binaryInstall()
	case "nix":
		return i.nixInstall()
	case "brew":
		return i.brewInstall()
	default:
		return fmt.Errorf("unsupported installation method: %s", i.Method)
	}
}

func (i *Installer) autoInstall() error {
	switch runtime.GOOS {
	case "darwin":
		if i.commandExists("brew") {
			fmt.Println("Detected Homebrew, using brew installation...")
			return i.brewInstall()
		}
		return i.binaryInstall()
	case "linux":
		if i.commandExists("nix-env") || i.commandExists("nix") {
			fmt.Println("Detected Nix, using nix installation...")
			return i.nixInstall()
		}
		return i.binaryInstall()
	case "windows":
		return i.binaryInstall()
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func (i *Installer) binaryInstall() error {
	binDir := i.getBinaryInstallDir()
	binPath := filepath.Join(binDir, i.getBinaryName())

	if !i.Force && i.fileExists(binPath) {
		fmt.Printf("gogen is already installed at %s\n", binPath)
		fmt.Println("Use --force to reinstall")
		return nil
	}

	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	downloadURL := i.getDownloadURL()
	fmt.Printf("Downloading from %s...\n", downloadURL)

	if err := i.downloadFile(downloadURL, binPath); err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(binPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}
	}

	fmt.Printf("gogen installed successfully to %s\n", binPath)
	i.printPathInstructions(binDir)

	return nil
}

func (i *Installer) nixInstall() error {
	if !i.commandExists("nix-env") && !i.commandExists("nix") {
		return fmt.Errorf("nix is not installed on this system")
	}

	fmt.Println("Nix package not available yet, using binary installation...")
	return i.binaryInstall()
}

func (i *Installer) brewInstall() error {
	if !i.commandExists("brew") {
		return fmt.Errorf("homebrew is not installed on this system")
	}

	fmt.Println("Homebrew formula not available yet, using binary installation...")
	return i.binaryInstall()
}

func (i *Installer) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (i *Installer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (i *Installer) getBinaryInstallDir() string {
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

func (i *Installer) getBinaryName() string {
	if runtime.GOOS == "windows" {
		return "gogen.exe"
	}
	return "gogen"
}

func (i *Installer) getDownloadURL() string {
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

func (i *Installer) downloadFile(url, filepath string) error {
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

func (i *Installer) printPathInstructions(binDir string) {
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
