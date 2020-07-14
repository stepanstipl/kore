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

package model

// ResourceReference is the reference to the resource type - i.e. cluster, team, projectclaim etc
type ResourceReference struct {
	// ResourceGroup is the group of the resource
	ResourceGroup string `sql:"DEFAULT:''"`
	// ResourceVersion is the version of the resource
	ResourceVersion string `sql:"DEFAULT:''"`
	// ResourceKind is the kind of the resource
	ResourceKind string `sql:"DEFAULT:''"`
	// ResourceNamespace is the namespace of the resource
	ResourceNamespace string `sql:"DEFAULT:''"`
	// ResourceName is the name of the resource
	ResourceName string `sql:"DEFAULT:''"`
}
