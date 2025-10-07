//nolint:depguard
package wizard

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
)

const (
	templateAPI = "api"
	templateWeb = "web"
	templateCLI = "cli"
)

var projectNameRe = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

func RunInteractive() (*ProjectConfig, error) {
	return RunInteractiveWithPrompter(NewSurveyPrompter())
}

// RunInteractiveWithPrompter is the same interactive flow but accepts an injected
// Prompter implementation for easier testing.
func RunInteractiveWithPrompter(p Prompter) (*ProjectConfig, error) {
	cfg := &ProjectConfig{}

	if err := askProjectName(p, cfg); err != nil {
		return nil, err
	}

	if err := askTemplate(p, cfg); err != nil {
		return nil, err
	}

	if err := askRouter(p, cfg); err != nil {
		return nil, err
	}

	if err := askFrontendOptions(p, cfg); err != nil {
		return nil, err
	}

	if err := askEditor(p, cfg); err != nil {
		return nil, err
	}

	if err := askDocker(p, cfg); err != nil {
		return nil, err
	}

	if err := askModule(p, cfg); err != nil {
		return nil, err
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

func handleInterrupt(err error) error {
	if errors.Is(err, terminal.InterruptErr) {
		return fmt.Errorf("aborted by user")
	}
	return err
}

func askProjectName(p Prompter, cfg *ProjectConfig) error {
	name, err := p.Input("Project name:", withValidate(validateProjectName))
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Name = name
	cfg.Dir = name // Directory defaults to name
	return nil
}

func askTemplate(p Prompter, cfg *ProjectConfig) error {
	templateChoices := []Choice{
		{Value: templateAPI, Label: templateAPI, Description: "REST API server"},
		{Value: templateWeb, Label: templateWeb, Description: "Web server (optionally with frontend)"},
		{Value: templateCLI, Label: templateCLI, Description: "Command-line application"},
	}
	template, err := p.Select("Choose a template:", templateChoices)
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Template = template
	return nil
}

func askRouter(p Prompter, cfg *ProjectConfig) error {
	// Router for api or web
	if cfg.Template == templateAPI || cfg.Template == templateWeb {
		routerChoices := []Choice{
			{Value: "stdlib", Label: "stdlib", Description: "Go standard library"},
			{Value: "chi", Label: "chi", Description: "Lightweight router"},
			{Value: "gorilla", Label: "gorilla", Description: "Gorilla Mux"},
			{Value: "httprouter", Label: "httprouter", Description: "High performance"},
		}
		router, err := p.Select("Select a router:", routerChoices)
		if err != nil {
			return handleInterrupt(err)
		}
		cfg.Router = router
	}
	return nil
}

func askFrontendOptions(p Prompter, cfg *ProjectConfig) error {
	if cfg.Template != templateWeb {
		return nil
	}

	hasFrontend, err := p.Confirm("Add a frontend framework?", true)
	if err != nil {
		return handleInterrupt(err)
	}
	if !hasFrontend {
		return nil
	}

	feChoices := []Choice{
		{Value: "react", Label: "react", Description: "React + Vite"},
		{Value: "vue", Label: "vue", Description: "Vue 3 + Vite"},
		{Value: "svelte", Label: "svelte", Description: "Svelte + Vite"},
		{Value: "solidjs", Label: "solidjs", Description: "SolidJS + Vite"},
		{Value: "angular", Label: "angular", Description: "Angular CLI"},
	}
	frontend, err := p.Select("Choose a frontend:", feChoices)
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Frontend = frontend

	// Runtime
	runtime, err := p.Select("JS runtime:", []Choice{
		{Value: "node", Label: "node", Description: "Default Node.js"},
		{Value: "bun", Label: "bun", Description: "Bun runtime"},
	})
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Runtime = runtime

	if cfg.Frontend != "angular" { // Angular = TS by default
		ts, err := p.Confirm("Use TypeScript?", true)
		if err != nil {
			return handleInterrupt(err)
		}
		cfg.TypeScript = ts
	} else {
		cfg.TypeScript = true
	}

	tailwind, err := p.Confirm("Add Tailwind CSS?", false)
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Tailwind = tailwind

	return nil
}

func askEditor(p Prompter, cfg *ProjectConfig) error {
	if cfg.Template == templateCLI { // Editor integration might be useful more broadly, optional
		return nil
	}

	addEditor, err := p.Confirm("Add editor AI template (cursor/vscode/jetbrains)?", false)
	if err != nil {
		return handleInterrupt(err)
	}
	if !addEditor {
		return nil
	}

	editor, err := p.Input("Editor (cursor/vscode/jetbrains):", withValidate(validateEditor))
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Editor = editor
	return nil
}

func askDocker(p Prompter, cfg *ProjectConfig) error {
	docker, err := p.Confirm("Add Docker support?", true)
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Docker = docker
	return nil
}

func askModule(p Prompter, cfg *ProjectConfig) error {
	useModule, err := p.Confirm("Specify a custom Go module path?", false)
	if err != nil {
		return handleInterrupt(err)
	}
	if !useModule {
		return nil
	}

	module, err := p.Input("Module path (e.g. github.com/user/"+cfg.Name+"):", withValidate(validateModule))
	if err != nil {
		return handleInterrupt(err)
	}
	cfg.Module = module
	return nil
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
