/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
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

package schema

import (
	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gkev1alpha1 "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	kschema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

var (
	// hs is the schema for hub resources
	hs *runtime.Scheme
)

func init() {
	log.Info("registering the schema with kube-apiserver")

	// @step: we start by registering the default core apigroups
	hs = scheme.Scheme

	builder := runtime.NewSchemeBuilder(
		awsv1alpha1.AddToScheme,
		clustersv1.AddToScheme,
		configv1.AddToScheme,
		corev1.AddToScheme,
		gkev1alpha1.AddToScheme,
		orgv1.AddToScheme,
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
