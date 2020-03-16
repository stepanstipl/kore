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

package adminproject

import (
	"context"
	"errors"

	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
)

// EnsureProject is responsible for ensuring the project exists or is created from the oauth token
func (t *gcpCtrl) EnsureProject(
	ctx context.Context,
	project *gcp.GCPAdminProject) error {

	logger := log.WithFields(log.Fields{
		"name":      project.Name,
		"namespace": project.Namespace,
	})
	logger.Debug("ensuring the gcp admin project exists")

	// @logic
	// - if we have gcp credentials we take those and check for the project (simple)
	// - if we don't have gcp credentials we check for oauth token
	// - if we don't have a token we fail
	// - if we do have a token we use to provision a project and create the service account credentials

	stage := "provision"
	secret := CreateSecretRef(project.Spec.CredentialsRef.Namespace, project.Spec.CredentialsRef.Name)

	// @step: do we have credentials for the project?
	_, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), secret)
	if err != nil {
		logger.WithError(err).Error("trying to check for gcp credentials")

		project.Status.Conditions.SetCondition(corev1.Component{
			Detail:  err.Error(),
			Name:    stage,
			Message: "Checking for the GCP credentials secret",
			Status:  corev1.FailureStatus,
		})

		return err
	}

	//

	// @step: check the project exists

	return nil
}

// EnsureCredentials is responsible for checking, creating the credentials
func (t *gcpCtrl) EnsureCredentials(ctx context.Context, project *gcp.GCPAdminProject) (*configv1.Secret, error) {
	logger := log.WithFields(log.Fields{
		"name":      project.Name,
		"namespace": project.Namespace,
	})

	stage := "provision"
	secret := CreateSecretRef(project.Spec.CredentialsRef.Namespace, project.Spec.CredentialsRef.Name)

	// @step: do we have credentials for the project?
	found, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), secret)
	if err != nil {
		logger.WithError(err).Error("trying to check for gcp credentials")

		project.Status.Conditions.SetCondition(corev1.Component{
			Detail:  err.Error(),
			Name:    stage,
			Message: "Checking for the GCP credentials secret",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	// @step: if the credentials secret does not exist, what about the oauth token?
	if !found {
		logger.Debug("gcp credentials do not exist, checking if we have oauth token")

		// @point we don't have a credentials secret or a token reference
		if project.Spec.TokenRef == nil || (project.Spec.TokenRef.Namespace == "" || project.Spec.TokenRef.Name == "") {

			project.Status.Conditions.SetCondition(corev1.Component{
				Detail:  "no credentials supplied",
				Name:    stage,
				Message: "You haven't specified either GCP credentials secret or a oauth token",
				Status:  corev1.FailureStatus,
			})

			return nil, errors.New("no credentials supplied")
		}

		// @step: we have a reference to a oauth token, lets retrieve it and check
		secret := CreateSecretRef(project.Spec.TokenRef.Namespace, project.Spec.TokenRef.Name)

		found, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), secret)
		if err != nil {
			logger.WithError(err).Error("trying to check for the gcp oauth token")

			project.Status.Conditions.SetCondition(corev1.Component{
				Detail:  err.Error(),
				Name:    stage,
				Message: "Unable to check for oauth token credentials secret",
				Status:  corev1.FailureStatus,
			})

			return nil, err
		}
		if !found {
			project.Status.Conditions.SetCondition(corev1.Component{
				Detail:  "no credentials supplied",
				Name:    stage,
				Message: "The reference to the oauth token is invalid, no credentials secret found",
				Status:  corev1.FailureStatus,
			})

			return nil, errors.New("no credentials supplied")
		}

		// else we have the oauth token to create the admin project
		logger.Debug("attempting to create the gcp admin project from oauth token")
	}

	// @step: we need to check the credentials are valid

	return secret, nil
}

// EnsureCredentialPermissions is responsible for checking the permissions are correct
func (t gcpCtrl) EnsureCredentialPermissions(
	ctx context.Context,
	project *gcp.GCPAdminProject,
	credentials *configv1.Secret) error {

	logger := log.WithFields(log.Fields{
		"name":      project.Name,
		"namespace": project.Namespace,
	})
	logger.Debug("checking the credentials for the admin project are correct")

	//stage := "permissions"

	return nil
}
