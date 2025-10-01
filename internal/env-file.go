package internal

import (
	"fmt"
	"os"
	"path/filepath"

	constants "github.com/luigimorel/gogen/consants"
)

func (pg *ProjectGenerator) CreateEnvFile(dirType, dirName string) error {
	var envContent string
	if dirType == constants.APIDir {
		envContent = `PORT=8080`
	} else {
		envContent = `VITE_API_URL=http://localhost:8080
VITE_API_BASE_PATH=/api

# Development
VITE_NODE_ENV=development
`
	}

	envExamplePath := filepath.Join(dirName, ".env.example")
	envPath := filepath.Join(dirName, ".env")

	if dirName == "." {
		envExamplePath = ".env.example"
		envPath = ".env"
	}

	if err := os.MkdirAll(filepath.Dir(envExamplePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for env files: %w", err)
	}

	if err := os.WriteFile(envExamplePath, []byte(envContent), 0600); err != nil {
		return fmt.Errorf("failed to create .env.example: %w", err)
	}

	if err := os.WriteFile(envPath, []byte(envContent), 0600); err != nil {
		return fmt.Errorf("failed to create .env: %w", err)
	}

	return nil
}

func (pg *ProjectGenerator) CreateEnvConfig(dirName, framework string, useTypeScript bool) error {
	if framework == angular {
		return nil
	}

	fileExt := "js"
	if useTypeScript {
		fileExt = "ts"
	}

	configContent := `/// <reference types="vite/client" />
export const config = {
  apiUrl: import.meta.env.VITE_API_URL,
  apiBasePath: import.meta.env.VITE_API_BASE_PATH,
  nodeEnv: import.meta.env.VITE_NODE_ENV,
};

export default config;
`

	if err := os.WriteFile(filepath.Join(dirName, "src", "config."+fileExt), []byte(configContent), 0600); err != nil {
		return fmt.Errorf("failed to create env config file: %w", err)
	}

	return nil
}

func (pg *ProjectGenerator) CreateGitignoreFile(dirType, dirName string) error {
	var gitignoreContent string

	switch dirType {
	case "api":
		gitignoreContent = `.env
.env.local
.env.production.local
.env.*.local
tmp
`
	case "cli":
		gitignoreContent = `# Binaries for programs and plugins
*.exe
tmp
main`
	default:
		gitignoreContent = `logs
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

	if err := os.MkdirAll(filepath.Dir(gitignorePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for .gitignore: %w", err)
	}

	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0600); err != nil {
		return fmt.Errorf("failed to create .gitignore in %s: %w", dirName, err)
	}

	return nil
}
