package internal

import (
	"fmt"
	"os"
)

type DirectoryManager struct {
	originalDir string
}

func NewDirectoryManager() (*DirectoryManager, error) {
	originalDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}
	return &DirectoryManager{originalDir: originalDir}, nil
}

func (dm *DirectoryManager) ChangeToDir(dir string) error {
	return os.Chdir(dir)
}

func (dm *DirectoryManager) RootDir() error {
	return os.Chdir(dm.originalDir)
}
