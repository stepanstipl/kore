/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package server

import (
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/register"
	"github.com/appvia/kore/pkg/utils/crds"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// registerCustomResources is used to register the CRDs with the kubernetes api
func registerCustomResources(cc client.Interface) error {
	filtered := []schema.GroupVersionKind{
		{
			Group:   orgv1.GroupVersion.Group,
			Version: orgv1.GroupVersion.Version,
			Kind:    "TeamMember",
		},
		{
			Group:   orgv1.GroupVersion.Group,
			Version: orgv1.GroupVersion.Version,
			Kind:    "TeamInvitation",
		},
		{
			Group:   orgv1.GroupVersion.Group,
			Version: orgv1.GroupVersion.Version,
			Kind:    "User",
		},
	}

	isFiltered := func(crd *apiextensions.CustomResourceDefinition, list []schema.GroupVersionKind) bool {
		for _, x := range list {
			if x.Group == crd.Spec.Group && x.Version == crd.Spec.Version && x.Kind == crd.Spec.Names.Kind {
				return true
			}
		}

		return false
	}

	list, err := register.GetCustomResourceDefinitions()
	if err != nil {
		return err
	}

	for _, x := range list {
		if isFiltered(x, filtered) {
			continue
		}
		if err := crds.ApplyCustomResourceDefinition(cc, x); err != nil {
			return err
		}
	}

	return nil
}
