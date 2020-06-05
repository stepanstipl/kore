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
// +build awse2e

package aws

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEnsureIRSA will test that we can create the IAM association with a known test cluster...
func TestEnsureIRSA(t *testing.T) {
	iamClient := NewIamClient(Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("KORETEST_AWS_ACCOUNTID"),
	})

	err := iamClient.EnsureIRSA(os.Getenv("KORETEST_CLUSTER_ARN"), os.Getenv("KORETEST_CLUSTER_OIDC"))
	require.NoError(t, err)
}
