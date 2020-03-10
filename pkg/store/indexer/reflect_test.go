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

package indexer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReflectQuery(t *testing.T) {
	cs := []struct {
		Query    document
		Expected string
	}{
		{},
		{
			Query: document{
				Age:          10,
				Kind:         "pods",
				IgnoreLabels: map[string]int{"test": 1},
				Labels: map[string]string{
					"env": "prod",
				},
				Name: "test1",
			},
			Expected: "+age:10 +kind:pods +labels.env:prod +name:test1",
		},
	}
	for _, c := range cs {
		query, err := buildReflectQuery(&c.Query)
		require.NoError(t, err)
		assert.Equal(t, c.Expected, query)
	}
}

func BenchmarkReflectQuery(b *testing.B) {
	doc := &document{
		Age:       10,
		Kind:      "pods",
		Labels:    map[string]string{"env": "prod", "live": "true"},
		Name:      "test1",
		Namespace: "default",
	}

	for n := 0; n < b.N; n++ {
		_, _ = buildReflectQuery(doc)
	}
}
