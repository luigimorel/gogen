package fsutils

import (
	"fmt"
	"os"
	"path/filepath"
)

type CreateFilesJob struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type CreateJob []*CreateFilesJob

func (cj CreateJob) Execute() error {
	for _, job := range cj {
		if err := os.MkdirAll(filepath.Dir(job.Filename), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		err := os.WriteFile(job.Filename, []byte(job.Content), 0644)
		switch {
		case os.IsExist(err):
			fmt.Printf("Warning: file already exists at %s\n", job.Filename)
		case err != nil:
			return fmt.Errorf("failed to create file: %w", err)
		}
	}
	return nil
}
