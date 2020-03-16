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

package kore

// GCP is the gke interface
type GCP interface {
	// ProjectClaims returns the claims interface

	// Organizations return the organizations interface
}

type gcpImpl struct {
	*cloudImpl
	// team is the request team
	team string
}

/*

// ProjectClaims is responsible for deleting a gke environment
func (h *gcpImpl) ProjectClaims() ProjectClaims {
	return nil
}

// Organizations return the organizations interface
func (h *gcpImpl) Organizations() Organizations {
	return nil
}
*/
