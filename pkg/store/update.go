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
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// updateObjectStore is responsible for updating the store
func (s *storeImpl) updateObjectStore(q *Query, o metav1.Object) error {
	timed := prometheus.NewTimer(setLatency)
	defer timed.ObserveDuration()

	err := func() error {
		uid := fmt.Sprintf("%s/%s/%s/%s/%s", q.UID, q.APIVersion, q.Kind, q.Namespace, q.Name)

		// @step: create a document for indexing
		if err := s.search.Index(uid, q); err != nil {
			return err
		}
		s.cache.Store(uid, o)

		return nil
	}()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("unable to add the document into search index")

		errorCounter.Inc()
	}

	return err
}
