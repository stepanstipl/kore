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

package store

import (
	"errors"
	"fmt"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// nameRegex is a regex to validate the naming parameter
	nameRegex = regexp.MustCompile(`^[a-zA-Z0-9\-_\*]*$`)
)

// queryBuilder is a query builer for the store's index
type queryBuilder struct {
	// the store contract
	store *storeImpl
	// query is the built query
	query *Query
}

// newQueryBuilder returns a query builder for us
func newQueryBuilder(store *storeImpl) *queryBuilder {
	return &queryBuilder{
		store: store,
		query: &Query{
			Labels: make(map[string]string),
		},
	}
}

// APIVersion adds the api version
func (q *queryBuilder) APIVersion(name string) Interface {
	q.query.APIVersion = name

	return q
}

// UID set the object uid
func (q *queryBuilder) UID(name string) Interface {
	q.query.UID = name

	return q
}

// Name sets the name of the resource
func (q *queryBuilder) Name(name string) Interface {
	q.query.Name = name

	return q
}

// Label adds a label to the query
func (q *queryBuilder) Label(k, v string) Interface {
	q.query.Labels[k] = v

	return q
}

// Labels sets all labels as one
func (q *queryBuilder) Labels(v map[string]string) Interface {
	q.query.Labels = v

	return q
}

// Namespace sets the query namespace
func (q *queryBuilder) Namespace(name string) Interface {
	q.query.Namespace = name

	return q
}

// Kind sets the query resource kind
func (q *queryBuilder) Kind(name string) Interface {
	q.query.Kind = name

	return q
}

// Delete removes a object from the store
func (q *queryBuilder) Delete(name string) error {
	q.query.Name = name

	if err := q.query.IsValid(); err != nil {
		return err
	}

	return q.store.deleteObjectStore(q.query)
}

// Set adds an object to the store
func (q *queryBuilder) Set(name string, o metav1.Object) error {
	q.query.Name = name

	if err := q.query.IsValid(); err != nil {
		return err
	}

	return q.store.updateObjectStore(q.query, o)
}

// Has checks if the resource exists in the store
func (q *queryBuilder) Has(name string) (bool, error) {
	e, err := q.Get(name)
	if err != nil {
		return false, err
	}

	return e != nil, nil
}

// Get retrieves a resource from the store
func (q *queryBuilder) Get(name string) (metav1.Object, error) {
	q.query.Name = name

	if err := q.query.IsValid(); err != nil {
		return nil, err
	}

	items, err := q.store.searchObjectStore(q.query)
	if err != nil {
		return nil, err
	}
	if len(items) > 1 {
		return nil, errors.New("too many results returns for query")
	}
	if len(items) == 0 {
		return nil, nil
	}

	return items[0], nil
}

// List retrieves a list of resources from the store
func (q *queryBuilder) List() ([]metav1.Object, error) {
	if err := q.query.IsValid(); err != nil {
		return nil, err
	}

	return q.store.searchObjectStore(q.query)
}

func (q *queryBuilder) Client() Client {
	return q.store.Client()
}

// IsValid checks if the query is valid
func (f *Query) IsValid() error {
	if f.Kind == "" {
		return errors.New("resource kind not set")
	}
	if f.Namespace != "" && !nameRegex.MatchString(f.Namespace) {
		return fmt.Errorf("namespace: %q is invalid", f.Namespace)
	}
	if f.Kind != "" && !nameRegex.MatchString(f.Kind) {
		return fmt.Errorf("kind: %q is invalid", f.Kind)
	}
	if f.Name != "" && !nameRegex.MatchString(f.Name) {
		return fmt.Errorf("name: %q is invalid", f.Name)
	}

	return nil
}
