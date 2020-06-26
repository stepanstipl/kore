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

package test

import (
	"os"
	"testing"
)

var (
	// AccessKeyID is used for testing aws crds
	AccessKeyID = os.Getenv("KORETEST_AWSE2E_AWS_ACCESS_KEY_ID")
	// SecretAccessKey is used for testing aws creds
	SecretAccessKey = os.Getenv("KORETEST_AWSE2E_AWS_SECRET_ACCESS_KEY")
	// AccountID is used for testing aws creds
	AccountID = os.Getenv("KORETEST_AWS_ACCOUNT_ID")
)

// SkipTestIfEnvNotSet will decide if this test should run
func SkipTestIfEnvNotSet(testName string, t *testing.T) {
	n := "KORETEST_AWSE2E_" + testName
	v := GetEnvAndTestLog(n, t)
	if v != "true" {
		t.Skipf("skipping test as %s is not set to 'true' but '%s'", n, v)
	}
}

// GetEnvAndTestLog get an environment variable and log it
func GetEnvAndTestLog(envName string, t *testing.T) string {
	v := os.Getenv(envName)
	t.Logf("%s=%s", envName, v)
	return v
}
