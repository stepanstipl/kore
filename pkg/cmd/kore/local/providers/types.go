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

package providers

import (
	"context"
	"io"
)

// Logger provides a logging interface to the providers
type Logger interface {
	// Infof provides local logging
	Infof(string, ...interface{})
	// Info provides local logging
	Info(string, ...interface{})
	// Stdout returns the writer
	Stdout() io.Writer
}

// CreateOptions are configurable for the provider creation
type CreateOptions struct {
	// DisableUI indicates the UI is disabled
	DisableUI bool
}

// Interface is the contract for a provider
type Interface interface {
	// Create is responsible for starting the cluster
	Create(ctx context.Context, name string, options CreateOptions) error
	// Destroy is responsible for destroy the cluster
	Destroy(ctx context.Context, name string) error
	// Export is responsible for exporting the kubeconfig
	Export(ctx context.Context, name string) (string, error)
	// Preflight checks we have the requirements
	Preflight(ctx context.Context) error
	// Stop is responsible for halting the provider
	Stop(ctx context.Context, name string) error
}
