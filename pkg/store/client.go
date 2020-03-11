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
	"context"
	"errors"
	"fmt"
	"reflect"

	hubschema "github.com/appvia/kore/pkg/schema"

	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// rclient is a ephermal struct to carry the
type rclient struct {
	store Store
	// client is the controller-runtime client to use
	client client.Client
	// index is the search query
	index *queryBuilder
	// value in the object we are mapping into
	value runtime.Object
	// current is the state of the current value
	current runtime.Object
	// useCache indicates we can use the index
	useCache bool
	// withCreate indicates if on update we check for exists
	withCreate bool
	// withPatch indicates we should check if there, patch
	withPatch bool
	// withForceApply indicates we should force apply the updates
	withForceApply bool
}

// Get is responsible for retrieving an object from the api / index
func (r *rclient) Get(ctx context.Context, options ...GetOptionFunc) error {
	// @step: apply the get options
	GetOptions.apply(r, options...)

	// @step: check we have what we need
	if r.value == nil {
		return errors.New("expected an runtime.Object set")
	}
	if r.index.query.Name == "" {
		return errors.New("no name set for object")
	}

	// @step: should we check the cache
	if r.useCache {
		found, err := func() (bool, error) {
			// @step: update the query from the object
			if _, err := r.updateQueryFromObject(r.value); err != nil {
				return false, err
			}
			// @step: query the cache for an answer
			object, err := r.index.Get(r.index.query.Name)
			if err != nil {
				return false, err
			}
			if object != nil {
				cacheHitCounter.Inc()
				if err := ObjectToType(r.value, object); err != nil {
					return false, err
				}
				return true, nil
			}
			cacheMissCounter.Inc()

			return true, nil
		}()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("failed to perform get request on the cache")
		}
		if found {
			return nil
		}
	}

	// @step create a object reference
	reference := types.NamespacedName{
		Namespace: r.index.query.Namespace,
		Name:      r.index.query.Name,
	}

	if err := r.client.Get(ctx, reference, r.value); err != nil {
		return err
	}

	// @step: retrieve the object the kube-api
	return r.client.Get(ctx, reference, r.value)
}

// Create is responsible for creating an object in the api - enuring it doesn't exist already
func (r *rclient) Create(ctx context.Context, options ...CreateOptionFunc) error {
	// @step: apply the options
	CreateOptions.apply(r, options...)

	if r.value == nil {
		return errors.New("no value has been set")
	}

	// @step: attempt to inject the object into the api
	err := r.client.Create(ctx, r.value)
	if err != nil {
		return err
	}

	// @step: check if the placeholder is an typed or unstructured list and
	// if so, we can just return
	if reflect.ValueOf(r.value).Elem().Type() == reflect.TypeOf(unstructured.Unstructured{}) {
		return nil
	}

	// @step: attempt to inject the resource direct
	err = func() error {
		object, err := r.updateQueryFromObject(r.value)
		if err != nil {
			return err
		}
		if err := r.index.Set(object.GetName(), object); err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to update internal cache on create")
	}

	return nil
}

// Delete is responsible for deleting an object from the index / api
func (r *rclient) Delete(ctx context.Context, options ...DeleteOptionFunc) error {
	// @step: apply the options
	DeleteOptions.apply(r, options...)

	if r.value == nil {
		return errors.New("no value has been set")
	}

	return r.client.Delete(ctx, r.value)
}

// DeleteAll is responsible for deleting a series of objects
func (r *rclient) DeleteAll(ctx context.Context, options ...DeleteAllOptionFunc) error {
	// @step: apply the options
	DeleteAllOptions.apply(r, options...)

	if r.value == nil {
		return errors.New("no value has been set")
	}
	if r.index.query.Namespace == "" {
		return errors.New("no namespace set")
	}

	var opts []client.DeleteAllOfOption
	if r.index.query.Namespace != "" {
		opts = append(opts, client.InNamespace(r.index.query.Namespace))
	}
	if len(r.index.query.Labels) > 0 {
		opts = append(opts, client.MatchingLabels(r.index.query.Labels))
	}
	if _, err := r.updateQueryFromObject(r.value); err != nil {
		return err
	}

	return r.client.DeleteAllOf(ctx, r.value, opts...)
}

