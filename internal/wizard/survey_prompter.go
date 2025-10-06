package wizard

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

type SurveyPrompter struct{}

func NewSurveyPrompter() *SurveyPrompter { return &SurveyPrompter{} }

func (sp *SurveyPrompter) Input(label string, opts ...InputOption) (string, error) {
	cfg := &inputConfig{}
	for _, o := range opts {
		o(cfg)
	}
	prompt := &survey.Input{
		Message: label,
		Default: cfg.Default,
	}
	var answer string
	err := survey.AskOne(prompt, &answer, survey.WithValidator(func(ans interface{}) error {
		s := ans.(string)
		if s == "" && !cfg.AllowEmpty {
			return fmt.Errorf("value required")
		}
		if cfg.Validate != nil {
			return cfg.Validate(s)
		}
		return nil
	}))
	return answer, err
}

func (sp *SurveyPrompter) Select(label string, choices []Choice) (string, error) {
	options := make([]string, len(choices))
	valMap := make(map[string]string)
	for i, c := range choices {
		display := fmt.Sprintf("%s  –  %s", c.Label, c.Description)
		options[i] = display
		valMap[display] = c.Value
	}
	var selected string
	err := survey.AskOne(&survey.Select{
		Message: label,
		Options: options,
	}, &selected)
	if err != nil {
		return "", err
	}
	return valMap[selected], nil
}

func (sp *SurveyPrompter) Confirm(label string, defaultYes bool) (bool, error) {
	var ok bool
	err := survey.AskOne(&survey.Confirm{
		Message: label,
		Default: defaultYes,
	}, &ok)
	return ok, err
}