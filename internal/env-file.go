package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func (pg *ProjectGenerator) CreateEnvFile(dirName string) error {
	var envContent string
	if dirName == "api" {
		envContent = `PORT=8080`
	} else if dirName == "frontend" {
		envContent = `
VITE_API_URL=http://localhost:8080
VITE_API_BASE_PATH=/api

# Development
VITE_NODE_ENV=development
`
	} else if dirName == "." {
		// Root level .env.example for API template
		envContent = `PORT=8080`
	}

	envPath := filepath.Join(dirName, ".env.example")
	if dirName == "." {
		envPath = ".env.example"
	}

	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		fmt.Printf("Warning: failed to create .env.example: %v\n", err)
	}

	return nil
}

func (pg *ProjectGenerator) CreateGitignoreFile(dirName, template string) error {
	var gitignoreContent string

	if dirName == "api" {
		gitignoreContent = `
.env
.env.local
.env.production.local
.env.*.local
`
	} else if dirName == "frontend" {
		gitignoreContent = `# Logs
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
pnpm-debug.log*

node_modules
dist
dist-ssr
*.local

# Environment variables
.env
.env.local
.env.production.local
.env.*.local
`
	}

	gitignorePath := filepath.Join(dirName, ".gitignore")
	if dirName == "." {
		gitignorePath = ".gitignore"
	}

	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore in %s: %w", dirName, err)
	}

	return nil
}
