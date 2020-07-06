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

package awsaccount

import (
	"context"
	"errors"
	"fmt"
	"strings"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/cloud/aws"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ComponentOrganizationVerify is the text for the stage of AWS Organization verification
	ComponentOrganizationVerify = "AWS Organization Verification"
	// ComponentAccountCreation describes when the account is being created
	ComponentAccountCreation = "AWS Account Creation"
	// ComponentAccountMasterAccess describes the stage when stacksets deploy the initial admin role
	ComponentAccountMasterAccess = "AWS Master Account Access"
	// ComponentAccountTeamAccess is the stage when a team account is created for provisioning clusters
	ComponentAccountTeamAccess = "AWS Team Account Access"
)

func (t *awsCtrl) EnsurePermitted(ctx context.Context, account *awsv1alpha1.AWSAccount) error {
	// @step: we check if the aws organization has been allocated to us
	permitted, err := t.Teams().Team(account.Namespace).Allocations().IsPermitted(ctx, account.Spec.Organization)
	if err != nil {
		return err
	}
	if !permitted {
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentOrganizationVerify,
			Message: "AWS Organization has not been allocated to you",
			Status:  corev1.FailureStatus,
		})

		return errors.New("AWS organization has not been allocated to team")
	}

	return nil
}

// EnsureAccountUnclaimed is responsible for making sure the aws account is unclaimed
func (t *awsCtrl) EnsureAccountUnclaimed(ctx context.Context, account *awsv1alpha1.AWSAccount) error {
	logger := log.WithFields(log.Fields{
		"name":      account.Name,
		"namespace": account.Namespace,
	})

	// @step: check if the projected account has already been claimed else where
	alreadyUsed, err := t.IsAccountAlreadyUsed(ctx, account)
	if err != nil {
		logger.WithError(err).Error("trying to check if the account is already used")

		account.Status.Status = corev1.FailureStatus
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentOrganizationVerify,
			Message: "Unable to fulfil request, cannot check if account name has already been used in the organization",
			Status:  corev1.FailureStatus,
		})

		return errors.New("failed to check if aws account is already in use")
	}
	if alreadyUsed {
		logger.Warn("attempting to use an aws account which has already been provisioned")

		account.Status.Status = corev1.FailureStatus
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentOrganizationVerify,
			Message: "AWS account has already been used by another team in kore",
			Status:  corev1.FailureStatus,
		})

		return errors.New("aws account name already provisioned")
	}

	return nil
}

// EnsureOrganization is responsible for checking and retrieving the aws org
func (t *awsCtrl) EnsureOrganization(ctx context.Context, account *awsv1alpha1.AWSAccount) (*awsv1alpha1.AWSOrganization, error) {
	org := &awsv1alpha1.AWSOrganization{}

	key := types.NamespacedName{
		Namespace: account.Spec.Organization.Namespace,
		Name:      account.Spec.Organization.Name,
	}

	if err := t.mgr.GetClient().Get(ctx, key, org); err != nil {
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentOrganizationVerify,
			Detail:  err.Error(),
			Message: "Attempting to retrieve the AWS Organization resources from API",
			Status:  corev1.FailureStatus,
		})

		return nil, err
	}

	// @step: check if the aws organization exists and if successful (Validated OU)
	if org.Status.Status != corev1.SuccessStatus {
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentOrganizationVerify,
			Detail:  "resource is in failing state",
			Message: "AWS Organisation is in a failing state, cannot provision account",
			Status:  corev1.FailureStatus,
		})

		return nil, errors.New("aws organisation cannot be verified or failed")
	}

	return org, nil
}

