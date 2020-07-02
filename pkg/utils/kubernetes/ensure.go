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

package kubernetes

import (
	"context"
	"errors"
	"fmt"

	"github.com/appvia/kore/pkg/utils/jsonutils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/banzaicloud/k8s-objectmatcher/patch"
	log "github.com/sirupsen/logrus"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckIfExists checks if a runtime object exist in the api
func CheckIfExists(ctx context.Context, cc client.Client, object runtime.Object) (bool, error) {
	return GetIfExists(ctx, cc, object.DeepCopyObject())
}

// GetIfExists retrieves an object if it exists
func GetIfExists(ctx context.Context, cc client.Client, object runtime.Object) (bool, error) {
	key, err := client.ObjectKeyFromObject(object)
	if err != nil {
		return false, err
	}

	if err := cc.Get(ctx, key, object); err != nil {
		if !kerrors.IsNotFound(err) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

// DeleteIfExists removes the resource if it exists
func DeleteIfExists(ctx context.Context, cc client.Client, object runtime.Object) error {
	if err := cc.Delete(ctx, object); err != nil {
		if !kerrors.IsNotFound(err) {
			return err
		}

		return nil
	}

	return nil
}

// CreateOrUpdateObject is shorthand for the below object we dont return the object though
func CreateOrUpdateObject(ctx context.Context, cc client.Client, object runtime.Object) error {
	_, err := CreateOrUpdate(ctx, cc, object)

	return err
}

// CreateOrUpdate is responsible for updating a resource
// extended to carry out updates on specific types when patching is not suitable
func CreateOrUpdate(ctx context.Context, cc client.Client, object runtime.Object) (runtime.Object, error) {
	supported, object, err := TypeSpecificUpdate(ctx, cc, object)
	if supported {
		return object, err
	}

	original := object.DeepCopyObject()

	existing, err := GetIfExists(ctx, cc, original)
	if err != nil {
		return nil, err
	}

	if existing {
		nobj, ok := object.(metav1.Object)
		if !ok {
			return nil, errors.New("object does not implement the metav1.Object")
		}
		old, ok := original.(metav1.Object)
		if !ok {
			return nil, errors.New("original object does not implement the metav1.Object")
		}
		nobj.SetResourceVersion(old.GetResourceVersion())

		return object, cc.Patch(ctx, object, client.MergeFrom(original))
	}

	// @step: we first try and create the resource
	return object, cc.Create(ctx, object)
}

// PatchOrReplace is responsible for updating a resouce and optionally updating if patching fails
// Due to issues with some immutable objects set on the server, this will allow a fall back for now
// See similar related issue with patching in Helm: https://github.com/helm/helm/issues/7516
func PatchOrReplace(ctx context.Context, cc client.Client, object runtime.Object) (runtime.Object, error) {
	// need to create or replace for deployments and services
	gvk := object.GetObjectKind().GroupVersionKind()
	log.Debugf("deciding if to delete first for kind %s", gvk.Kind)
	switch gvk.GroupKind().String() {
	case "Service", "Deployment.apps":
		// Try replace instead
		object, err := CreateOrReplace(ctx, cc, object)
		if err != nil {
			return object, err
		}
	default:
		object, err := CreateOrUpdate(ctx, cc, object)
		if err != nil {
			return object, err
		}
	}
	return object, nil
}

// CreateOrReplace works for services and deployments...
// - until we can address https://github.com/appvia/kore/issues/78
func CreateOrReplace(ctx context.Context, cc client.Client, object runtime.Object) (runtime.Object, error) {
	objMeta, _ := GetMeta(object)
	log.Debugf("deleting %s/%s as part of replace", objMeta.Namespace, objMeta.Name)
	// check type - can't go deleting CRD's / namespaces etc...
	if err := DeleteIfExists(ctx, cc, object); err != nil {
		return nil, fmt.Errorf(
			"error on delete as part of replace operation for item %s/%s - %s",
			objMeta.Namespace,
			objMeta.Name,
			err,
		)
	}
	object, err := CreateOrUpdate(ctx, cc, object)
	if err != nil {
		return nil, fmt.Errorf(
			"error on create as part of replace operation for item %s/%s - %s",
			objMeta.Namespace,
			objMeta.Name,
			err,
		)
	}
	return object, nil
}

// TypeSpecificUpdate will update where relevant logic is required
func TypeSpecificUpdate(ctx context.Context, cc client.Client, object runtime.Object) (bool, runtime.Object, error) {
	gvk := object.GetObjectKind().GroupVersionKind()
	switch gvk.GroupKind().String() {
	case "Service":
		err := DeleteIfExists(ctx, cc, object)
		if err != nil {
			return true, object, err
		}
		// use existing create logic
		return false, object, nil
	case "Deployment.apps", "ConfigMap":
		object, err := CreateOrForceUpdate(ctx, cc, object)
		if err != nil {
			return true, object, err
		}
		return true, object, nil
	default:
		return false, object, nil
	}
}

func CreateOrForceUpdate(ctx context.Context, cc client.Client, obj runtime.Object) (runtime.Object, error) {
	if err := cc.Create(ctx, obj); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}

		objMeta, err := meta.Accessor(obj)
		if err != nil {
			return nil, err
		}

		key := types.NamespacedName{
			Namespace: objMeta.GetNamespace(),
			Name:      objMeta.GetName(),
		}
		current := obj.DeepCopyObject()
		if err := cc.Get(ctx, key, current); err != nil {
			return nil, err
		}

		currentMeta, err := meta.Accessor(obj)
		if err != nil {
			return nil, err
		}

		objMeta.SetResourceVersion(currentMeta.GetResourceVersion())
		objMeta.SetGeneration(currentMeta.GetGeneration())

		return obj, cc.Update(ctx, obj)
	}

	return obj, nil
}

// UpdateIfChangedSinceLastUpdate updates the object only if the object data has changed since we've last applied it
// The last applied data will be saved in the `kore.appvia.io/last-applied` annotation and it will be used for comparison.
func UpdateIfChangedSinceLastUpdate(ctx context.Context, client client.Client, object, existing runtime.Object) (bool, error) {
	patchAnnotator := patch.NewAnnotator("kore.appvia.io/last-applied")
	var patchResult *patch.PatchResult

	if existing != nil {
		var err error
		patchResult, err = patch.NewPatchMaker(patchAnnotator).Calculate(
			existing,
			object,
			patch.IgnoreStatusFields(),
			// This is an ugly hack to compare the spec.configuration field as raw JSON strings
			// The strategic merge patch throws an error for some non-existing struct fields, in this case for the JSON struct:
			// > Failed to generate strategic merge patch: unable to find api field in struct JSON for the json field "resourceSelector"
			compareFieldAsRawJSON("spec.configuration"),
		)
		if err != nil {
			return false, err
		}
	}

	if patchResult == nil || !patchResult.IsEmpty() {
		if err := patchAnnotator.SetLastAppliedAnnotation(object); err != nil {
			return false, err
		}

		if _, err := CreateOrUpdate(ctx, client, object); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func compareFieldAsRawJSON(path string) patch.CalculateOption {
	return func(current, modified []byte) ([]byte, []byte, error) {
		var err error
		var compacted []byte

		val := gjson.GetBytes(current, path)
		if val.Exists() {
			if compacted, err = jsonutils.Compact([]byte(val.Raw)); err != nil {
				return nil, nil, err
			}
			if current, err = sjson.SetBytes(current, path, string(compacted)); err != nil {
				return nil, nil, err
			}
		}

		val = gjson.GetBytes(modified, path)
		if val.Exists() {
			if compacted, err = jsonutils.Compact([]byte(val.Raw)); err != nil {
				return nil, nil, err
			}
			if modified, err = sjson.SetBytes(modified, path, string(compacted)); err != nil {
				return nil, nil, err
			}
		}

		return current, modified, nil
	}
}
