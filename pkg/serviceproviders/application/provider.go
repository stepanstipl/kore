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

package application

import (
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"

	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ kore.ServiceProvider = Provider{}

const (
	Type               = "application"
	ServiceKindApp     = "app"
	ServiceKindHelmApp = "helm-app"
)

type Provider struct {
	name  string
	plans []servicesv1.ServicePlan
}

func (p Provider) Name() string {
	return p.name
}

func (p Provider) Kinds() []servicesv1.ServiceKind {
	return []servicesv1.ServiceKind{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       servicesv1.ServiceKindGVK.Kind,
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      ServiceKindApp,
				Namespace: kore.HubNamespace,
			},
			Spec: servicesv1.ServiceKindSpec{
				DisplayName: "Kubernetes Application",
				Summary:     "Kubernetes Application",
				Enabled:     false,
			},
		},
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       servicesv1.ServiceKindGVK.Kind,
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      ServiceKindHelmApp,
				Namespace: kore.HubNamespace,
			},
			Spec: servicesv1.ServiceKindSpec{
				DisplayName: "Kubernetes Helm Application",
				Summary:     "Kubernetes Helm Application",
				Enabled:     false,
			},
		},
	}
}

func (p Provider) Plans() []servicesv1.ServicePlan {
	return p.plans
}

func (p Provider) AdminServices() []servicesv1.Service {
	cluster := corev1.Ownership{
		Group:     clustersv1.ClusterGroupVersionKind.Group,
		Version:   clustersv1.ClusterGroupVersionKind.Version,
		Kind:      clustersv1.ClusterGroupVersionKind.Kind,
		Namespace: "kore-admin",
		Name:      "kore",
	}

	var services []servicesv1.Service
	for _, servicePlan := range p.plans {
		if servicePlan.Annotations[kore.AnnotationSystem] != "true" {
			continue
		}

		// A small hack to generate the final plan name
		sp := servicePlan.DeepCopy()
		sp.Name = sp.Spec.Kind + "-" + sp.Name
		services = append(services, CreateSystemServiceFromPlan(*sp, cluster, servicePlan.Name, kore.HubAdminTeam))
	}
	return services
}

func (p Provider) PlanJSONSchema(kind string, _ string) (string, error) {
	switch kind {
	case ServiceKindApp:
		return AppSchema, nil
	case ServiceKindHelmApp:
		return HelmAppSchema, nil
	default:
		panic(fmt.Errorf("unexpected service kind: %s", kind))
	}
}

func (p Provider) CredentialsJSONSchema(_, _ string) (string, error) {
	return "", fmt.Errorf("can not create credentials for kubernetes services")
}

func (p Provider) RequiredCredentialTypes(_ string) ([]schema.GroupVersionKind, error) {
	return nil, nil
}

func (p Provider) ReconcileCredentials(
	ctx kore.ServiceProviderContext,
	service *servicesv1.Service,
	creds *servicesv1.ServiceCredentials,
) (reconcile.Result, map[string]string, error) {
	return reconcile.Result{}, nil, fmt.Errorf("can not create credentials for kubernetes services")
}

func (p Provider) DeleteCredentials(
	ctx kore.ServiceProviderContext,
	service *servicesv1.Service,
	creds *servicesv1.ServiceCredentials,
) (reconcile.Result, error) {
	return reconcile.Result{}, fmt.Errorf("can not create credentials for kubernetes services")
}
