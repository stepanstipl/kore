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

package kore

import (
	"context"
	"errors"
	"fmt"

	"github.com/appvia/kore/pkg/utils/validation"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore/authentication"
	koreschema "github.com/appvia/kore/pkg/schema"
	"github.com/appvia/kore/pkg/utils"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IsValidResourceName checks the resource name is valid
// https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
func IsValidResourceName(subject, name string) error {
	if !ResourceNameFilter.MatchString(name) {
		return validation.NewError("%s has failed validation", subject).WithFieldErrorf(
			"name",
			validation.Pattern,
			"must comply with %s",
			ResourceNameFilter.String(),
		)
	}

	if len(name) > 63 {
		return validation.NewError("%s has failed validation", subject).WithFieldError(
			"name",
			validation.MaxLength,
			"length must be less than 64 characters",
		)
	}

	return nil
}

// TeamsToList returns an array of teams
func TeamsToList(list *orgv1.TeamList) []string {
	items := make([]string, len(list.Items))

	for i := 0; i < len(list.Items); i++ {
		items[i] = list.Items[i].Name
	}

	return items
}

// HasAccessToTeam checks if the user has access to the team
func HasAccessToTeam(ctx context.Context, team string) bool {
	user := authentication.MustGetIdentity(ctx)

	if user.IsGlobalAdmin() {
		return true
	}

	return utils.Contains(team, user.Teams())
}

// IsGlobalTeam checks if the namespace is global
func IsGlobalTeam(name string) bool {
	return name == HubAdminTeam
}

// IsOwn checks the ownership are the same
func IsOwn(a, b corev1.Ownership) bool {
	fields := map[string]string{
		a.Group:     b.Group,
		a.Version:   b.Version,
		a.Kind:      b.Kind,
		a.Namespace: b.Namespace,
		a.Name:      b.Name,
	}
	for k, v := range fields {
		if k != v {
			return false
		}
	}

	return true
}

// IsResourceOwner checks if the object is pointed to by the ownership reference
func IsResourceOwner(o runtime.Object, ownership corev1.Ownership) (bool, error) {
	if o == nil {
		return false, errors.New("no object defined")
	}
	mo, ok := o.(metav1.Object)
	if !ok {
		return false, errors.New("object does not implement metav1.Object")
	}

	gvk, found, err := koreschema.GetGroupKindVersion(o)
	if err != nil {
		return false, err
	}
	if !found {
		return false, errors.New("resource not found in registered schema")
	}

	switch {
	case mo.GetName() != ownership.Name:
		return false, nil
	case mo.GetNamespace() != ownership.Namespace:
		return false, nil
	case gvk.Group != ownership.Group:
		return false, nil
	case gvk.Version != ownership.Version:
		return false, nil
	case gvk.Kind != ownership.Kind:
		return false, nil
	}

	return true, nil
}

// ResourceExists checks if some resource exists
func ResourceExists(client client.Client, resource corev1.Ownership) (bool, error) {
	// @step: convert to an unstructured
	u, err := ToUnstructuredFromOwnership(resource)
	if err != nil {
		return false, err
	}

	// @step: check if the resource exists
	if err := client.Get(context.Background(), types.NamespacedName{
		Namespace: resource.Namespace,
		Name:      resource.Name,
	}, u); err != nil {
		if kerrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// ToUnstructuredFromOwnership converts an ownership to an unstructured type
func ToUnstructuredFromOwnership(resource corev1.Ownership) (*unstructured.Unstructured, error) {
	if err := IsOwnershipValid(resource); err != nil {
		return nil, err
	}

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   resource.Group,
		Version: resource.Version,
		Kind:    resource.Kind,
	})
	u.SetName(resource.Name)
	u.SetNamespace(resource.Namespace)

	return u, nil
}

// IsProviderBacked checks if the kubernetes cluster is back by the provider
func IsProviderBacked(cluster *clustersv1.Kubernetes) bool {
	return HasOwnership(cluster.Spec.Provider)
}

// HasOwnership checks if the ownership is set
func HasOwnership(owner corev1.Ownership) bool {
	// @step: if any of fields are set we assume use
	fields := []string{
		owner.Group,
		owner.Version,
		owner.Kind,
		owner.Namespace,
		owner.Name,
	}
	for _, x := range fields {
		if x != "" {
			return true
		}
	}

	return false
}

// IsOwner checks if the ownerships matches
func IsOwner(obj runtime.Object, ownership corev1.Ownership) bool {
	gvk := obj.GetObjectKind().GroupVersionKind()

	mo, ok := obj.(metav1.Object)
	if !ok {
		return false
	}

	switch {
	case gvk.Group != ownership.Group:
		return false
	case gvk.Version != ownership.Version:
		return false
	case gvk.Kind != ownership.Kind:
		return false
	case mo.GetNamespace() != ownership.Namespace:
		return false
	case mo.GetName() != ownership.Name:
		return false
	}

	return true
}

// IsOwnershipValid checks the ownership is filled in
func IsOwnershipValid(owner corev1.Ownership) error {
	fields := map[string]string{
		"group":     owner.Group,
		"version":   owner.Version,
		"kind":      owner.Kind,
		"namespace": owner.Namespace,
		"name":      owner.Name,
	}
	for k, v := range fields {
		if v == "" {
			return fmt.Errorf("%s field in ownership is not defined", k)
		}
	}

	return nil
}

// UnstructuredKind returns an unstructured kind
func UnstructuredKind(gvk schema.GroupVersionKind) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(gvk)

	return u
}

// IsValidGVK checks if the GVK is valid
func IsValidGVK(gvk schema.GroupVersionKind) error {
	if gvk.Group == "" {
		return errors.New("missing apigroup")
	}
	if gvk.Version == "" {
		return errors.New("missing apigroup version")
	}
	if gvk.Kind == "" {
		return errors.New("missing apigroup kind")
	}

	return nil
}

// Label returns a kore label on a resource
func Label(tag string) string {
	return fmt.Sprintf("kore.appvia.io/%s", tag)
}

// EmptyUser returns an empty user
func EmptyUser(username string) *orgv1.User {
	return &orgv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name:      username,
			Namespace: HubNamespace,
		},
		Spec: orgv1.UserSpec{Username: username},
	}
}
