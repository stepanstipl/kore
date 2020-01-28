/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of hub-apiserver.
 *
 * hub-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * hub-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with hub-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */
// Code generated by register-gen. DO NOT EDIT.

package v1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName specifies the group name used to register the objects.
const GroupName = "clusters.compute.hub.appvia.io"

// GroupVersion specifies the group and the version used to register the objects.
var GroupVersion = v1.GroupVersion{Group: GroupName, Version: "v1"}

// SchemeGroupVersion is group version used to register these objects
// Deprecated: use GroupVersion instead.
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// localSchemeBuilder and AddToScheme will stay in k8s.io/kubernetes.
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	// Depreciated: use Install instead
	AddToScheme = localSchemeBuilder.AddToScheme
	Install     = localSchemeBuilder.AddToScheme
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Kubernetes{},
		&KubernetesList{},
		&ManagedClusterRole{},
		&ManagedClusterRoleBinding{},
		&ManagedClusterRoleBindingList{},
		&ManagedClusterRoleList{},
		&ManagedConfig{},
		&ManagedConfigList{},
		&ManagedPodSecurityPolicy{},
		&ManagedPodSecurityPolicyList{},
		&ManagedRole{},
		&ManagedRoleList{},
		&NamepacePolicy{},
		&NamepacePolicyList{},
		&NamespaceClaim{},
		&NamespaceClaimList{},
	)
	// AddToGroupVersion allows the serialization of client types like ListOptions.
	v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
