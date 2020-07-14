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

package awsorganization

import (
	"context"

	awsv1alpha1 "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	"github.com/appvia/kore/pkg/utils/cloud/aws"

	log "github.com/sirupsen/logrus"
)

// ValidateRoleAndOUName is responsible for checking the permissions are correct
func (t awsCtrl) ValidateRoleAndOUName(
	ctx context.Context,
	org *awsv1alpha1.AWSOrganization,
	credentials *aws.Credentials) error {

	logger := log.WithFields(log.Fields{
		"name":      org.Name,
		"namespace": org.Namespace,
	})
	logger.Debug("checking the credentials and role for the aws organization are correct")

	s := aws.AssumeRoleFromCreds(*credentials, org.Spec.RoleARN, org.Spec.Region, org.Spec.Region)
	_, err := aws.GetOUID(s, org.Spec.OuName)
	if err != nil {
		// TODO set conditions

		return err
	}
	// All good - ensure this gets saved
	org.Status.AccountID = credentials.AccountID

	return nil
}
