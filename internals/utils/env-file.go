package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateEnvFile(dirName string) error {
	var envContent string
	if dirName == "api" {
		envContent = `
PORT=8080
 `
	} else if dirName == "frontend" {
		envContent = `
VITE_API_URL=http://localhost:8080
VITE_API_BASE_PATH=/api

# Development
VITE_NODE_ENV=development
`
	}
	envPath := filepath.Join(dirName, ".env.example")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		fmt.Printf("Warning: failed to create .env.example: %v\n", err)
	}

	return nil
}
