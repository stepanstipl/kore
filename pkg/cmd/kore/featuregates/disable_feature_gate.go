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

package featuregates

import (
	"errors"
	"fmt"
	"strings"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"

	"github.com/spf13/cobra"
)

type DisableFeatureGateOption struct {
	cmdutil.Factory
	Name string
}

// NewCmdDisableFeatureGate disables the given feature gate
func NewCmdDisableFeatureGate(factory cmdutil.Factory) *cobra.Command {
	o := &DisableFeatureGateOption{Factory: factory}

	return &cobra.Command{
		Use:     "disable",
		Short:   "Disables the given feature gate",
		Run:     cmdutil.DefaultRunFunc(o),
		Example: "kore alpha feature-gates disable <FEATURE>",
	}
}

// Validate is called to validate the options
func (o *DisableFeatureGateOption) Validate() error {
	name := strings.ToLower(o.Name)

	if name == "" {
		return errors.New("feature gate name is required")
	}
	defaultFeatureGates := kore.DefaultFeatureGates()
	if _, exists := defaultFeatureGates[name]; !exists {
		return fmt.Errorf("feature gate %q is invalid", name)
	}
	return nil
}

// Run implements the action
func (o *DisableFeatureGateOption) Run() error {
	name := strings.ToLower(o.Name)

	if o.Config().FeatureGates == nil {
		o.Config().FeatureGates = map[string]bool{}
	}
	o.Config().FeatureGates[name] = false

	if err := o.UpdateConfig(); err != nil {
		return err
	}

	o.Println("%q feature gate was disabled", name)
	return nil
}
