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
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGetRuntimeField(t *testing.T) {
	value := ""
	expected := "test"

	cases := []struct {
		Object   runtime.Object
		Path     string
		Expected interface{}
		Value    interface{}
	}{
		{
			Object: (&v1.Pod{
				Status: v1.PodStatus{
					Message:           "test",
					NominatedNodeName: "test",
				},
			}),
			Path:     "status",
			Expected: &v1.PodStatus{Message: "test", NominatedNodeName: "test"},
			Value:    &v1.PodStatus{},
		},
		{
			Object: (&v1.Pod{
				Status: v1.PodStatus{
					Message:           "test",
					NominatedNodeName: "test",
				},
			}),
			Path:     "status.message",
			Expected: &expected,
			Value:    &value,
		},
	}
	for _, c := range cases {
		assert.NoError(t, GetRuntimeField(c.Object, c.Path, c.Value))
		assert.Equal(t, c.Expected, c.Value)
	}
}
