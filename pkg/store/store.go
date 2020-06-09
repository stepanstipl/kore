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
	"sync"
	"time"

	"github.com/appvia/kore/pkg/store/indexer"
	"github.com/appvia/kore/pkg/store/informer"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	log "github.com/sirupsen/logrus"
)

// @TODO we need to remove the former interface used to call
// the index directly

// storeImpl is the implementation of the store
type storeImpl struct {
	sync.RWMutex

	// cache is local cache for object
	cache sync.Map
	// client is the kubernetes client to use
	client kubernetes.Interface
	// runtime is the controller-runtime client for raw access to the api
	runtime client.Client
	// factories is a provider of shared informers
	factories []informers.SharedInformerFactory
	// search is the search index for the resources
	search indexer.Interface
	// watchers is a list of consumers for events on resources
	watchers map[string][]*Listener
	// watching is a map of resource to informer currently being watched
	watching map[string]informer.Informer
}

// New creates and returns a resource store: we receive a kubernete clients and a list
// of resource type to watch and add into the store
func New(client kubernetes.Interface, cc client.Client) (Store, error) {
	log.Info("initializing the store data access layer")

	// @step: create an indexer for the resources
	search, err := indexer.New()
	if err != nil {
		return nil, err
	}

	factories := make([]informers.SharedInformerFactory, 0)
	factories = append(factories, informers.NewSharedInformerFactoryWithOptions(client, 30*time.Second))

	// @step: create a the store service
	return &storeImpl{
		cache:     sync.Map{},
		client:    client,
		runtime:   cc,
		factories: factories,
		watchers:  make(map[string][]*Listener),
		watching:  make(map[string]informer.Informer),
		search:    search,
	}, nil
}

// GetFactories returns a list of factories
func (s *storeImpl) GetFactories() []informers.SharedInformerFactory {
	s.RLock()
	defer s.RUnlock()

	return s.factories
}

// AddEventListener is responsible for adding a listener to one or more resources
func (s *storeImpl) AddEventListener(l *Listener) error {
	// @step: default nothing to a wildcard
	if len(l.Resources) == 0 {
		l.Resources = append(l.Resources, "*")
	}

	for _, x := range l.Resources {
		log.WithFields(log.Fields{
			"resource": x,
		}).Debug("adding a event listener for resource")

		if err := s.addWatcher(l, x); err != nil {
			return err
		}
	}

	return nil
}

// WatchResource adds a resource to the list of resources being watched
func (s *storeImpl) WatchResource(resource string) error {
	// @check if the resource is ok and convert to schema
	version, err := informer.ToSchema(resource)
	if err != nil {
		return err
	}

	// @check if the resource is already being watched
	if s.Watching(version) {
		log.WithFields(log.Fields{
			"group":    version.Group,
			"resource": version.Resource,
			"version":  version.Version,
		}).Debug("resource is already being watched")

		return nil
	}

	// @step: we create an informer for this resource type
	inf, err := informer.New(&informer.Config{
		Factories: s.GetFactories(),
		Resource:  resource,

		// add actions methods
		AddFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			s.handleObject(version, nil, object, Created)
		},
		DeleteFunc: func(version schema.GroupVersionResource, object metav1.Object) {
			s.handleObject(version, nil, object, Deleted)
		},
		UpdateFunc: func(version schema.GroupVersionResource, before, after metav1.Object) {
			s.handleObject(version, before, after, Updated)
		},
		// add the downstream error method
		ErrorFunc: func(version schema.GroupVersionResource, err error) {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"group":    version.Group,
				"resource": version.Resource,
				"version":  version.Version,
			}).Error("resource informer has encountered an error while watching")

			errorCounter.Inc()
		},
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err.Error(),
			"group":    version.Group,
			"resource": version.Resource,
			"version":  version.Version,
		}).Error("failed to create resource informer")

		return err
	}

	// @step: add the informer to the list
	s.addWatching(version, inf)

	return nil
}

// Client returns a store query client
func (s *storeImpl) Client() Client {
	return &rclient{
		client: s.runtime,
		index:  newQueryBuilder(s),
		store:  s,
	}
}

// Watching checks if the resource type is being watched
func (s *storeImpl) Watching(version schema.GroupVersionResource) bool {
	s.RLock()
	defer s.RUnlock()

	if _, found := s.watching[informer.NiceVersion(version)]; found {
		return true
	}

	return false
}

// Namespace returns a namespece scoped client request
func (s *storeImpl) Namespace(name string) Interface {
	return newQueryBuilder(s).Namespace(name)
}

func (s *storeImpl) APIVersion(name string) Interface {
	return newQueryBuilder(s).APIVersion(name)
}

// Kind returns a kind scope client request
func (s *storeImpl) Kind(name string) Interface {
	return newQueryBuilder(s).Kind(name)
}

// Stop is responsible for releasing the resources
func (s *storeImpl) Stop() error {
	// @step: start by kill off the all informers
	s.Lock()
	defer s.Unlock()

	for _, x := range s.watching {
		if err := x.Stop(); err != nil {
			return err
		}
	}

	return nil
}

// addWatcher adds an informer to the watcher map
func (s *storeImpl) addWatching(version schema.GroupVersionResource, inf informer.Informer) {
	s.Lock()
	defer s.Unlock()

	s.watching[informer.NiceVersion(version)] = inf
}

// addWatcher adds an listener to the list
func (s *storeImpl) addWatcher(l *Listener, resource string) error {
	if l == nil {
		return errors.New("listener no defined")
	}
	if l.EventHandlers == nil {
		return errors.New("no event handlers defined")
	}

	s.Lock()
	defer s.Unlock()

	s.watchers[resource] = append(s.watchers[resource], l)

	return nil
}

// getWatchers returns the list of listeners on a resource
func (s *storeImpl) getWatchers(version schema.GroupVersionResource) []*Listener {
	s.RLock()
	defer s.RUnlock()

	var listeners []*Listener

	if v, found := s.watchers["*"]; found {
		listeners = append(listeners, v...)
	}

	if v, found := s.watchers[informer.NiceVersion(version)]; found {
		listeners = append(listeners, v...)
	}

	return listeners
}

// RuntimeClient returns with the runtime client
func (s *storeImpl) RuntimeClient() client.Client {
	return s.runtime
}
