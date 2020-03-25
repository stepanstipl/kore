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

package v1

// Ownership indicates the ownership of a resource
// +k8s:openapi-gen=true
type Ownership struct {
	// Group is the api group
	Group string `json:"group"`
	// Version is the group version
	Version string `json:"version"`
	// Kind is the name of the resource under the group
	Kind string `json:"kind"`
	// Namespace is the location of the object
	Namespace string `json:"namespace"`
	// Name is name of the resource
	Name string `json:"name"`
}

func (o Ownership) IsSameType(o2 Ownership) bool {
	return o.Group == o2.Group &&
		o.Version == o2.Version &&
		o.Kind == o2.Kind &&
		o.Namespace == o2.Namespace
}
