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

package utils

import (
	accountsv1beta1 "github.com/appvia/kore/pkg/apis/accounts/v1beta1"
	aksv1alpha1 "github.com/appvia/kore/pkg/apis/aks/v1alpha1"
	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	monitoringv1beta1 "github.com/appvia/kore/pkg/apis/monitoring/v1beta1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"
)

var (
	// ResourceList is a list of supported resources
	ResourceList = []Resource{
		{
			Name:         "accountmanagement",
			GroupVersion: accountsv1beta1.GroupVersion.String(),
			Kind:         "Account",
			Scope:        GlobalScope,
			ShortName:    "acc",
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Provider", "spec.provider", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "alert",
			APIName:      "monitoring/alerts",
			GroupVersion: monitoringv1beta1.SchemeGroupVersion.String(),
			Kind:         "Alert",
			Scope:        DualScope,
			ShortName:    "alerts",
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Resource", "status.rule.spec.resource.name", ""},
				{"Kind", "status.rule.spec.resource.kind", ""},
				{"Detail", "spec.summary", ""},
			},
		},
		{
			Name:         "aks",
			APIName:      "aks",
			GroupVersion: aksv1alpha1.GroupVersion.String(),
			Kind:         "AKS",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Credentials", "spec.credentials.name", ""},
				{"Cluster", "spec.cluster.name", ""},
				{"Version", "spec.version", ""},
				{"Region", "spec.region", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "akscredential",
			GroupVersion: aksv1alpha1.GroupVersion.String(),
			Kind:         "AKSCredentials",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Status", "status.status", ""},
				{"Verified", "status.verified", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "allocation",
			GroupVersion: configv1.GroupVersion.String(),
			Kind:         "Allocation",
			Scope:        TeamScope,
			ShortName:    "allo",
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Description", "spec.summary", ""},
				{"Owned By", "metadata.namespace", ""},
				{"Resource", "spec.resource.kind", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "audit",
			GroupVersion: orgv1.GroupVersion.String(),
			Kind:         "AuditEvent",
			Scope:        DualScope,
			Printer: []Column{
				{"Time", "metadata.creationTimestamp", ""},
				{"Operation", "spec.operation", ""},
				{"URI", "spec.resourceURI", ""},
				{"User", "spec.user", ""},
				{"Team", "spec.team", ""},
				{"Result", "spec.responseCode", ""},
			},
		},
		{
			Name:         "cluster",
			GroupVersion: clustersv1.GroupVersion.String(),
			Kind:         "Cluster",
			Scope:        TeamScope,
			ShortName:    "cs",
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Provider", "spec.kind", ""},
				{"Plan", "spec.plan", ""},
				{"Endpoint", "status.authProxyEndpoint", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "eks",
			APIName:      "eks",
			GroupVersion: eks.GroupVersion.String(),
			Kind:         "EKS",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Credentials", "spec.credentials.name", ""},
				{"Cluster", "spec.cluster.name", ""},
				{"Version", "spec.version", ""},
				{"Region", "spec.region", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "ekscredential",
			GroupVersion: eks.GroupVersion.String(),
			Kind:         "EKSCredentials",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Status", "status.status", ""},
				{"Verified", "status.verified", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "eksnodegroup",
			GroupVersion: eks.GroupVersion.String(),
			Kind:         "EKSNodeGroup",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Credentials", "spec.credentials.name", ""},
				{"Cluster", "spec.cluster.name", ""},
				{"Desired Size", "spec.desiredSize", ""},
				{"Instance Type", "spec.instanceType", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "eksvpc",
			GroupVersion: eks.GroupVersion.String(),
			Kind:         "EKSVPC",
			Scope:        TeamScope,
			ShortName:    "vpc",
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Credentials", "spec.credentials.name", ""},
				{"Cluster", "spec.cluster.name", ""},
				{"Network", "spec.privateIPV4Cidr", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "identity",
			APIName:      "identities",
			GroupVersion: orgv1.GroupVersion.String(),
			Kind:         "Identity",
			Scope:        GlobalScope,
			ShortName:    "ident",
			Printer: []Column{
				{"Name", "metadata.name", ""},
			},
		},
		{
			Name:         "gke",
			GroupVersion: gke.GroupVersion.String(),
			Kind:         "GKE",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Region", "spec.region", ""},
				{"Endpoint", "status.endpoint", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "gkecredential",
			GroupVersion: gke.GroupVersion.String(),
			Kind:         "GKECredentials",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Project", "spec.project", ""},
				{"Status", "status.status", ""},
				{"Verified", "status.verified", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "kubernetes",
			APIName:      "kubernetes",
			GroupVersion: clustersv1.GroupVersion.String(),
			Kind:         "Kubernetes",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Provider", "spec.provider.kind", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "member",
			GroupVersion: orgv1.GroupVersion.String(),
			Kind:         "Member",
			Scope:        TeamScope,
			ShortName:    "mb",
			Printer: []Column{
				{"Username", ".", ""},
			},
		},
		{
			Name:         "namespaceclaim",
			GroupVersion: clustersv1.GroupVersion.String(),
			Kind:         "NamespaceClaim",
			Scope:        TeamScope,
			Printer: []Column{
				{"Resource", "metadata.name", ""},
				{"Namespace", "spec.name", ""},
				{"Cluster", "spec.cluster.name", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "plan",
			GroupVersion: configv1.GroupVersion.String(),
			Kind:         "Plan",
			Scope:        GlobalScope,
			Printer: []Column{
				{"Resource", "metadata.name", ""},
				{"Summary", "spec.summary", ""},
				{"Description", "spec.description", ""},
				{"Kind", "spec.kind", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "planpolicy",
			GroupVersion: configv1.GroupVersion.String(),
			Kind:         "PlanPolicy",
			Scope:        GlobalScope,
			Printer: []Column{
				{"Resource", "metadata.name", ""},
				{"Summary", "spec.summary", ""},
				{"Description", "spec.description", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "organization",
			GroupVersion: gke.GroupVersion.String(),
			Kind:         "Organization",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "awsorganization",
			GroupVersion: aws.GroupVersion.String(),
			Kind:         "AWSOrganization",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "projectclaim",
			GroupVersion: gke.GroupVersion.String(),
			Kind:         "ProjectClaim",
			Scope:        TeamScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Organization", "spec.organization.name", ""},
				{"Owned By", "spec.organization.namespace", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "secret",
			GroupVersion: configv1.GroupVersion.String(),
			Kind:         "Secret",
			Scope:        TeamScope,
			ShortName:    "sc",
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Type", "spec.type", ""},
				{"Description", "spec.description", ""},
				{"Verified", "status.verified", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "team",
			GroupVersion: orgv1.GroupVersion.String(),
			Kind:         "Team",
			Scope:        GlobalScope,
			ShortName:    "tm",
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Description", "spec.description", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "user",
			GroupVersion: orgv1.GroupVersion.String(),
			Kind:         "User",
			Scope:        GlobalScope,
			Printer: []Column{
				{"Username", "metadata.name", ""},
				{"Email", "spec.email", ""},
				{"Disabled", "spec.disabled", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "service",
			GroupVersion: servicesv1.GroupVersion.String(),
			Kind:         "Service",
			Scope:        TeamScope,
			ShortName:    "svc",
			FeatureGate:  kore.FeatureGateServices,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Kind", "spec.kind", ""},
				{"Plan", "spec.plan", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "servicekind",
			GroupVersion: servicesv1.GroupVersion.String(),
			Kind:         "ServiceKind",
			Scope:        GlobalScope,
			ShortName:    "svcp",
			FeatureGate:  kore.FeatureGateServices,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Title", "spec.displayName", ""},
				{"Summary", "spec.summary", ""},
				{"Platform", "metadata.labels.kore\\.appvia\\.io/platform", ""},
				{"Enabled", "spec.enabled", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "serviceplan",
			GroupVersion: servicesv1.GroupVersion.String(),
			Kind:         "ServicePlan",
			Scope:        GlobalScope,
			ShortName:    "svcp",
			FeatureGate:  kore.FeatureGateServices,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Title", "spec.displayName", ""},
				{"Summary", "spec.summary", ""},
				{"Kind", "spec.kind", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "servicecredential",
			GroupVersion: servicesv1.GroupVersion.String(),
			Kind:         "ServiceCredentials",
			Scope:        TeamScope,
			FeatureGate:  kore.FeatureGateServices,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Status", "status.status", ""},
				{"Service", "spec.service.name", ""},
				{"Cluster", "spec.cluster.name", ""},
				{"Namespace", "spec.clusterNamespace", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "serviceprovider",
			GroupVersion: servicesv1.GroupVersion.String(),
			Kind:         "ServiceProvider",
			Scope:        GlobalScope,
			FeatureGate:  kore.FeatureGateServices,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Type", "spec.type", ""},
				{"Summary", "spec.summary", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
		{
			Name:         "config",
			GroupVersion: configv1.GroupVersion.String(),
			Kind:         "Config",
			Scope:        GlobalScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
			},
		},
		{
			Name:         "korefeature",
			GroupVersion: configv1.GroupVersion.String(),
			Kind:         "KoreFeature",
			Scope:        GlobalScope,
			Printer: []Column{
				{"Name", "metadata.name", ""},
				{"Type", "spec.featureType", ""},
				{"Status", "status.status", ""},
				{"Age", "metadata.creationTimestamp", "age"},
			},
		},
	}
)
