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

func TestGetSubnetsFromCidr(t *testing.T) {
	type SubNetTest struct {
		Cidr          string
		BitMaskSize   int
		Count         int
		ExpectedCidrs []string
		ExpectedErr   error
	}
	tests := []SubNetTest{
		{
			Cidr:        "10.0.0.0/16",
			BitMaskSize: 24,
			Count:       3,
			ExpectedCidrs: []string{
				"10.0.0.0/24",
				"10.0.1.0/24",
				"10.0.2.0/24",
			},
		},
		{
			Cidr:        "192.168.0.0/16",
			BitMaskSize: 18,
			Count:       2,
			ExpectedCidrs: []string{
				"192.168.0.0/18",
				"192.168.64.0/18",
			},
		},
	}
	for _, a := range tests {
		nets, err := GetSubnetsFromCidr(a.Cidr, a.BitMaskSize, a.Count)
		assert.Equal(t, a.ExpectedErr, err)
		assert.Len(t, nets, a.Count)
		for i, net := range nets {
			assert.Equal(t, a.ExpectedCidrs[i], net.String())
		}
	}
}
