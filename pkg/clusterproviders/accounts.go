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

package clusterproviders

import (
	"fmt"

	accountsv1beta1 "github.com/appvia/kore/pkg/apis/accounts/v1beta1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IsAccountManaged returns if we are managing accounts for this cluster
func IsAccountManaged(owner corev1.Ownership) bool {
	if owner.Group != accountsv1beta1.GroupVersion.Group {
		return false
	}
	if owner.Version != accountsv1beta1.GroupVersion.Version {
		return false
	}
	if owner.Kind != "AccountManagement" {
		return false
	}

	return true
}

// FindAccountManagement gets an accounting object from ownership
func FindAccountManagement(ctx kore.Context, owner corev1.Ownership) (*accountsv1beta1.AccountManagement, error) {
	account := &accountsv1beta1.AccountManagement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      owner.Name,
			Namespace: owner.Namespace,
		},
	}
	found, err := kubernetes.GetIfExists(ctx, ctx.Client(), account)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("accounting resource %q does not exist", owner.Name)
	}

	return account, nil
}

// FindAccountingRule will discover an account rule from a plan name
func FindAccountingRule(account *accountsv1beta1.AccountManagement, plan string) (*accountsv1beta1.AccountsRule, bool) {
	for _, x := range account.Spec.Rules {
		if utils.Contains(plan, x.Plans) {
			return x, true
		}
	}

	return nil, false
}
