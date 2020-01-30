/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package kubernetes

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CheckIfExists checks if a runtime object exist in the api
func CheckIfExists(ctx context.Context, cc client.Client, object runtime.Object) (bool, error) {
	return GetIfExists(ctx, cc, object)
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

// CreateOrUpdate is responsible for updating a resouce
// extended to carry out updates on specific types when patching is not suitable
func CreateOrUpdate(ctx context.Context, cc client.Client, object runtime.Object) (runtime.Object, error) {
	supported, object, err := TypeSpecificUpdate(ctx, cc, object)
	if supported {
		return object, err
	}

	// default case
	key, err := client.ObjectKeyFromObject(object)
	if err != nil {
		return nil, err
	}

	// @step: we first try and create the resource
	if err := cc.Create(ctx, object); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return nil, err
		}
		// @step: we need to retrieve the current one
		original := object.DeepCopyObject()

		if err := cc.Get(ctx, key, original); err != nil {
			return nil, err
		}

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

	return object, nil
}

// PatchOrReplace is responsible for updating a resouce and optionally updating if patching fails
// Due to issues with some immutable objects set on the server, this will allow a fall back for now
// See similar related issue with patching in Helm: https://github.com/helm/helm/issues/7516
func PatchOrReplace(ctx context.Context, cc client.Client, object runtime.Object) (runtime.Object, error) {
	// need to create or replace for deployments and services
	gvk := object.GetObjectKind().GroupVersionKind()
	log.Debugf("deciding if to delete first for kind %s", gvk.Kind)
	switch gvk.Kind {
	case "Service", "Deployment":
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
	switch gvk.Kind {
	case "Service":
		err := DeleteIfExists(ctx, cc, object)
		if err != nil {
			return true, object, err
		}
		// use existing create logic
		return false, object, nil
	case "Deployment":
		deploy := object.(*appsv1.Deployment)
		object, err := CreateOrUpdateDeployment(ctx, cc, deploy)
		if err != nil {
			return true, object, err
		}
		return true, object, nil
	case "ConfigMap":
		cm := object.(*corev1.ConfigMap)
		object, err := CreateOrUpdateConfigMap(ctx, cc, cm)
		if err != nil {
			return true, object, err
		}
		return true, object, nil
	default:
		return false, object, nil
	}
}
