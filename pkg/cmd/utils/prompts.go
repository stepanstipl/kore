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

package utils

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type Prompt struct {
	Id          string
	LabelSuffix string
	ErrMsg      string
	Value       *string
}

func (p *Prompt) Do() error {
	var value string
	if p.Value != nil {
		value = *p.Value
	}
	runner := promptui.Prompt{
		Label:     p.Id + " " + p.LabelSuffix,
		AllowEdit: true,
		Default:   value,
		Validate: func(in string) error {
			if len(in) == 0 {
				return fmt.Errorf(p.ErrMsg, p.Id)
			}
			return nil
		},
	}

	gathered, err := runner.Run()
	if err != nil {
		return err
	}

	*p.Value = strings.TrimSpace(gathered)
	return nil
}

type Prompts []*Prompt

func (p Prompts) Collect() error {
	for _, p := range p {
		if err := p.Do(); err != nil {
			return err
		}
	}
	return nil
}
