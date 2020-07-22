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

package security

import (
	"context"
	"testing"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestNewEmbeddedRules(t *testing.T) {
	e, err := NewEmbeddedRules()
	require.NoError(t, err)
	require.NotNil(t, e)
}

func TestNewEmbeddedRulesList(t *testing.T) {
	e, err := NewEmbeddedRules()
	require.NoError(t, err)
	require.NotNil(t, e)
	assert.NotEmpty(t, e.List())
}

func TestEmbeddedRules(t *testing.T) {
	plan := &configv1.Plan{
		TypeMeta: metav1.TypeMeta{
			Kind: "Plan",
		},
		Spec: configv1.PlanSpec{
			Kind: "GKE",
			Configuration: apiextv1.JSON{Raw: []byte(`
				{
					"enableSometning": true
				}
			`)},
		},
	}

	e, err := NewEmbeddedRules()
	require.NoError(t, err)
	require.NotNil(t, e)

	for _, x := range e.List() {
		x.(PlanRule).CheckPlan(context.TODO(), fake.NewFakeClient(), plan)
	}
}
