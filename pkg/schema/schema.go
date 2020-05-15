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

package schema

import (
	accountsv1beta1 "github.com/appvia/kore/pkg/apis/accounts/v1beta1"
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	eksv1alpha1 "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	gcpv1alpha1 "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	gkev1alpha1 "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	kschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	// hs is the schema for kore resources
	hs *runtime.Scheme
)

func init() {
	// @step: we start by registering the default core apigroups
	hs = scheme.Scheme

	builder := runtime.NewSchemeBuilder(
		apiextv1.AddToScheme,
		apiextv1beta1.AddToScheme,
		applicationv1beta.AddToScheme,
		accountsv1beta1.AddToScheme,
		clustersv1.AddToScheme,
		configv1.AddToScheme,
		corev1.AddToScheme,
		eksv1alpha1.AddToScheme,
		gcpv1alpha1.AddToScheme,
		gkev1alpha1.AddToScheme,
		orgv1.AddToScheme,
		servicesv1.AddToScheme,
	)
	if err := builder.AddToScheme(hs); err != nil {
		log.WithError(err).Fatal("failed to register the schema")
	}
}

// IsVersioned checks if the type is versioned in the scheme
func IsVersioned(object runtime.Object) bool {
	gvks, _, err := hs.ObjectKinds(object)
	if err != nil {
		return false
	}
	if len(gvks) >= 1 {
		return hs.Recognizes(gvks[0])
	}

	return false
}

// GetGroupKindVersion returns a schema for any registered type
func GetGroupKindVersion(object runtime.Object) (kschema.GroupVersionKind, bool, error) {
	possible, _, err := GetScheme().ObjectKinds(object)
	if err != nil {
		return kschema.GroupVersionKind{}, false, err
	}
	if len(possible) == 1 {
		return possible[0], IsVersioned(object), nil
	}

	return kschema.GroupVersionKind{}, IsVersioned(object), nil
}

// GetScheme returns a copy of the scheme
func GetScheme() *runtime.Scheme {
	return hs
}
