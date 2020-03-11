/**
 * Copyright 2020 Appvia Ltd <info@appvia.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package korectl

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type prompt struct {
	id          string
	labelSuffix string
	errMsg      string
	value       *string
}

func (p *prompt) do() error {
	var value string
	if p.value != nil {
		value = *p.value
	}
	runner := promptui.Prompt{
		Label:     p.id + " " + p.labelSuffix,
		AllowEdit: true,
		Default:   value,
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

	*p.value = strings.TrimSpace(gathered)
	return nil
}

type prompts []*prompt

func (p prompts) collect() error {
	for _, p := range p {
		if err := p.do(); err != nil {
			return err
		}
	}
	return nil
}
