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

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateManagedSecret returns a managed secret
func CreateManagedSecret(ctx context.Context, owner runtime.Object, cc client.Client, secret *configv1.Secret) error {
	if secret == nil {
		return errors.New("no secret defined")
	}

	// @step: create the object in the api
	if _, err := kubernetes.CreateOrUpdate(ctx, cc, secret); err != nil {
		return err
	}
	original := secret.DeepCopy()

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return utils.Retry(ctx, 5, true, 1000*time.Millisecond, func() (bool, error) {
		found, err := kubernetes.GetIfExists(ctx, cc, secret)
		if err != nil || !found {
			return false, nil
		}
		if secret.Status.SystemManaged != nil {
			return true, nil
		}
		secret.Status.SystemManaged = utils.TruePtr()

		if err := cc.Status().Patch(ctx, secret, client.MergeFrom(original)); err == nil {
			return false, nil
		}

		return true, nil
	})
}
