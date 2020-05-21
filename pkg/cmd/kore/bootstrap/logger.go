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

package bootstrap

import (
	"github.com/appvia/kore/pkg/cmd/kore/bootstrap/providers"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
)

type providerLogger struct {
	cmdutil.Factory
}

// newProviderLogger provides a logger
func newProviderLogger(factory cmdutil.Factory) providers.Logger {
	return &providerLogger{Factory: factory}
}

// Info prints a logging line from the provider
func (p *providerLogger) Info(message string, args ...interface{}) {
	p.Println("   ◉ "+message, args...)
}
