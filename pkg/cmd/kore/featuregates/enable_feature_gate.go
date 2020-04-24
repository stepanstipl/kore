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

	"github.com/appvia/kore/pkg/kore"

	cmdutil "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
)

type EnableFeatureGateOption struct {
	cmdutil.Factory
	Name string
}

// NewCmdEnabledFeatureGate enables the given feature gate
func NewCmdEnabledFeatureGate(factory cmdutil.Factory) *cobra.Command {
	o := &EnableFeatureGateOption{Factory: factory}

	return &cobra.Command{
		Use:     "enable",
		Short:   "Enables the given feature gate",
		Run:     cmdutil.DefaultRunFunc(o),
		Example: "kore alpha feature-gate enables <FEATURE>",
	}
}

// Validate is called to validate the options
func (o *EnableFeatureGateOption) Validate() error {
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
func (o *EnableFeatureGateOption) Run() error {
	name := strings.ToLower(o.Name)

	if o.Config().FeatureGates == nil {
		o.Config().FeatureGates = map[string]bool{}
	}
	o.Config().FeatureGates[name] = true

	if err := o.UpdateConfig(); err != nil {
		return err
	}

	o.Println("%q feature gate was enabled", name)
	return nil
}
