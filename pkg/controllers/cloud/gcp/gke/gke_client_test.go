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

package gke

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpgradeRequired(t *testing.T) {
	cases := []struct {
		Desired  string
		Current  string
		Expected bool
	}{
		{"", "1.15.9-gke.9", false},
		{"-", "1.15.9-gke.9", false},
		{"latest", "1.15.9-gke.9", false},
		{"1.14", "1.15.9-gke.9", false},
		{"1.15", "1.15.9-gke.9", false},
		{"1.15.1", "1.15.9-gke.9", false},
		{"1.15.9", "1.15.9-gke.9", false},
		{"1.15.9-gke.9", "1.15.9-gke.9", false},
		{"1.16", "1.15.9-gke.9", true},
		{"1.15.10", "1.15.9-gke.9", true},
		{"1.15.9-gke.10", "1.15.9-gke.9", true},
	}
	for i, c := range cases {
		upgrade, _ := UpgradeRequired(c.Current, c.Desired)
		assert.Equal(t, c.Expected, upgrade, "case %d, upgrade: %s to %s (expected: %t, actual: %t) not as expected", i, c.Current, c.Desired, c.Expected, upgrade)
	}
}
