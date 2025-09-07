package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func InitGitRepository(projectName, template string) error {
	if err := exec.Command("git", "init").Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	gitignoreContent := createGitignoreContent(template)
	if err := os.WriteFile(".gitignore", []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	if err := exec.Command("git", "add", ".").Run(); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	commitMessage := fmt.Sprintf("Initial commit: %s %s project", projectName, template)
	if err := exec.Command("git", "commit", "-m", commitMessage).Run(); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	return nil
}

func createGitignoreContent(template string) string {
	baseIgnore := `
# Environment variables
.env
.env.local
.env.*.local
`

	switch template {
	case "web":
		return baseIgnore + `
# Frontend dependencies and build files

frontend/.env
frontend/.env.local
frontend/.env.*.local
`
	default:
		return baseIgnore
	}
}
