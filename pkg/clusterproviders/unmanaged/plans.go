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

package unmanaged

import (
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var plans = []configv1.Plan{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kore",
			Annotations: map[string]string{
				"kore.appvia.io/system":   "true",
				"kore.appvia.io/readonly": "true",
			},
		},
		Spec: configv1.PlanSpec{
			Kind:        Kind,
			Summary:     "Default cluster plan for Kore",
			Description: "Default cluster plan for Kore",
			Configuration: apiextv1.JSON{
				Raw: []byte(`{}`),
			},
		},
	},
}
