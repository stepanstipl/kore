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

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/appvia/kore/pkg/schema"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	rschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewMemClient(scheme *runtime.Scheme, objects ...runtime.Object) (*MemClient, error) {
	m := &MemClient{
		scheme:  scheme,
		objects: map[rschema.GroupVersionKind]map[types.NamespacedName][]byte{},
	}

	for _, o := range objects {
		if err := m.Create(context.Background(), o); err != nil {
			return nil, err
		}
	}

	return m, nil
}

type MemClient struct {
	objects map[rschema.GroupVersionKind]map[types.NamespacedName][]byte
	scheme  *runtime.Scheme
}

func (m MemClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	gvk, err := getGVK(obj)
	if err != nil {
		return err
	}

	raw, found := m.objects[gvk][key]
	if !found {
		return apierrors.NewNotFound(rschema.GroupResource{
			Group:    gvk.Group,
			Resource: strings.ToLower(gvk.Kind),
		}, key.Name)
	}

	if _, err := schema.DecodeJSON(raw, obj); err != nil {
		return err
	}

	return nil
}

func (m MemClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	gvk, err := getGVK(list)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(gvk.Kind, "List") {
		return fmt.Errorf("non-list type %T (kind %q) passed as output", list, gvk)
	}

	gvk.Kind = gvk.Kind[:len(gvk.Kind)-4]

	var res []runtime.Object
	for _, raw := range m.objects[gvk] {
		o, err := schema.DecodeJSON(raw, nil)
		if err != nil {
			return err
		}
		res = append(res, o)
	}

	listOpts := client.ListOptions{}
	listOpts.ApplyOptions(opts)
	if listOpts.LabelSelector != nil {
		if res, err = filterWithLabels(res, listOpts.LabelSelector); err != nil {
			return err
		}
	}

	if err := meta.SetList(list, res); err != nil {
		return err
	}

	return nil
}

func (m MemClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	gvk, err := getGVK(obj)
	if err != nil {
		return err
	}

	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	if accessor.GetName() == "" {
		return apierrors.NewInvalid(
			gvk.GroupKind(),
			accessor.GetName(),
			field.ErrorList{field.Required(field.NewPath("metadata.name"), "name is required")})
	}
	if accessor.GetResourceVersion() != "" {
		return apierrors.NewBadRequest("resourceVersion can not be set for Create requests")
	}
	accessor.SetResourceVersion("1")

	raw, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	if _, ok := m.objects[gvk]; !ok {
		m.objects[gvk] = map[types.NamespacedName][]byte{}
	}
	m.objects[gvk][types.NamespacedName{
		Namespace: accessor.GetNamespace(),
		Name:      accessor.GetNamespace(),
	}] = raw

	return nil
}

func (m MemClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	panic("implement me")
}

func (m MemClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	panic("implement me")
}

func (m MemClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	panic("implement me")
}

func (m MemClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	panic("implement me")
}

func (m MemClient) Status() client.StatusWriter {
	return memStatusClient{client: m}
}

type memStatusClient struct {
	client MemClient
}

func (m memStatusClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	panic("implement me")
}

func (m memStatusClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	panic("implement me")
}

func getGVK(obj runtime.Object) (rschema.GroupVersionKind, error) {
	gvk, found, err := schema.GetGroupKindVersion(obj)
	if err != nil {
		return rschema.GroupVersionKind{}, err
	}

	if !found {
		return rschema.GroupVersionKind{}, fmt.Errorf("GVK not found for object %T", obj)
	}

	return gvk, nil
}

func filterWithLabels(objs []runtime.Object, selector labels.Selector) ([]runtime.Object, error) {
	res := make([]runtime.Object, 0, len(objs))
	for _, obj := range objs {
		m, err := meta.Accessor(obj)
		if err != nil {
			return nil, err
		}
		if !selector.Matches(labels.Set(m.GetLabels())) {
			continue
		}
		res = append(res, obj)
	}
	return res, nil
}
