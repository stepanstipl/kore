/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
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

package hub

import (
	"context"
	"errors"
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// TeamsToList returns an array of teams
func TeamsToList(list *orgv1.TeamList) []string {
	items := make([]string, len(list.Items))

	for i := 0; i < len(list.Items); i++ {
		items[i] = list.Items[i].Name
	}

	return items
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

// IsProvider checks if the kubernetes cluster is back by the provider
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

// Label returns a hub label on a resource
func Label(tag string) string {
	return fmt.Sprintf("hub.appvia.io/%s", tag)
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
