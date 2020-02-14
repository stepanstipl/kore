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

package clusterapp

import (
	"fmt"

	mutwebhookv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetClientOptions gets client options suitable for interactying with client apps
func GetClientOptions() (client.Options, error) {
	appSchema := runtime.NewScheme()
	err := addAllToScheme(appSchema)
	if err != nil {
		return client.Options{}, err
	}
	options := client.Options{
		Scheme: appSchema,
	}
	return options, nil
}

// AddAllToScheme allows us to interact with most types
func addAllToScheme(appScheme *runtime.Scheme) error {
	// We just have to add all the API's we intend to use here...

	// most bits
	err := v1.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding v1 schema - %s", err)
	}

	err = appsv1.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding apps v1 schema - %s", err)
	}

	err = rbacv1.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding batch schema - %s", err)
	}
	// Batch
	err = batchv1.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding v1role schema - %s", err)
	}

	// Supports CRD's etc...
	err = apiextv1.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding apiextv1 schema - %s", err)
	}
	err = apiextv1beta1.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding apiextv1beta1 schema - %s", err)
	}

	// Application kind
	err = applicationv1beta.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding crd schema for applications - %s", err)
	}

	// Mutating webhooks
	err = mutwebhookv1beta1.AddToScheme(appScheme)
	if err != nil {
		return fmt.Errorf("error adding mutatingwebhook v1beta - %s", err)
	}

	return nil
}
