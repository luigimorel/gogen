//nolint:depguard
package wizard

import (
	"errors"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
)

// fakePrompter is a simple label->answer map implementation for Prompter used in tests.
type fakePrompter struct {
	inputs   map[string]string
	selects  map[string]string
	confirms map[string]bool
}

func (f *fakePrompter) Input(label string, opts ...InputOption) (string, error) {
	if v, ok := f.inputs[label]; ok {
		return v, nil
	}
	return "", nil
}

func (f *fakePrompter) Select(label string, choices []Choice) (string, error) {
	if v, ok := f.selects[label]; ok {
		return v, nil
	}
	return "", nil
}

func (f *fakePrompter) Confirm(label string, defaultYes bool) (bool, error) {
	if v, ok := f.confirms[label]; ok {
		return v, nil
	}
	return defaultYes, nil
}

// abortPrompter simulates a Ctrl+C / interrupt from the terminal.
type abortPrompter struct{}

func (a *abortPrompter) Input(label string, opts ...InputOption) (string, error) {
	return "", terminal.InterruptErr
}

func (a *abortPrompter) Select(label string, choices []Choice) (string, error) {
	return "", terminal.InterruptErr
}

func (a *abortPrompter) Confirm(label string, defaultYes bool) (bool, error) {
	return false, terminal.InterruptErr
}

func TestRunInteractive_WebReact_TS_Tailwind(t *testing.T) {
	fp := &fakePrompter{
		inputs: map[string]string{
			"Project name:": "my-app",
		},
		selects: map[string]string{
			"Choose a template:": "web",
			"Select a router:":   "chi",
			"Choose a frontend:": "react",
			"JS runtime:":        "bun",
		},
		confirms: map[string]bool{
			"Add a frontend framework?":                         true,
			"Use TypeScript?":                                   true,
			"Add Tailwind CSS?":                                 true,
			"Add editor AI template (cursor/vscode/jetbrains)?": false,
			"Add Docker support?":                               true,
			"Specify a custom Go module path?":                  false,
			"Proceed with generation?":                          true,
		},
	}

	cfg, err := RunInteractiveWithPrompter(fp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "my-app" {
		t.Fatalf("unexpected name: %s", cfg.Name)
	}
	if cfg.Template != "web" || cfg.Router != "chi" || cfg.Frontend != "react" {
		t.Fatalf("unexpected selections: %+v", cfg)
	}
	if !cfg.TypeScript || !cfg.Tailwind || cfg.Runtime != "bun" {
		t.Fatalf("frontend flags not set correctly: %+v", cfg)
	}

	got := generateCommand(cfg)
	want := "gogen new --name my-app --template web --router chi --frontend react --runtime bun --ts --tailwind --docker"
	if got != want {
		t.Fatalf("generateCommand mismatch:\n got: %s\n want: %s", got, want)
	}
}

func TestRunInteractive_CLI_Minimal(t *testing.T) {
	fp := &fakePrompter{
		inputs: map[string]string{
			"Project name:": "cliapp",
		},
		selects: map[string]string{
			"Choose a template:": "cli",
		},
		confirms: map[string]bool{
			"Add Docker support?":              false,
			"Specify a custom Go module path?": false,
			"Proceed with generation?":         true,
		},
	}

	cfg, err := RunInteractiveWithPrompter(fp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Template != "cli" {
		t.Fatalf("expected cli template, got %s", cfg.Template)
	}
	if cfg.Router != "" || cfg.Frontend != "" {
		t.Fatalf("router/frontend should be empty for cli: %+v", cfg)
	}
}

func TestRunInteractive_Abort(t *testing.T) {
	// simulate Ctrl+C by returning terminal.InterruptErr from the prompter
	ap := &abortPrompter{}
	_, err := RunInteractiveWithPrompter(ap)
	if err == nil {
		t.Fatalf("expected abort error, got nil")
	}
	if !errors.Is(err, terminal.InterruptErr) && err.Error() != "aborted by user" {
		// we expect the wrapper to return a friendly "aborted by user" message.
		t.Fatalf("unexpected error for abort: %v", err)
	}
}