// EnsureAWSOrganization is responsible for ensuring the AWS Organisation is valid
func (t *awsCtrl) EnsureAWSAccountProvisioned(
	ctx context.Context,
	account *awsv1alpha1.AWSAccount,
	org *awsv1alpha1.AWSOrganization,
	creds *aws.Credentials) (*aws.AccountClient, bool, error) {

	logger := log.WithFields(log.Fields{
		"name":      account.Name,
		"namespace": account.Namespace,
	})

	// @step: check the accounct exists
	client := aws.NewAccountClientFromCredsAndRole(*creds, org.Spec.RoleARN, org.Spec.Region, aws.Account{
		AccountEmail:              t.DeriveAccountEmail(org, account),
		ManagedOrganizationalUnit: org.Spec.OuName,
		NewAccountName:            account.Name,
		PrimaryResourceRegion:     account.Spec.Region,
		SSOUserEmail:              org.Spec.SsoUser.Email,
		SSOUserFirstName:          org.Spec.SsoUser.FirstName,
		SSOUserLastName:           org.Spec.SsoUser.LastName,
	})
	successComponent := corev1.Component{
		Name:    ComponentAccountCreation,
		Detail:  "Account Created",
		Message: "Account has been provisioned",
		Status:  corev1.SuccessStatus,
	}

	exists, err := client.Exists()
	if err != nil {
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentAccountCreation,
			Detail:  err.Error(),
			Message: "Attempting to retrieve the account from aws",
			Status:  corev1.FailureStatus,
		})

		return client, false, err
	}
	if !exists {
		// Create the account (if not already being provisioned)
		if account.Status.ServiceCatalogProvisioningID != "" {
			// is service catalog provisoning still in progress?
			ready, err := client.IsAccountReady(account.Status.ServiceCatalogProvisioningID)
			if err != nil {

				logger.Warnf("unable to check for provisioning record from service catalog %s - %s", account.Status.ServiceCatalogProvisioningID, err)

				account.Status.Conditions.SetCondition(corev1.Component{
					Name:    ComponentAccountCreation,
					Detail:  err.Error(),
					Message: "Atempting to verify account factory service catalogue record",
					Status:  corev1.FailureStatus,
				})

				return nil, false, err
			}
			if !ready {

				account.Status.Conditions.SetCondition(corev1.Component{
					Name:    ComponentAccountCreation,
					Detail:  "Account creation in progress",
					Message: "Account has been requested through the control tower service catalog",
					Status:  corev1.PendingStatus,
				})
				// Back off and wait (provisioning)
				return client, false, nil
			}
			// account is provisioned so stop tracking the service catalog ID
			account.Status.ServiceCatalogProvisioningID = ""

			// account creation complete, next step
			account.Status.Conditions.SetCondition(successComponent)
			return client, true, nil
		}
		// We need to provision now...
		account.Status.ServiceCatalogProvisioningID, err = client.CreateNewAccount()
		if err != nil {
			logger.Warnf("unable to create account %s - %s", account.Spec.AccountName, err)

			account.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentAccountCreation,
				Detail:  err.Error(),
				Message: "Problem when creating account",
				Status:  corev1.FailureStatus,
			})

			return client, false, err
		}
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentAccountCreation,
			Detail:  "Account creation started",
			Message: "Account requested",
			Status:  corev1.PendingStatus,
		})

		// We now need to back off and wait
		return client, false, nil
	}
	account.Status.Conditions.SetCondition(successComponent)

	// Account already exists so provisioning complete
	return client, true, nil
}

// EnsureAccessFromMasterRole will create any missing stacksets to gain access to an account
func (t *awsCtrl) EnsureAccessFromMasterRole(
	ctx context.Context,
	account *awsv1alpha1.AWSAccount,
	client aws.AccountClienter) (bool, error) {

	logger := log.WithFields(log.Fields{
		"name":      account.Name,
		"namespace": account.Namespace,
	})
	err := client.EnsureInitialAccessCreated()

	if err != nil {
		logger.Warnf("unable to ensure account access provisioned for %s - %s", account.Spec.AccountName, err)

		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentAccountMasterAccess,
			Detail:  err.Error(),
			Message: "Problem when deploying access to account",
			Status:  corev1.FailureStatus,
		})

		return false, err
	}

	ready, err := client.IsInitialAccessReady()
	if err != nil {
		logger.Warnf("unable to check for account access for %s - %s", account.Spec.AccountName, err)

		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentAccountMasterAccess,
			Detail:  err.Error(),
			Message: "Problem when checking access for account",
			Status:  corev1.FailureStatus,
		})

		return false, err
	}
	if !ready {
		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentAccountMasterAccess,
			Detail:  "Waiting for stacksets to provision initial access",
			Message: "Initial access from master account is being provisioned",
			Status:  corev1.PendingStatus,
		})

		return false, nil
	}
	account.Status.Conditions.SetCondition(corev1.Component{
		Name:    ComponentAccountMasterAccess,
		Detail:  "Account access is granted from master account",
		Message: "Access has been verified",
		Status:  corev1.SuccessStatus,
	})

	return true, nil
}

// EnsureCredentials is responsible for checking, creating the credentials
func (t *awsCtrl) EnsureCredentials(ctx context.Context, account *awsv1alpha1.AWSAccount, client aws.AccountClienter) error {
	logger := log.WithFields(log.Fields{
		"name":      account.Name,
		"namespace": account.Namespace,
	})
	secretName := t.GetAccountCredentialsSecretName(account)
	secret := &configv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: account.Namespace,
		},
	}
	if account.Status.CredentialRef == nil {
		account.Status.CredentialRef = &v1.SecretReference{
			Name:      secretName,
			Namespace: account.Namespace,
		}
	}

	found, err := kubernetes.GetIfExists(ctx, t.mgr.GetClient(), secret)
	if err != nil {
		logger.Warnf("unable to lookup account credentials secret %s/%s- %s", account.Namespace, account.Status.AccountID, err)

		account.Status.Conditions.SetCondition(corev1.Component{
			Name:    ComponentAccountTeamAccess,
			Detail:  err.Error(),
			Message: "Attempting to retrieve the AWS Account credentials",
			Status:  corev1.FailureStatus,
		})

		return err
	}
	if !found {
		// TODO: check if account user and then credentials need to be re-created for a user we deleted just from kore
		// 		 see - https://github.com/appvia/kore/issues/1051

		// Create the account credentials
		creds, err := client.CreateAccountCredentials()
		if err != nil {
			logger.Warnf("unable to create aws account credentials for %s - %s", account.Status.AccountID, err)

			account.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentAccountTeamAccess,
				Detail:  err.Error(),
				Message: "AWS Account credentials could not be created in aws",
				Status:  corev1.FailureStatus,
			})

			return errors.New("credentials not created")
		}
		// Persist the account credentials for project access
		secret = CreateAWSCredentialsSecret(account, secretName, map[string]string{
			"access_key_id":     creds.AccessKeyID,
			"access_secret_key": creds.SecretAccessKey,
		})

		if _, err := kubernetes.CreateOrUpdate(ctx, t.mgr.GetClient(), secret.Encode()); err != nil {
			logger.WithError(err).Error("trying to update or create the credentials")

			account.Status.Conditions.SetCondition(corev1.Component{
				Name:    ComponentAccountTeamAccess,
				Detail:  err.Error(),
				Message: "AWS Account credentials could not be persisted",
				Status:  corev1.FailureStatus,
			})

			return errors.New("credentials not persisted")
		}
	}
	account.Status.Conditions.SetCondition(corev1.Component{
		Name:    ComponentAccountTeamAccess,
		Detail:  "AWS account team access granted",
		Message: "AWS account credentials persisted",
		Status:  corev1.SuccessStatus,
	})

	return nil
}

