package wizard

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
)

var projectNameRe = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

func RunInteractive() (*ProjectConfig, error) {
	return RunInteractiveWithPrompter(NewSurveyPrompter())
}

// RunInteractiveWithPrompter is the same interactive flow but accepts an injected
// Prompter implementation for easier testing.
func RunInteractiveWithPrompter(p Prompter) (*ProjectConfig, error) {
	cfg := &ProjectConfig{}

	var err error
	if cfg.Name, err = p.Input("Project name:", withValidate(validateProjectName)); err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			return nil, fmt.Errorf("aborted by user")
		}
		return nil, err
	}
	// Directory defaults to name
	cfg.Dir = cfg.Name

	templateChoices := []Choice{
		{Value: "api", Label: "api", Description: "REST API server"},
		{Value: "web", Label: "web", Description: "Web server (optionally with frontend)"},
		{Value: "cli", Label: "cli", Description: "Command-line application"},
	}
	if cfg.Template, err = p.Select("Choose a template:", templateChoices); err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			return nil, fmt.Errorf("aborted by user")
		}
		return nil, err
	}

	// Router for api or web
	if cfg.Template == "api" || cfg.Template == "web" {
		routerChoices := []Choice{
			{Value: "stdlib", Label: "stdlib", Description: "Go standard library"},
			{Value: "chi", Label: "chi", Description: "Lightweight router"},
			{Value: "gorilla", Label: "gorilla", Description: "Gorilla Mux"},
			{Value: "httprouter", Label: "httprouter", Description: "High performance"},
		}
		if cfg.Router, err = p.Select("Select a router:", routerChoices); err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				return nil, fmt.Errorf("aborted by user")
			}
			return nil, err
		}
	}

	if cfg.Template == "web" {
		hasFrontend, err := p.Confirm("Add a frontend framework?", true)
		if err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				return nil, fmt.Errorf("aborted by user")
			}
			return nil, err
		}
		if hasFrontend {
			feChoices := []Choice{
				{Value: "react", Label: "react", Description: "React + Vite"},
				{Value: "vue", Label: "vue", Description: "Vue 3 + Vite"},
				{Value: "svelte", Label: "svelte", Description: "Svelte + Vite"},
				{Value: "solidjs", Label: "solidjs", Description: "SolidJS + Vite"},
				{Value: "angular", Label: "angular", Description: "Angular CLI"},
			}
			if cfg.Frontend, err = p.Select("Choose a frontend:", feChoices); err != nil {
				if errors.Is(err, terminal.InterruptErr) {
					return nil, fmt.Errorf("aborted by user")
				}
				return nil, err
			}
			// Runtime
			if cfg.Runtime, err = p.Select("JS runtime:", []Choice{
				{Value: "node", Label: "node", Description: "Default Node.js"},
				{Value: "bun", Label: "bun", Description: "Bun runtime"},
			}); err != nil {
				if errors.Is(err, terminal.InterruptErr) {
					return nil, fmt.Errorf("aborted by user")
				}
				return nil, err
			}
			if cfg.Frontend != "angular" { // Angular = TS by default
				if cfg.TypeScript, err = p.Confirm("Use TypeScript?", true); err != nil {
					if errors.Is(err, terminal.InterruptErr) {
						return nil, fmt.Errorf("aborted by user")
					}
					return nil, err
				}
			} else {
				cfg.TypeScript = true
			}
			if cfg.Tailwind, err = p.Confirm("Add Tailwind CSS?", false); err != nil {
				if errors.Is(err, terminal.InterruptErr) {
					return nil, fmt.Errorf("aborted by user")
				}
				return nil, err
			}
		}
	}

	if cfg.Template != "cli" { // Editor integration might be useful more broadly, optional
		addEditor, err := p.Confirm("Add editor AI template (cursor/vscode/jetbrains)?", false)
		if err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				return nil, fmt.Errorf("aborted by user")
			}
			return nil, err
		}
		if addEditor {
			if cfg.Editor, err = p.Input("Editor (cursor/vscode/jetbrains):", withValidate(validateEditor)); err != nil {
				if errors.Is(err, terminal.InterruptErr) {
					return nil, fmt.Errorf("aborted by user")
				}
				return nil, err
			}
		}
	}

	if cfg.Docker, err = p.Confirm("Add Docker support?", true); err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			return nil, fmt.Errorf("aborted by user")
		}
		return nil, err
	}

	// Optional module path
	useModule, err := p.Confirm("Specify a custom Go module path?", false)
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			return nil, fmt.Errorf("aborted by user")
		}
		return nil, err
	}
	if useModule {
		if cfg.Module, err = p.Input("Module path (e.g. github.com/user/"+cfg.Name+"):", withValidate(validateModule)); err != nil {
			if errors.Is(err, terminal.InterruptErr) {
				return nil, fmt.Errorf("aborted by user")
			}
			return nil, err
		}
	}

	printSummary(cfg)
	ok, err := p.Confirm("Proceed with generation?", true)
	if err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			return nil, fmt.Errorf("aborted by user")
		}
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("aborted by user")
	}

	return cfg, nil
}

func printSummary(c *ProjectConfig) {
	var b strings.Builder
	b.WriteString("\nConfiguration Summary:\n")
	b.WriteString("  Name: " + c.Name + "\n")
	b.WriteString("  Template: " + c.Template + "\n")
	if c.Router != "" {
		b.WriteString("  Router: " + c.Router + "\n")
	}
	if c.Frontend != "" {
		b.WriteString("  Frontend: " + c.Frontend + "\n")
		b.WriteString(fmt.Sprintf("  TypeScript: %v\n", c.TypeScript))
		b.WriteString(fmt.Sprintf("  Tailwind: %v\n", c.Tailwind))
		b.WriteString("  Runtime: " + c.Runtime + "\n")
	}
	if c.Editor != "" {
		b.WriteString("  Editor: " + c.Editor + "\n")
	}
	b.WriteString(fmt.Sprintf("  Docker: %v\n", c.Docker))
	if c.Module != "" {
		b.WriteString("  Module: " + c.Module + "\n")
	}
	b.WriteString("\nGenerated Command:\n")
	b.WriteString("  " + generateCommand(c) + "\n")
	fmt.Println(b.String())
}

func generateCommand(c *ProjectConfig) string {
	args := []string{"gogen", "new", "--name", c.Name, "--template", c.Template}
	if c.Router != "" {
		args = append(args, "--router", c.Router)
	}
	if c.Frontend != "" {
		args = append(args, "--frontend", c.Frontend, "--runtime", c.Runtime)
		if c.TypeScript && c.Frontend != "angular" {
			args = append(args, "--ts")
		}
		if c.Tailwind {
			args = append(args, "--tailwind")
		}
	}
	if c.Editor != "" {
		args = append(args, "--editor", c.Editor)
	}
	if c.Docker {
		args = append(args, "--docker")
	}
	if c.Module != "" {
		args = append(args, "--module", c.Module)
	}
	return strings.Join(args, " ")
}

func withValidate(f func(string) error) InputOption {
	return func(ic *inputConfig) {
		ic.Validate = f
	}
}

func validateProjectName(name string) error {
	if !projectNameRe.MatchString(name) {
		return fmt.Errorf("invalid project name (only alphanumeric, dash, underscore)")
	}
	return nil
}

func validateEditor(e string) error {
	switch e {
	case "cursor", "vscode", "jetbrains":
		return nil
	default:
		return fmt.Errorf("unsupported editor")
	}
}

func validateModule(m string) error {
	if strings.Count(m, "/") < 2 {
		return fmt.Errorf("module path should look like github.com/username/repo")
	}
	return nil
}
