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

package server

import (
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/register"
	"github.com/appvia/kore/pkg/utils/crds"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	filtered = []schema.GroupVersionKind{
		{
			Group:   clustersv1.GroupVersion.Group,
			Version: clustersv1.GroupVersion.Version,
			Kind:    "KubernetesCredentials",
		},
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
		{
			Group:   orgv1.GroupVersion.Group,
			Version: orgv1.GroupVersion.Version,
			Kind:    "AuditEvent",
		},
		{
			Group:   securityv1.GroupVersion.Group,
			Version: securityv1.GroupVersion.Version,
			Kind:    "SecurityRule",
		},
		{
			Group:   securityv1.GroupVersion.Group,
			Version: securityv1.GroupVersion.Version,
			Kind:    "SecurityRuleList",
		},
		{
			Group:   securityv1.GroupVersion.Group,
			Version: securityv1.GroupVersion.Version,
			Kind:    "SecurityOverview",
		},
		{
			Group:   securityv1.GroupVersion.Group,
			Version: securityv1.GroupVersion.Version,
			Kind:    "SecurityScanResult",
		},
		{
			Group:   securityv1.GroupVersion.Group,
			Version: securityv1.GroupVersion.Version,
			Kind:    "SecurityScanResultList",
		},
	}
)

// registerCustomResources is used to register the CRDs with the kubernetes api
func registerCustomResources(cc client.Interface) error {
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
