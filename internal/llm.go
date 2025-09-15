package internal

import (
	"fmt"
	"os"
)

type LLMTemplate struct {
}

func NewLLMTemplate() *LLMTemplate {
	return &LLMTemplate{}
}

func (lt *LLMTemplate) CreateTemplate(template, frontendFramework, runtime, router string) error {
	var content string
	var filePath string

	switch template {
	case "cursor":
		content = lt.generateCursorContent(frontendFramework, runtime, router)
		filePath = ".cursorrules"
	case "vscode":
		if err := os.MkdirAll(".vscode", 0755); err != nil {
			return fmt.Errorf("failed to create .vscode directory: %w", err)
		}
		content = lt.generateVSCodeContent(frontendFramework, runtime, router)
		filePath = ".vscode/settings.json"
	case "jetbrains":
		content = lt.generateJetbrainsContent(frontendFramework, runtime, router)
		filePath = ".aiassistant"
	default:
		return fmt.Errorf("unsupported template: %s", template)
	}

	return os.WriteFile(filePath, []byte(content), 0600)
}

func (lt *LLMTemplate) generateCursorContent(frontendFramework, runtime, router string) string {
	baseContent := `# Cursor AI Rules for Go Development

## Project Context
This is a Go project using:
- Go modules for dependency management
- Standard Go project structure
- Web development with ` + router + ` router 
- CLI applications using urfave/cli/v2`

	if frontendFramework != "" {
		baseContent += fmt.Sprintf(`
- Frontend: %s framework`, frontendFramework)

		switch frontendFramework {
		case react:
			baseContent += `
- React components with hooks and modern patterns
- JSX/TSX for component templates`
		case vue:
			baseContent += `
- Vue 3 composition API
- Single File Components (SFC)`
		case svelte:
			baseContent += `
- Svelte components with reactive statements
- SvelteKit for full-stack applications`
		case solidjs:
			baseContent += `
- SolidJS with fine-grained reactivity
- JSX templating with solid patterns`
		case angular:
			baseContent += `
- Angular with TypeScript
- Component-based architecture with dependency injection`
		}
	}

	if runtime != "" && runtime != node {
		baseContent += fmt.Sprintf(`
- JavaScript runtime: %s`, runtime)

		if runtime == bun {
			baseContent += `
- Fast package management and bundling with Bun
- TypeScript support out of the box`
		}
	}

	content := baseContent + `

## Development Guidelines

### Code Style
- Follow Go conventions and best practices
- Use gofmt for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Keep functions small and focused

### Project Structure
- Follow standard Go project layout
- Use cmd/ for main applications
- Use internal/ for private application code
- Use pkg/ for library code that can be imported by external applications

### Error Handling
- Always handle errors appropriately
- Use fmt.Errorf for error wrapping
- Return errors rather than panicking in most cases
- Log errors when appropriate

### Testing
- Write unit tests for all public functions
- Use table-driven tests when appropriate
- Follow Go testing conventions
- Aim for good test coverage

### Dependencies
- Minimize external dependencies
- Use standard library when possible
- Keep go.mod clean and up to date
- Use go mod tidy regularly

## Specific Project Rules
- When working with web handlers, ensure proper HTTP status codes
- Use proper middleware patterns for common functionality
- Follow RESTful API conventions when applicable
- Handle graceful shutdowns for server applications
- Use environment variables for configuration

## Code Generation
- When generating boilerplate code, follow the existing patterns in the project
- Ensure generated code is idiomatic Go
- Add appropriate comments and documentation
- Consider edge cases and error conditions`

	if frontendFramework != "" {
		content += `

## Frontend Development Guidelines

### General Frontend Rules
- Maintain separation between API and frontend concerns
- Use proper error handling for API calls
- Implement loading states and error boundaries`

		switch frontendFramework {
		case react:
			content += `

### React-Specific Rules
- Use functional components with hooks
- Follow React best practices for state management
- Use proper key props for list items
- Implement proper cleanup in useEffect
- Use TypeScript for better type safety (if enabled)`
		case vue:
			content += `

### Vue-Specific Rules
- Use Composition API for new components
- Follow Vue 3 best practices
- Use proper reactive references and computed properties
- Implement proper component lifecycle management`
		case svelte:
			content += `

### Svelte-Specific Rules
- Use reactive statements ($:) appropriately
- Follow Svelte best practices for component communication
- Use stores for global state management
- Implement proper component lifecycle`
		case solidjs:
			content += `

### SolidJS-Specific Rules
- Use signals and effects properly
- Follow SolidJS patterns for reactivity
- Implement proper resource management
- Use JSX patterns specific to SolidJS`
		case angular:
			content += `

### Angular-Specific Rules
- Use Angular CLI for code generation
- Follow Angular style guide conventions
- Use proper dependency injection patterns
- Implement proper component lifecycle hooks`
		}

		if runtime == bun {
			content += `

### Bun Runtime Guidelines
- Leverage Bun's fast package installation
- Use Bun's built-in bundler when appropriate
- Take advantage of Bun's TypeScript support`
		}
	}

	return content
}

