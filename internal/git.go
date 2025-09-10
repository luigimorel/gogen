package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func (pg *ProjectGenerator) InitGitRepository(projectName, template string) error {
	if err := exec.Command("git", "init").Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
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

func (pg *ProjectGenerator) RemoveGitRepository(dirName string) error {
	if err := os.Chdir(dirName); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return nil
	}

	if err := os.RemoveAll(".git"); err != nil {
		return fmt.Errorf("failed to remove .git directory: %w", err)
	}

	return nil
}
