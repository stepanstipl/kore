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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlural(t *testing.T) {
	cases := []struct {
		Name     string
		Expected string
	}{
		{
			Name:     "gke",
			Expected: "gkes",
		},
		{
			Name:     "ingress",
			Expected: "ingresses",
		},
		{
			Name:     "kubernetes",
			Expected: "kubernetes",
		},
		{
			Name:     "gkecredentials",
			Expected: "gkecredentials",
		},
		{
			Name:     "team",
			Expected: "teams",
		},
		{
			Name:     "teams",
			Expected: "teams",
		},
		{
			Name:     "managedpodsecuritypolicys",
			Expected: "managedpodsecuritypolicies",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.Expected, ToPlural(c.Name))
	}
}
