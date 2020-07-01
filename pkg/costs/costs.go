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

package costs

type Costs interface {
	// Metadata returns the interface to the pricing metadata service
	Metadata() Metadata
	// Estimates returns the interface to the cost estimation service
	Estimates() Estimates
	// Actuals returns the interface to the actual costs service
	Actuals() Actuals
}

// New returns a new instance of the costs API
func New(config *Config) Costs {
	cloudinfo := NewCloudInfo(config.CloudinfoURL)
	metadata := NewMetadata(cloudinfo)
	estimates := NewEstimates(metadata)
	actuals := NewActuals()
	return &costsImpl{
		metadata,
		estimates,
		actuals,
	}
}

var _ Costs = &costsImpl{}

type costsImpl struct {
	metadata  Metadata
	estimates Estimates
	actuals   Actuals
}

func (c *costsImpl) Metadata() Metadata {
	return c.metadata
}

func (c *costsImpl) Estimates() Estimates {
	return c.estimates
}

func (c *costsImpl) Actuals() Actuals {
	return c.actuals
}
