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
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// deleteObjectStore is responsible for deleting an object from the store
func (s *storeImpl) deleteObjectStore(query *Query) error {
	timed := prometheus.NewTimer(deleteLatency)
	defer timed.ObserveDuration()

	err := func() error {
		var hits []string
		// @step: search for object matching the query
		if hits, err := s.search.Query(query); err != nil {
			return err
		} else if len(hits) <= 0 {
			return nil
		}

		// @step: delete the items from the search index
		if _, err := s.search.DeleteByQuery(query); err != nil {
			return err
		}

		// @step: delete items from the cache
		for _, x := range hits {
			s.cache.Delete(x)
		}

		return nil
	}()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to delete from object store")

		errorCounter.Inc()
	}

	return err
}
