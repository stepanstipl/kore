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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// searchObjectStore is resposible for handling a search of the store for objects
func (s *storeImpl) searchObjectStore(query *Query) ([]metav1.Object, error) {
	timed := prometheus.NewTimer(getLatency)
	defer timed.ObserveDuration()

	// @step: query the index for the objects
	hits, err := s.search.Query(query)
	if err != nil {
		return []metav1.Object{}, err
	}

	var list []metav1.Object

	for _, x := range hits {
		if object, found := s.cache.Load(x); found {
			list = append(list, object.(metav1.Object))
			continue
		}

		log.WithFields(log.Fields{
			"query": query,
		}).Warn("cache key was not found in local store cache")

		errorCounter.Inc()
	}

	return list, nil
}
