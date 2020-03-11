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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type document struct {
	Age          int               `json:"age"`
	Kind         string            `json:"kind"`
	IgnoreLabels map[string]int    `json:"ignore-labels"`
	Labels       map[string]string `json:"labels"`
	Modified     string            `json:"modified"`
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
}

func newTestIndex(t *testing.T) *indexer {
	i, err := New()
	require.NotNil(t, i)
	require.NoError(t, err)

	namespaces := []string{"default", "test", "frontend"}
	services := []string{"svc1", "svc2"}
	pods := []string{"test1", "test2"}
	things := []string{"thing1", "thing2"}

	for j, x := range namespaces {
		id := fmt.Sprintf("namespace-%d", j)
		require.NoError(t, i.Index(id, &document{Kind: "namespace", Name: x}))
	}

	for _, n := range namespaces {
		for j, p := range pods {
			id := fmt.Sprintf("pod-%s-%d", n, j)
			require.NoError(t, i.Index(id, &document{Kind: "pod", Name: p, Namespace: n}))
		}
	}

	for _, n := range namespaces {
		for j, p := range services {
			id := fmt.Sprintf("svc-%s-%d", n, j)
			require.NoError(t, i.Index(id, &document{Kind: "service", Name: p, Namespace: n}))
		}
	}

	for _, n := range things {
		require.NoError(t, i.Index(n, &document{
			Kind:   "thing",
			Labels: map[string]string{"class": n},
		}))
	}

	return i.(*indexer)
}

func TestNew(t *testing.T) {
	i, err := New()
	assert.NotNil(t, i)
	assert.NoError(t, err)
}

func TestSize(t *testing.T) {
	i := newTestIndex(t)
	size, err := i.Size()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0x11), size)
}

func TestIndex(t *testing.T) {
	_ = newTestIndex(t)
}

func TestDelete(t *testing.T) {
	i := newTestIndex(t)
	hits, err := i.QueryRaw("+name:test1 +namespace:default +kind:pod")
	require.NoError(t, err)
	require.NotEmpty(t, hits)
	require.Equal(t, 1, len(hits))

	require.NoError(t, i.Delete(hits[0]))
	hits, err = i.QueryRaw("+name:test1 +namespace:default +kind:pod")
	require.NoError(t, err)
	require.Empty(t, hits)
	require.Equal(t, 0, len(hits))
}

func TestDeleteByQuery(t *testing.T) {
	i := newTestIndex(t)
	hits, err := i.QueryRaw("+name:test1 +namespace:default +kind:pod")
	require.NoError(t, err)
	require.NotEmpty(t, hits)
	require.Equal(t, 1, len(hits))

	filter := &document{
		Kind:      "pod",
		Namespace: "default",
		Name:      "test1",
	}
	num, err := i.DeleteByQuery(filter)
	require.Equal(t, 1, num)
	require.NoError(t, err)

	hits, err = i.QueryRaw("+name:test1 +namespace:default +kind:pod")
	require.NoError(t, err)
	require.Empty(t, hits)
	require.Equal(t, 0, len(hits))
}

func TestDeleteNamespaceByQuery(t *testing.T) {
	i := newTestIndex(t)
	query := "+namespace:default"
	hits, err := i.QueryRaw(query)
	require.NoError(t, err)
	require.NotEmpty(t, hits)
	require.Equal(t, 4, len(hits))

	num, err := i.DeleteByQueryRaw(query)
	require.NoError(t, err)
	require.Equal(t, 4, num)
	hits, err = i.QueryRaw(query)
	require.NoError(t, err)
	require.Empty(t, hits)
	require.Equal(t, 0, len(hits))
}

func TestQueries(t *testing.T) {
	cs := []struct {
		Query         *document
		RawQuery      string
		ExpectedCount int
	}{
		{
			ExpectedCount: 0,
			RawQuery:      "",
			Query:         &document{},
		},
		{
			ExpectedCount: 6,
			Query:         &document{Kind: "pod[s]"},
			RawQuery:      "+kind:pod[s]",
		},
		{
			ExpectedCount: 1,
			Query:         &document{Labels: map[string]string{"class": "thing1"}},
			RawQuery:      "+labels.class:thing1",
		},
		{
			ExpectedCount: 2,
			Query:         &document{Labels: map[string]string{"class": "*"}},
			RawQuery:      "+labels.class:*",
		},
		{
			ExpectedCount: 2,
			Query:         &document{Kind: "thing*", Labels: map[string]string{"class": "*"}},
			RawQuery:      "+kind:thing* +labels.class:*",
		},
		{
			ExpectedCount: 2,
			Query:         &document{Namespace: "default", Kind: "pod[s]"},
			RawQuery:      "+namespace:default +kind:pod[s]",
		},
		{
			ExpectedCount: 1,
			Query:         &document{Namespace: "default", Kind: "pod[s]", Name: "test1"},
			RawQuery:      "+namespace:default +kind:pod[s] +name=test1",
		},
		{
			ExpectedCount: 0,
			Query:         &document{Namespace: "default", Kind: "pod[s]", Name: "test0"},
			RawQuery:      "+namespace:default +kind:pod[s] +name=test0",
		},
		{
			ExpectedCount: 0,
			Query:         &document{Namespace: "default", Modified: fmt.Sprintf(">%d", time.Now().Add(-1*time.Minute).Unix())},
			RawQuery:      "+namespace:default +modified:>" + fmt.Sprintf("%d", time.Now().Add(-1*time.Minute).Unix()),
		},
	}
	s := newTestIndex(t)
	for i, c := range cs {
		hits, err := s.QueryRaw(c.RawQuery)
		require.NoError(t, err)
		require.Equal(t, c.ExpectedCount, len(hits), "case %d, raw: %s, expected: %d, got: %d",
			i, c.RawQuery, c.ExpectedCount, len(hits))

		hits, err = s.Query(c.Query)
		require.NoError(t, err)
		require.Equal(t, c.ExpectedCount, len(hits), "case %d, query: %v, expected: %d, got: %d",
			i, c.Query, c.ExpectedCount, len(hits))
	}
}
