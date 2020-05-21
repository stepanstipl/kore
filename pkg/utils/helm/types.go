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

package helm

import (
	"context"

	"helm.sh/helm/v3/pkg/chart"
)

// Interface is a generic interface for helm client
type Interface interface {
	// Deploy is responsible for deploying the chart
	Deploy(ctx context.Context, chart *chart.Chart)
	// GetReleases returns a list of releases
	GetReleases(ctx context.Context, namespace string) ([]chart.Chart, error)
	// LoadChart loads a chart for us
	LoadChart(ctx context.Context, path string) (*chart.Chart, error)
}

type DeploymentOptions interface {
	// Namespace is the name of the namespace to deploy
	Namespace(string) DeploymentOptions
	// Name is the release name
	Name(string) DeploymentOptions
	// Wait indicates if we should wait for it
	Wait(bool) DeploymentOptions
}
