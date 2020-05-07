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

package kubernetes

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestPatch(t *testing.T) {
	cases := []struct {
		Object   runtime.Object
		Expected runtime.Object
		Field    string
		Patch    string
	}{
		{
			Object: &v1.Pod{},
			Expected: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "test1",
						},
					},
				},
			},
			Field: "spec",
			Patch: `
			{
				"containers": [{
					"name": "test1"
				}]
			}
			`,
		},
	}
	for _, c := range cases {
		assert.NoError(t, PatchHelper(c.Object, c.Field, strings.NewReader(c.Patch)))
		json.NewEncoder(os.Stdout).Encode(c.Expected)
		json.NewEncoder(os.Stdout).Encode(c.Object)
		assert.Equal(t, c.Expected, c.Object)
	}
}
