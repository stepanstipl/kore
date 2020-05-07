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

package cluster

import (
	"context"
	"errors"
	"fmt"

	accounts "github.com/appvia/kore/pkg/apis/accounts/v1beta1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SetClusterRevision set the revision annotation on the resource
func SetClusterRevision(object runtime.Object, revision string) {
	kubernetes.SetRuntimeAnnotation(object, ClusterRevisionName, revision)
}

// GetClusterRevision retrieve the revision from the resource
func GetClusterRevision(object runtime.Object) string {
	return kubernetes.GetRuntimeAnnotation(object, ClusterRevisionName)
}

// GetObjectStatus attempts to inspect the resource for a status
func GetObjectStatus(object runtime.Object) (corev1.Status, error) {
	var status corev1.Status

	return status, kubernetes.GetRuntimeField(object, "status.status", &status)
}

// GetObjectReasonForFailure try's to get the reason for failure
func GetObjectReasonForFailure(object runtime.Object) (corev1.Condition, error) {
	c, err := GetObjectComponents(object)
	if err == nil {
		return c, nil
	}

	return GetObjectConditions(object)
}

// GetObjectComponents attempts to inspect the resource components
func GetObjectComponents(object runtime.Object) (corev1.Condition, error) {
	var components corev1.Components
	if err := kubernetes.GetRuntimeField(object, "status.components", &components); err != nil {
		return corev1.Condition{}, err
	}

	if components != nil && len(components) > 0 {
		return corev1.Condition{
			Detail:  components[0].Detail,
			Message: components[0].Message,
		}, nil
	}

	return corev1.Condition{}, kubernetes.ErrFieldNotFound
}

// GetObjectConditions returns the conditions on a resource
func GetObjectConditions(object runtime.Object) (corev1.Condition, error) {
	var conditions []corev1.Condition

	if err := kubernetes.GetRuntimeField(object, "status.conditions", &conditions); err != nil {
		return corev1.Condition{}, err
	}

	return conditions[0], nil
}

// IsDeleting check if the resource is being deleted
func IsDeleting(object runtime.Object) bool {
	mo, _ := object.(metav1.Object)

	return !mo.GetDeletionTimestamp().IsZero()
}

// SetRuntimeNamespace is used to apply the namespace
func SetRuntimeNamespace(object runtime.Object, namespace string) {
	mo, ok := object.(metav1.Object)
	if ok {
		mo.SetNamespace(namespace)
	}
}

// FindAccountManagement returns the account management
func FindAccountManagement(ctx context.Context, cc client.Client, owner corev1.Ownership) (*accounts.AccountManagement, error) {
	account := &accounts.AccountManagement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      owner.Name,
			Namespace: owner.Namespace,
		},
	}
	found, err := kubernetes.GetIfExists(ctx, cc, account)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("accounting resource %q does not exist", owner.Name)
	}

	return account, nil
}

// FindAccountingRule matches the account rule
func FindAccountingRule(account *accounts.AccountManagement, plan string) (*accounts.AccountsRule, bool) {
	for _, x := range account.Spec.Rules {
		if utils.Contains(plan, x.Plans) {
			return x, true
		}
	}

	return nil, false
}

// ComponentToUnstructured converts the component to a runtime reference
func ComponentToUnstructured(component *corev1.Component) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   component.Resource.Group,
		Version: component.Resource.Version,
		Kind:    component.Resource.Kind,
	})
	u.SetNamespace(component.Resource.Namespace)
	u.SetName(component.Resource.Name)

	return u
}

// IsComponentReferenced checks if the component is required
func IsComponentReferenced(component *corev1.Component, components *Components) (bool, error) {
	list, err := components.Walk()
	if err != nil {
		return false, err
	}

	for _, x := range list {
		resource := component.Resource
		if resource == nil {
			return false, controllers.NewCriticalError(errors.New("resource is nil"))
		}
		if kore.IsOwner(x.Object, *resource) {
			return true, nil
		}
	}

	return false, nil
}