func (lt *LLMTemplate) generateVSCodeContent(frontendFramework, runtime, router string) string {
	baseSettings := `{
    "github.copilot.enable": {
        "*": true,
        "yaml": true,
        "plaintext": true,
        "markdown": true,
        "go": true`

	if frontendFramework != "" {
		baseSettings += `,
        "javascript": true,
        "typescript": true,
        "json": true,
        "html": true,
        "css": true`

		switch frontendFramework {
		case react:
			baseSettings += `,
        "javascriptreact": true,
        "typescriptreact": true`
		case vue:
			baseSettings += `,
        "vue": true`
		case svelte:
			baseSettings += `,
        "svelte": true`
		case angular:
			baseSettings += `,
        "html": true`
		}
	}

	baseSettings += `
    },
    "github.copilot.chat.localeOverride": "en",
    "github.copilot.advanced": {
        "debug.overrideEngine": "gpt-4",
        "length": 3000
    },
    "go.toolsManagement.autoUpdate": true,
    "go.useLanguageServer": true,
    "go.formatTool": "gofmt",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "go.buildTags": "integration",
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
        "source.organizeImports": true
    },
    "files.associations": {
        "*.go": "go",
        "go.mod": "go.mod",
        "go.sum": "go.sum"
    },
    "gopls": {
        "ui.completion.usePlaceholders": true,
        "ui.diagnostic.analyses": {
            "fieldalignment": false,
            "shadow": true
        }
    }`

	if frontendFramework != "" {
		switch frontendFramework {
		case react:
			baseSettings += `,
    "typescript.preferences.includePackageJsonAutoImports": "on",
    "typescript.suggest.autoImports": true,
    "javascript.suggest.autoImports": true,
    "emmet.includeLanguages": {
        "javascript": "javascriptreact",
        "typescript": "typescriptreact"
    },
    "emmet.triggerExpansionOnTab": true`
		case vue:
			baseSettings += `,
    "vetur.validation.template": false,
    "vetur.validation.script": false,
    "vetur.validation.style": false,
    "volar.takeOverMode": true`
		case svelte:
			baseSettings += `,
    "svelte.enable-ts-plugin": true,
    "typescript.preferences.includePackageJsonAutoImports": "on"`
		case angular:
			baseSettings += `,
    "typescript.preferences.includePackageJsonAutoImports": "on",
    "angular.enableCodeCompletion": true`
		}

		if runtime == bun {
			baseSettings += `,
    "terminal.integrated.defaultProfile.linux": "bash",
    "npm.packageManager": "bun"`
		}
	}

	baseSettings += `
}`

	return baseSettings
}

func (lt *LLMTemplate) generateJetbrainsContent(frontendFramework, runtime, router string) string {
	baseContent := `# JetBrains AI Assistant Rules for Go Development

## Project Context
This is a Go project with the following characteristics:
- Modular Go application with multiple components
- Web development with ` + router + ` router 
- CLI application support
- Standard Go project structure`

	if frontendFramework != "" {
		baseContent += fmt.Sprintf(`
- Frontend framework: %s`, frontendFramework)

		switch frontendFramework {
		case react:
			baseContent += ` with modern hooks and functional components`
		case vue:
			baseContent += ` with Composition API and SFC`
		case svelte:
			baseContent += ` with reactive statements and SvelteKit`
		case solidjs:
			baseContent += ` with fine-grained reactivity`
		case angular:
			baseContent += ` with TypeScript and dependency injection`
		}

		if runtime == bun {
			baseContent += fmt.Sprintf(`
- JavaScript runtime: %s for fast package management and execution`, runtime)
		}
	}

	content := baseContent + `

## Development Guidelines

### Code Quality
- Maintain high code quality with proper error handling
- Use Go idioms and conventions consistently
- Write self-documenting code with clear variable names
- Follow the principle of least astonishment

### Architecture
- Keep the architecture simple and maintainable
- Use dependency injection where appropriate
- Separate concerns properly
- Follow SOLID principles where applicable to Go

### Performance
- Profile code when performance is critical
- Use appropriate data structures
- Avoid premature optimization
- Be mindful of memory allocations

### Security
- Validate all inputs
- Use proper authentication and authorization
- Follow security best practices for web applications
- Handle sensitive data appropriately

### Documentation
- Write clear and concise comments
- Document public APIs thoroughly
- Keep README.md up to date
- Use godoc conventions for documentation

## AI Assistant Preferences
- Suggest Go-idiomatic solutions
- Prefer standard library over third-party packages when possible
- Focus on readability and maintainability
- Consider error handling in all suggestions
- Recommend testing strategies for new code
- Follow the existing project patterns and conventions

## Code Review Focus
- Error handling completeness
- Resource cleanup (defer statements)
- Concurrency safety where applicable
- Interface usage and design
- Performance implications of suggested changes`

	if frontendFramework != "" {
		content += `

## Frontend Development Guidelines

### API Integration
- Ensure proper separation between backend and frontend
- Use appropriate HTTP status codes and error handling
- Implement proper CORS configuration when needed
- Follow RESTful API conventions

### Frontend Best Practices`

		switch frontendFramework {
		case react:
			content += `
- Use functional components with hooks
- Implement proper error boundaries
- Use React.memo for performance optimization
- Follow React testing library best practices
- Use proper TypeScript types when applicable`
		case vue:
			content += `
- Use Composition API for new components
- Implement proper reactive state management
- Use Vue 3 best practices for component communication
- Follow Vue testing utils conventions
- Ensure proper component lifecycle management`
		case svelte:
			content += `
- Use reactive statements ($:) appropriately
- Implement proper store patterns for state management
- Follow SvelteKit conventions for routing and data loading
- Use proper component communication patterns`
		case solidjs:
			content += `
- Use signals and effects properly
- Implement proper resource management
- Follow SolidJS patterns for reactivity
- Use proper JSX patterns specific to SolidJS`
		case angular:
			content += `
- Follow Angular style guide conventions
- Use proper dependency injection patterns
- Implement proper component lifecycle hooks
- Use Angular CLI for consistent code generation
- Follow RxJS best practices for reactive programming`
		}

		if runtime == bun {
			content += `

### Bun Runtime Optimization
- Leverage Bun's fast package installation
- Use Bun's built-in bundler for optimal performance
- Take advantage of Bun's native TypeScript support
- Consider Bun-specific APIs for enhanced performance`
		}
	}

	return content
}
