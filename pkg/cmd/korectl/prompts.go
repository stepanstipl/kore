package korectl

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

type prompt struct {
	id          string
	labelSuffix string
	errMsg      string
	value       string
}

func (p *prompt) do() error {
	runner := promptui.Prompt{
		Label: p.id + " " + p.labelSuffix,
		Validate: func(in string) error {
			if len(in) == 0 {
				return fmt.Errorf(p.errMsg, p.id)
			}
			return nil
		},
	}

	gathered, err := runner.Run()
	if err != nil {
		return err
	}

	p.value = gathered
	return nil
}

type prompts struct {
	prompts []*prompt
}

func (p *prompts) collect() error {
	for _, p := range p.prompts {
		if err := p.do(); err != nil {
			return err
		}
	}
	return nil
}

func (p *prompts) getValue(id string) string {
	for _, p := range p.prompts {
		if p.id == id {
			return p.value
		}
	}
	return ""
}
