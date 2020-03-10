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

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// handleObject is called by an informer when an resource has been created, deleted or updated
func (s *storeImpl) handleObject(version schema.GroupVersionResource, before, object metav1.Object, eventType EventType) {
	fields := log.Fields{
		"event":     string(eventType),
		"group":     version.Group,
		"kind":      version.Resource,
		"name":      object.GetName(),
		"namespace": object.GetNamespace(),
		"resource":  object.GetResourceVersion(),
		"version":   version.Version,
	}
	log.WithFields(fields).Debug("processing incoming informer event")

	dc, ok := object.(runtime.Object)
	if !ok {
		log.WithFields(log.Fields{
			"object": object,
		}).Fatal("object does not support the runtime.Objects interface")
	}
	group := dc.GetObjectKind().GroupVersionKind()
	apiGroup := fmt.Sprintf("%s/%s", group.Group, group.Version)

	switch eventType {
	case Created:
		createCounter.Inc()
		fallthrough
	case Updated:
		updateCounter.Inc()
		if err := s.APIVersion(apiGroup).Namespace(object.GetNamespace()).Kind(group.Kind).Set(object.GetName(), object); err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Error("unable to update or create resource in the store")
		}
	default:
		deleteCounter.Inc()
		if err := s.APIVersion(apiGroup).Namespace(object.GetNamespace()).Kind(group.Kind).Delete(object.GetName()); err != nil {
			fields["error"] = err.Error()
			log.WithFields(fields).Error("unable to delete from the store")
		}
	}
	s.handleEventListeners(version, before, object, eventType)
}

// handleEventListeners is responsible for handling the event listeners
func (s *storeImpl) handleEventListeners(version schema.GroupVersionResource, before, object metav1.Object, eventType EventType) {
	// @step: iterate the listeners and fire off the callback
	for _, x := range s.getWatchers(version) {
		go func(l *Listener) {
			switch eventType {
			case Created:
				l.Created(object)
			case Deleted:
				l.Deleted(object)
			case Updated:
				l.Updated(before, object)
			}
		}(x)
	}
}
