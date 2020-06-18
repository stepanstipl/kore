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
	"context"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEnsureIRSA will test that we can create the IAM association with a known test cluster...
func TestEnsureIRSA(t *testing.T) {
	if !doAWSE2ETest() {
		t.Skip("skipping test - requires KORETEST_AWSE2E to run")
	}
	iamClient := getIamClientFromEnv()
	err := iamClient.EnsureIRSA(os.Getenv("KORETEST_CLUSTER_ARN"), os.Getenv("KORETEST_CLUSTER_OIDC_URL"))
	require.NoError(t, err)
}

func TestEnsureClusterAutoscalingRoleAndPolicies(t *testing.T) {
	if !doAWSE2ETest() {
		t.Skip("skipping test - requires KORETEST_AWSE2E to run")
	}
	iamClient := getIamClientFromEnv()

	ngas := []NodeGroupAutoScaler{
		{
			AutoScalingARN: os.Getenv("KORETEST_NG_1_AUTOSCALING_ARN"),
			NodeGroupName:  os.Getenv("KORETEST_NG_1_NAME"),
		},
		{
			AutoScalingARN: os.Getenv("KORETEST_NG_2_AUTOSCALING_ARN"),
			NodeGroupName:  os.Getenv("KORETEST_NG_2_NAME"),
		},
	}

	issuerURL, _ := url.Parse(os.Getenv("KORETEST_CLUSTER_OIDC_URL"))
	oidcIssuer := issuerURL.Hostname() + issuerURL.Path

	role, err := iamClient.EnsureClusterAutoscalingRoleAndPolicies(
		context.Background(),
		os.Getenv("KORETEST_CLUSTERNAME"),
		os.Getenv("KORETEST_AWS_ACCOUNTID"),
		oidcIssuer,
		ngas,
	)
	require.NoError(t, err)
	require.NotNil(t, role)
}

func getIamClientFromEnv() *IamClient {
	return NewIamClient(Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AccountID:       os.Getenv("KORETEST_AWS_ACCOUNTID"),
	})
}

func doAWSE2ETest() bool {
	return (os.Getenv("KORETEST_AWSE2E") == "true")
}
