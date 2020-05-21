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

package kind

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/appvia/kore/pkg/cmd/kore/bootstrap/providers"
)

var (
	// KindConfiguration is the configuration for kind
	KindConfiguration = `
apiVersion: kind.x-k8s.io/v1alpha4
kind: Cluster
nodes:
- role: control-plane
  image: kindest/node:v1.15.11@sha256:6cc31f3533deb138792db2c7d1ffc36f7456a06f1db5556ad3b6927641016f50
  extraPortMappings:
  - containerPort: 3000
    hostPort: 3000
    protocol: TCP
  - containerPort: 10080
    hostPort: 10080
    protocol: TCP
`
)

type providerImpl struct {
	providers.Logger
	// path is the file path to the kind binary
	path string
}

// New creates and returns a kind provider
func New(logger providers.Logger) (providers.Interface, error) {
	return &providerImpl{Logger: logger}, nil
}

// Destroy is responsible for deleting the cluster
func (p *providerImpl) Destroy(ctx context.Context, name string) error {
	found, err := p.Has(ctx, name)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}
	args := []string{
		"delete",
		"cluster",
		"--name", name,
	}
	p.Info("Deleting the kind cluster: %q", name)

	return exec.CommandContext(ctx, p.path, args...).Run()
}

// Create is responsible for provisioning a kind cluster
func (p *providerImpl) Create(ctx context.Context, name string) error {
	found, err := p.Has(ctx, name)
	if err != nil {
		return err
	}
	if found {
		p.Info("Kind cluster: %q already exists, skipping creation", name)

		return p.ensureRunning(ctx, name)
	}

	args := []string{
		"create",
		"cluster",
		"--name", name,
		"--wait", "10m",
		"--config=-",
	}
	p.Info("Provisioning a kind cluster: %q (usually takes 3-5mins)", name)

	cmd := exec.CommandContext(ctx, p.path, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if _, err := io.WriteString(stdin, KindConfiguration); err != nil {
		return err
	}
	stdin.Close()

	if combined, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s", combined)
	}

	return nil
}

// Export is responsible for exporting the kind kubeconfig
func (p *providerImpl) Export(ctx context.Context, name string) (string, error) {
	args := []string{
		"export",
		"kubeconfig",
		"--name", name,
	}
	p.Info("Exporting kubeconfig from kind cluster: %q", name)

	combined, err := exec.CommandContext(ctx, p.path, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("trying to export kubeconfig: %s", combined)
	}

	contextName := fmt.Sprintf("kind-%s", name)

	return contextName, nil
}

// Has checks if a kind cluster already exists
func (p *providerImpl) Has(ctx context.Context, name string) (bool, error) {
	args := []string{
		"get",
		"clusters",
	}
	p.Info("Checking if kind cluster: %q already exists", name)

	combined, err := exec.CommandContext(ctx, p.path, args...).CombinedOutput()
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(combined))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), name) {
			return true, nil
		}
	}

	return false, nil
}

// Stop is called to stop the kind instance
func (p *providerImpl) Stop(ctx context.Context, name string) error {
	found, err := p.Has(ctx, name)
	if err != nil {
		return err
	}
	if !found {
		return nil
	}

	args := []string{
		"stop",
		"kore-control-plane",
	}
	p.Info("Ensuring the kind cluster: %q is stopped", name)

	path, err := exec.LookPath("docker")
	if err != nil {
		return errors.New("missing binary: docker in $PATH")
	}

	return exec.CommandContext(ctx, path, args...).Run()
}

func (p *providerImpl) Preflight(ctx context.Context) error {
	path, err := exec.LookPath("kind")
	if err != nil {
		return errors.New("missing binary: kind in $PATH")
	}
	p.path = path

	path, err = exec.LookPath("docker")
	if err != nil {
		return errors.New("missing binary: docker in $PATH")
	}

	return nil
}

func (p *providerImpl) ensureRunning(ctx context.Context, name string) error {
	args := []string{
		"start",
		"kore-control-plane",
	}
	p.Info("Ensuring the kind cluster: %q is running", name)

	path, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	return exec.CommandContext(ctx, path, args...).Run()
}
