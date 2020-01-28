/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
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
