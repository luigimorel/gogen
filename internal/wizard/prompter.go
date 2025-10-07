package wizard

type Prompter interface {
	Input(label string, opts ...InputOption) (string, error)
	Select(label string, choices []Choice) (string, error)
	Confirm(label string, defaultYes bool) (bool, error)
}

type Choice struct {
	Value       string
	Label       string
	Description string
}

type InputOption func(*inputConfig)

type inputConfig struct {
	Default     string
	Validate    func(string) error
	AllowEmpty  bool
	Placeholder string
}
