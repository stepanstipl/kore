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

package informer

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// resourceInformer is a kubernetes resources informer
type resourceInformer struct {
	// config is the configuration for the service
	config *Config
	// informer is the underlining generic informer
	informer informers.GenericInformer
	// resourceVersion is the resource we are listening to
	resourceVersion schema.GroupVersionResource
	// stopCh is the stop channel for the informer
	stopCh chan struct{}
}

// New creates and returns a resource informer
func New(config *Config) (Informer, error) {
	// @check the resource is supported
	version, found := ResourceVersions()[config.Resource]
	if !found {
		return nil, fmt.Errorf("resource: %s is not a supported resource type", config.Resource)
	}

	var informer informers.GenericInformer

	// @step: we try each of the informers
	for _, x := range config.Factories {
		i, err := x.ForResource(version)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Warn("informer does not supported requested type")

			continue
		}
		informer = i
	}

	// @step: we need to create an dynamic informer for this
	if informer == nil {
		return nil, fmt.Errorf("no informaer found for resource type: %s", version)
	}

	svc := &resourceInformer{
		informer:        informer,
		config:          config,
		resourceVersion: version,
		stopCh:          make(chan struct{}),
	}

	// @step: start up informer and return
	if err := svc.start(); err != nil {
		return nil, err
	}

	return svc, nil
}

// Stop is called to close the informer
func (r *resourceInformer) Stop() error {
	r.stopCh <- struct{}{}

	return nil
}

// start is responsible for running the informing and updating the caches
func (r *resourceInformer) start() error {
	r.informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(before interface{}) {
			addCounter.WithLabelValues(r.config.Resource).Inc()
			r.handleAddObject(before)
		},
		DeleteFunc: func(before interface{}) {
			deleteCounter.WithLabelValues(r.config.Resource).Inc()
			r.handleDeleteObject(before)
		},
		UpdateFunc: func(before, after interface{}) {
			updateCounter.WithLabelValues(r.config.Resource).Inc()
			r.handleUpdateObject(before, after)
		},
	})

	// @step: start the shared index informer
	stopCh := make(chan struct{})
	go r.informer.Informer().Run(stopCh)

	log.Infof("synchronizing the informer cache with resources: %s", r.config.Resource)

	if !cache.WaitForCacheSync(stopCh, r.informer.Informer().HasSynced) {
		runtime.HandleError(fmt.Errorf("controller timed out waiting for caches to sync"))

		return fmt.Errorf("controller timed out waiting for cache sync")
	}

	// @step: wait for a signal to stop
	go func() {
		<-r.stopCh
		log.WithFields(log.Fields{
			"resource": r.config.Resource,
		}).Infof("closing down in the informer for resource")

		close(stopCh)
	}()

	return nil
}

// handleAddObject is responsible for handling the deletions
func (r *resourceInformer) handleAddObject(before interface{}) {
	if r.config.AddFunc == nil {
		return
	}
	if err := func() error {
		object, err := ensureMetaObject(before)
		if err != nil {
			return err
		}
		r.config.AddFunc(r.resourceVersion, object)

		return nil
	}(); err != nil {
		r.handleInformerError(err)
	}
}

// handleDeleteObject is responsible for handling the deletions
func (r *resourceInformer) handleDeleteObject(before interface{}) {
	if r.config.DeleteFunc == nil {
		return
	}
	if err := func() error {
		object, err := ensureMetaObject(before)
		if err != nil {
			return err
		}
		r.config.DeleteFunc(r.resourceVersion, object)

		return nil
	}(); err != nil {
		r.handleInformerError(err)
	}
}

// handleUpdateObject is resposible for handling an updated object
func (r *resourceInformer) handleUpdateObject(before, after interface{}) {
	if r.config.UpdateFunc == nil {
		return
	}
	if err := func() error {
		b, err := ensureMetaObject(before)
		if err != nil {
			return err
		}
		a, err := ensureMetaObject(after)
		if err != nil {
			return err
		}
		r.config.UpdateFunc(r.resourceVersion, b, a)

		return nil
	}(); err != nil {
		r.handleInformerError(err)
	}
}

// handleInformerError is responsible for pushing the error upstream
func (r *resourceInformer) handleInformerError(err error) {
	go func() {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Debug("resource informer has encountered an error")

		errorCounter.WithLabelValues(r.config.Resource).Inc()

		r.config.ErrorFunc(r.resourceVersion, err)
	}()
}

// ensureMetaObject checks to make sure the object is a meta object for us
func ensureMetaObject(object interface{}) (metav1.Object, error) {
	expected := "metav1.Object"

	// @check the object is not nil
	if object == nil {
		return nil, fmt.Errorf("object expected: %q, got: nil", expected)
	}

	// @check the object implements the metav1.Object interface
	if _, ok := object.(metav1.Object); !ok {
		return nil, fmt.Errorf("object not as expected: %q, got: %q",
			expected, reflect.TypeOf(object).String())
	}

	return object.(metav1.Object), nil
}