// EnsureCredentialsAllocation is responsible for creating an allocation to the credentials
func (t *awsCtrl) EnsureCredentialsAllocation(
	ctx context.Context,
	account *awsv1alpha1.AWSAccount) error {

	logger := log.WithFields(log.Fields{
		"name":      account.Name,
		"namespace": account.Namespace,
	})
	logger.Debug("attempting to create the allocation for the aws account")

	name := t.GetAllocationName(account)

	allocation := &configv1.Allocation{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.GroupVersion.String(),
			Kind:       "Allocation",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: account.Namespace,
		},
		Spec: configv1.AllocationSpec{
			Name:    "aws-" + account.Name,
			Summary: "Provides credentials to team AWS Account " + account.Spec.AccountName,
			Resource: corev1.Ownership{
				Group:     awsv1alpha1.GroupVersion.Group,
				Kind:      "AWSAccount",
				Name:      account.Name,
				Namespace: account.Namespace,
				Version:   awsv1alpha1.GroupVersion.Version,
			},
			Teams: []string{account.Namespace},
		},
	}

	// @step: check if the allocation exists
	found, err := kubernetes.CheckIfExists(ctx, t.mgr.GetClient(), allocation)
	if err != nil {
		logger.WithError(err).Error("trying to check for allocation")

		return err
	}
	if found {
		return nil
	}

	if _, err := kubernetes.CreateOrUpdate(ctx, t.mgr.GetClient(), allocation); err != nil {
		logger.WithError(err).Error("trying to create the aws account to team allocation")

		return err
	}

	return nil
}

// IsAccountAlreadyUsed checks if the aws account name has already been used by another team
func (t *awsCtrl) IsAccountAlreadyUsed(ctx context.Context, account *awsv1alpha1.AWSAccount) (bool, error) {
	list := &awsv1alpha1.AWSAccountList{}

	if err := t.mgr.GetClient().List(ctx, list, client.InNamespace("")); err != nil {
		return false, err
	}

	// @step: we iterate the list and look for any aws account with the same name
	// but NOT in our namespace
	for _, x := range list.Items {
		if x.Namespace == account.Namespace && x.Name == account.Name {
			continue
		}

		if x.Spec.AccountName == account.Spec.AccountName || x.Status.AccountID == account.Status.AccountID {
			return true, nil
		}
	}

	return false, nil
}

// GetAllocationName returns the name we should use for the account allocation
func (t *awsCtrl) GetAllocationName(account *awsv1alpha1.AWSAccount) string {
	return fmt.Sprintf("aws-%s", account.Name)
}

// GetProjectCredentialsSecretName returns the secret name of the credentials for this project
func (t *awsCtrl) GetAccountCredentialsSecretName(account *awsv1alpha1.AWSAccount) string {
	return fmt.Sprintf("aws-%s", account.Name)
}

func (t *awsCtrl) DeriveAccountEmail(org *awsv1alpha1.AWSOrganization, account *awsv1alpha1.AWSAccount) string {
	components := strings.Split(org.Spec.SsoUser.Email, "@")
	username, domain := components[0], components[1]

	return fmt.Sprintf("%s+%s@%s", username, account.Spec.AccountName, domain)
}

// CreateAWSCredentialsSecret returns a project credentials secret
func CreateAWSCredentialsSecret(account *awsv1alpha1.AWSAccount, name string, values map[string]string) *configv1.Secret {
	secret := &configv1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.GroupVersion.String(),
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: account.Namespace,
		},
		Spec: configv1.SecretSpec{
			Data:        values,
			Description: fmt.Sprintf("AWS Account credentials for team account: %s", account.Spec.AccountName),
			Type:        configv1.GenericSecret,
		},
	}

	return secret
}
