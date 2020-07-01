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

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestMapCloudSuccess(t *testing.T) {
	cases := []struct {
		cloud         string
		shouldSucceed bool
		expCloud      string
		expK8s        string
	}{
		{
			"gcp",
			true,
			"google",
			"gke",
		},
		{
			"aws",
			true,
			"amazon",
			"eks",
		},
		{
			"azure",
			true,
			"azure",
			"aks",
		},
		{
			"horse",
			false,
			"",
			"",
		},
	}
	for _, c := range cases {
		cicloud, service, err := mapCloud(c.cloud)
		if c.shouldSucceed {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
		assert.Equal(t, c.expCloud, cicloud)
		assert.Equal(t, c.expK8s, service)
	}
}
