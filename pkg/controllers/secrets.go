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

package controllers

import (
	"context"
	"errors"
	"time"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CreateManagedSecret returns a managed secret
func CreateManagedSecret(ctx context.Context, owner runtime.Object, cc client.Client, secret *configv1.Secret) error {
	if secret == nil {
		return errors.New("no secret defined")
	}

	mo, ok := owner.(metav1.Object)
	if !ok {
		return errors.New("object does not implement the metav1.Object")
	}

	secret.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: owner.GetObjectKind().GroupVersionKind().GroupVersion().String(),
			Kind:       owner.GetObjectKind().GroupVersionKind().Kind,
			Controller: utils.TruePtr(),
			UID:        mo.GetUID(),
			Name:       mo.GetName(),
		},
	}

	// @step: update the secret
	if _, err := kubernetes.CreateOrUpdate(ctx, cc, secret); err != nil {
		return err
	}

	// @step: update the secret to be managed
	for i := 0; i < 3; i++ {
		if found, err := kubernetes.GetIfExists(ctx, cc, secret); err != nil {
			continue
		} else if !found {
			return nil
		}
		original := secret.DeepCopy()

		if secret.Status.SystemManaged != nil {
			return nil
		}
		secret.Status.SystemManaged = utils.TruePtr()

		if err := cc.Status().Patch(ctx, secret, client.MergeFrom(original)); err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return errors.New("trying to update the system managed status")
}