// List is responsible for listing all the objects
func (r *rclient) List(ctx context.Context, options ...ListOptionFunc) error {
	// @step: apply the options
	ListOptions.apply(r, options...)

	if r.value == nil {
		return errors.New("no runtime.Object has been set")
	}

	// @step: should we check the cache
	if r.useCache {
		found, err := func() (bool, error) {
			// @step: update the query from the object
			if err := r.updateQueryFromList(r.value); err != nil {
				return false, err
			}
			// @step: query the cache for answers
			objects, err := r.index.List()
			if err != nil {
				return false, err
			}
			if len(objects) > 0 {
				cacheHitCounter.Inc()
				// @step: we reflect the cached objects
				if err := ObjectsToList(r.value, objects); err != nil {
					return false, err
				}
				return true, nil
			}
			cacheMissCounter.Inc()

			return true, nil
		}()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("failed to perform get request on the cache")
		}
		if found {
			return nil
		}
	}

	// @step: make the request direct to the kube api
	var opts []client.ListOption
	if r.index.query.Namespace != "" {
		opts = append(opts, client.InNamespace(r.index.query.Namespace))
	}
	if len(r.index.query.Labels) > 0 {
		opts = append(opts, client.MatchingLabels(r.index.query.Labels))
	}

	return r.client.List(ctx, r.value, opts...)
}

// Has is responsible for checking a resource exists
func (r *rclient) Has(ctx context.Context, options ...HasOptionFunc) (bool, error) {
	// @step: apply the options
	HasOptions.apply(r, options...)

	if err := r.Get(ctx); err != nil {
		if kerrors.IsNotFound(err); err != nil {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Update is responsible for either creating or updating the resource on patch
func (r *rclient) Update(ctx context.Context, options ...UpdateOptionFunc) error {
	// @step: apply the options
	UpdateOptions.apply(r, options...)

	if r.value == nil {
		return errors.New("no value set for update operation on store")
	}

	key, err := client.ObjectKeyFromObject(r.value)
	if err != nil {
		return err
	}
	if key.Name == "" {
		return errors.New("no name on resource")
	}

	// @step: check if the resource exists already
	original, found, err := func() (runtime.Object, bool, error) {
		current := r.value.DeepCopyObject()
		if err := r.client.Get(ctx, key, current); err != nil {
			if !kerrors.IsNotFound(err) {
				return nil, false, err
			}

			return nil, false, nil
		}

		return current, true, nil
	}()
	if err != nil {
		return err
	}

	if !found && r.withCreate {
		return r.client.Create(ctx, r.value)
	}
	if r.withForceApply {
		r.value.(metav1.Object).SetGeneration(original.(metav1.Object).GetGeneration())
		r.value.(metav1.Object).SetResourceVersion(original.(metav1.Object).GetResourceVersion())

		return r.client.Update(ctx, r.value)
	}

	if err := r.client.Patch(ctx, r.value, client.MergeFrom(original)); err != nil {
		return err
	}

	// @step: attempt to inject the resource direct
	err = func() error {
		object, err := r.updateQueryFromObject(r.value)
		if err != nil {
			return err
		}
		if err := r.index.Set(object.GetName(), object); err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to update internal cache on create")
	}

	return nil
}

// updateQueryFromList is responsible for filling in the scheme for query
func (r *rclient) updateQueryFromList(object runtime.Object) error {
	// @step: is the object versioned in our scheme?
	gvk, _, err := hubschema.GetGroupKindVersion(object)
	if err != nil {
		return err
	}
	// @step: make an exception for v1
	if gvk.Group == "" && gvk.Version == "v1" {
		r.index.APIVersion("v1")
	} else {
		r.index.APIVersion(fmt.Sprintf("%s/%s", gvk.Group, gvk.Version))
	}
	r.index.Kind(gvk.Kind[:len(gvk.Kind)-4])

	return nil
}

// updateQueryFromObject is responsible for filling in the scheme for query
func (r *rclient) updateQueryFromObject(object runtime.Object) (metav1.Object, error) {
	// @step: is the object versioned in our scheme?
	gvk, _, err := hubschema.GetGroupKindVersion(object)
	if err != nil {
		return nil, err
	}
	// @step: make an exception for v1
	if gvk.Group == "" && gvk.Version == "v1" {
		r.index.APIVersion("v1")
	} else {
		r.index.APIVersion(fmt.Sprintf("%s/%s", gvk.Group, gvk.Version))
	}
	r.index.Kind(gvk.Kind)

	// @step: we need to extract name, namespace and labels
	dc, ok := object.(metav1.Object)
	if !ok {

		return nil, errors.New("object does not support the metav1.Object interface")
	}
	if r.index.query.Name == "" {
		r.index.Name(dc.GetName())
	}
	if r.index.query.Namespace == "" {
		r.index.Namespace(dc.GetNamespace())
	}

	return dc, nil
}
