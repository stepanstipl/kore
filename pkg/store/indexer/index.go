/*
Copyright 2018 Appvia Ltd <info@appvia.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package indexer

import (
	"github.com/blevesearch/bleve"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// indexer is the service wrapper
type indexer struct {
	// store is the index interface
	store bleve.Index
}

// New returns a memory only index
func New() (Interface, error) {
	store, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return nil, err
	}

	return &indexer{store: store}, nil
}

// Delete is responsible for deleting a document from the index
func (i *indexer) Delete(id string) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	return i.store.Delete(id)
}

// DeleteByQuery deletes all the documents which match the query
func (i *indexer) DeleteByQuery(query interface{}) (int, error) {
	// @step: reflect the struct and build a query
	q, err := buildReflectQuery(query)
	if err != nil {
		return 0, err
	}

	return i.DeleteByQueryRaw(q)
}

// DeleteByQueryRaw deletes all the documents which match the query
func (i *indexer) DeleteByQueryRaw(q string) (int, error) {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	hits, err := i.QueryRaw(q)
	if err != nil {
		return 0, err
	}
	if len(hits) <= 0 {
		return 0, nil
	}

	for index, id := range hits {
		if err := i.Delete(id); err != nil {
			return index, err
		}
	}

	return len(hits), nil
}

// Index is responsible is add a document the index
func (i *indexer) Index(id string, doc interface{}) error {
	timed := prometheus.NewTimer(indexLatency)
	defer timed.ObserveDuration()

	log.WithFields(log.Fields{
		"id": id,
	}).Debug("indexing document into store")

	return i.store.Index(id, doc)
}

// Query is responsible for searching the index
func (i *indexer) Query(search interface{}) ([]string, error) {
	// @step: reflect the struct and build a query
	query, err := buildReflectQuery(search)
	if err != nil {
		return []string{}, err
	}

	return i.QueryRaw(query)
}

// QueryRaw is responsible for searching the index directly
func (i *indexer) QueryRaw(query string) ([]string, error) {
	timed := prometheus.NewTimer(searchLatency)
	defer timed.ObserveDuration()

	var list []string

	log.WithFields(log.Fields{
		"query": query,
	}).Debug("searching the index for item")

	resp, err := i.store.Search(bleve.NewSearchRequest(bleve.NewQueryStringQuery(query)))
	if err != nil {
		return list, err
	}

	for _, x := range resp.Hits {
		list = append(list, x.ID)
	}

	return list, nil
}

// Size returns the size of the index
func (i *indexer) Size() (uint64, error) {
	return i.store.DocCount()
}
