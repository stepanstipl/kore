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
	"context"
	"regexp"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	monitoring "github.com/appvia/kore/pkg/apis/monitoring/v1beta1"
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
			Group:   orgv1.GroupVersion.Group,
			Version: orgv1.GroupVersion.Version,
			Kind:    `^(TeamMember|TeamInvitation|User|AuditEvent|Identity)$`,
		},
		{
			Group:   securityv1.GroupVersion.Group,
			Version: securityv1.GroupVersion.Version,
			Kind:    `.*`,
		},
		{
			Group:   monitoring.GroupVersion.Group,
			Version: monitoring.GroupVersion.Version,
			Kind:    `.*`,
		},
		{
			Group:   corev1.GroupVersion.Group,
			Version: corev1.GroupVersion.Version,
			Kind:    `^(IDP|IDPClient)$`,
		},
	}
)

// registerCustomResources is used to register the CRDs with the kubernetes api
func registerCustomResources(ctx context.Context, cc client.Interface) error {
	isFiltered := func(crd *apiextensions.CustomResourceDefinition, list []schema.GroupVersionKind) bool {
		for _, x := range list {
			if x.Group == crd.Spec.Group {
				if x.Version == crd.Spec.Version {
					re := regexp.MustCompile(x.Kind)
					if re.MatchString(crd.Spec.Names.Kind) {
						return true
					}
				}
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
		if err := crds.ApplyCustomResourceDefinition(ctx, cc, x); err != nil {
			return err
		}
	}

	return nil
}
