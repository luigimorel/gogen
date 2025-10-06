package wizard

type ProjectConfig struct {
	Name       string
	Module     string
	Template   string
	Router     string
	Frontend   string
	TypeScript bool
	Tailwind   bool
	Runtime    string
	Editor     string
	Docker     bool
	Dir        string
}