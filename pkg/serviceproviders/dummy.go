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

package serviceproviders

import (
	"context"

	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	servicesv1 "github.com/appvia/kore/pkg/apis/services/v1"
	"github.com/appvia/kore/pkg/kore"

	"github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func init() {
	kore.DefaultServiceProviders.Register(Dummy{})
}

type Dummy struct {
}

func (d Dummy) Name() string {
	return "dummy"
}

func (d Dummy) Kinds() []string {
	return []string{"dummy"}
}

func (d Dummy) Plans() []servicesv1.ServicePlan {
	return []servicesv1.ServicePlan{
		{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServicePlan",
				APIVersion: servicesv1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "dummy",
				Namespace: "kore-admin",
			},
			Spec: servicesv1.ServicePlanSpec{
				Kind:          "dummy",
				Description:   "Used for testing",
				Summary:       "This is a dummy service plan",
				Configuration: v1beta1.JSON{Raw: []byte(`{"foo":"bar"}`)},
			},
		},
	}
}

func (d Dummy) JSONSchema(kind string) string {
	return `{
		"$id": "https://appvia.io/schemas/services/dummy/dummy.json",
		"$schema": "http://json-schema.org/draft-07/schema#",
		"description": "Dummy service plan schema",
		"type": "object",
		"additionalProperties": false,
		"required": [
			"foo"
		],
		"properties": {
			"foo": {
				"type": "string",
				"minLength": 1
			}
		}
	}`
}

func (d Dummy) Reconcile(_ context.Context, _ logrus.FieldLogger, service *servicesv1.Service) (reconcile.Result, error) {
	service.Status.Status = corev1.SuccessStatus
	return reconcile.Result{}, nil
}

func (d Dummy) Delete(_ context.Context, _ logrus.FieldLogger, service *servicesv1.Service) (reconcile.Result, error) {
	service.Status.Status = corev1.DeletedStatus
	return reconcile.Result{}, nil
}
