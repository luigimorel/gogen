# Changelog — feat/interactive-menu

This file summarizes the changes implemented for the interactive project creation wizard (branch: feat/interactive-menu).

## Summary
- Add an interactive "wizard" for `gogen new` to guide project creation (templates, router, frontend, runtime, TS, Tailwind, editor, Docker, module path).
- Make the interactive flow testable and robust to Ctrl+C (returns a friendly "aborted by user" error).
- Add unit tests for the wizard (web path, cli path, abort path).
- Wire prompts to a prompter abstraction so interactive behavior can be mocked.

## Added
- internal/wizard/prompter.go
  - Prompter interface and choice/input helper types.
- internal/wizard/survey_prompter.go
  - Survey-backed Prompter implementation (AlecAivazis/survey/v2).
- internal/wizard/config.go
  - ProjectConfig struct.
- internal/wizard/run.go (updated)
  - `RunInteractive()` and `RunInteractiveWithPrompter(p Prompter)` flow.
  - printSummary() and generateCommand() helper functions.
  - Terminal interrupt handling (maps survey terminal.InterruptErr → "aborted by user").
- internal/wizard/run_test.go
  - Tests with fakePrompter and abortPrompter to exercise main flows.

## Changed
- cmd/new.go
  - Added `--interactive` flag to the `gogen new` command and integrated the wizard (Action now calls `wizard.RunInteractive()` when interactive).
  - Adjusted flag validation so interactive mode is not blocked by "required" CLI flag validation (validation happens at runtime for non-interactive mode).
- go.mod / go.sum
  - Added dependency on github.com/AlecAivazis/survey/v2.

## Tests
- Unit tests added for:
  - Web + React + TypeScript + Tailwind + bun + docker generation path.
  - CLI minimal generation path.
  - Abort (Ctrl+C) behavior.
- Run tests:
  - go test ./internal/wizard -v

## Behavior notes
- When the user presses Ctrl+C during the prompts, the wizard exits with a consistent "aborted by user" error.
- The interactive flow prints a configuration summary and shows the exact generated CLI command before proceeding.
- The code supports injecting a fake Prompter for deterministic, non-interactive tests.

## How to verify locally
1. Tidy and build:
   - go mod tidy
   - go build -o ../gogen.exe .
2. Run unit tests:
   - go test ./internal/wizard -v
3. Try the interactive flow (use a temp dir):
   - ..\gogen.exe new --interactive --name temp
   - Fill the prompts and confirm generation.

## Suggested follow-ups (not included in this PR)
- Add integration tests that run the built binary and feed scripted answers (expect/pexpect).
- Persist defaults to `~/.config/gogen` (optional).
- Improve the SurveyPrompter to avoid display collisions (use indices) for robust mapping.

## PR / Commit suggestions
- Commit 1: feat(wizard): add prompter abstraction and survey implementation
- Commit 2: feat(wizard): interactive flow (RunInteractiveWithPrompter)
- Commit 3: test(wizard): add fake prompter tests + abort test
- Commit 4: refactor(new): wire --interactive into cmd/new.go
- Commit 5: docs: update README with interactive usage

---
Branch: feat/interactive-menu
