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

package aws

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/appvia/kore/pkg/utils/cloud/aws/test"

	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type testAccountData struct {
	region      string
	roleARN     string
	newRegion   string
	session     *session.Session
	accountName string
	ssoEmail    string
	account     Account
}

const (
	testAccountsPrefix = "ACCOUNTS"
)

var (
	accountCreationRecordID *string
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)
}

// TestCreateAccount will test that we can assume an AWS role
func TestCreateAccount(t *testing.T) {
	test.SkipTestIfEnvNotSet(testAccountsPrefix+"_CREATE", t)

	a := getAccountClientForTest(t)

	accountCreationRecordID, err := a.CreateNewAccount()
	t.Logf("token for account - %s", accountCreationRecordID)
	require.NoError(t, err)

	timeout, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()
	err = a.WaitForAccountAvailable(timeout, accountCreationRecordID)
	require.NoError(t, err)
}

func TestEnsureInitialAccess(t *testing.T) {
	test.SkipTestIfEnvNotSet(testAccountsPrefix, t)

	// Test bad account name
	td := getTestDataForAccountClient(t)
	td.account.NewAccountName = "notpresent"
	a := NewAccountClientFromSessionAndRole(td.session, td.roleARN, td.newRegion, td.account)
	err := a.EnsureInitialAccessCreated()
	require.Error(t, err)

	// Test bad OU Name
	td = getTestDataForAccountClient(t)
	td.account.ManagedOrganizationalUnit = "notpresent"
	a = NewAccountClientFromSessionAndRole(td.session, td.roleARN, td.newRegion, td.account)
	err = a.EnsureInitialAccessCreated()
	require.Error(t, err)

	// Test with valid data
	a = getAccountClientForTest(t)
	err = a.EnsureInitialAccessCreated()
	require.NoError(t, err)
}

func TestWaitForInitialAccess(t *testing.T) {
	test.SkipTestIfEnvNotSet(testAccountsPrefix, t)

	timeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Test with valid data
	a := getAccountClientForTest(t)
	err := a.WaitForInitialAccess(timeout)
	require.NoError(t, err)

	log.Debugf("Account ID created:%s", *a.account.id)
}

func TestCreateAccountCredentials(t *testing.T) {
	test.SkipTestIfEnvNotSet(testAccountsPrefix+"_CREATE_USER_CREDS", t)

	a := getAccountClientForTest(t)
	creds, err := a.CreateAccountCredentials()
	log.Debugf("creds:%v", creds)
	require.NoError(t, err)

	i := NewIamClient(*creds)
	time.Sleep(time.Second * 30)
	_, err = i.GetARN()

	require.NoError(t, err)
}

func getAccountClientForTest(t *testing.T) *AccountClient {
	td := getTestDataForAccountClient(t)
	return NewAccountClientFromSessionAndRole(td.session, td.roleARN, td.newRegion, td.account)
}

func getTestDataForAccountClient(t *testing.T) testAccountData {
	region := test.GetEnvAndTestLog("KORETEST_REGION", t)
	accountName := test.GetEnvAndTestLog("KORETEST_ACCOUNT_NAME", t)
	ssoEmail := test.GetEnvAndTestLog("KORETEST_ACCOUNT_SSO_EMAIL", t)
	components := strings.Split(ssoEmail, "@")
	username, domain := components[0], components[1]

	return testAccountData{
		region:    region,
		roleARN:   test.GetEnvAndTestLog("KORETEST_ASSUME_ROLE_ARN", t),
		newRegion: test.GetEnvAndTestLog("KORETEST_ASSUME_ROLE_NEW_REGION", t),
		session: getNewSession(
			Credentials{
				AccessKeyID:     test.AccessKeyID,
				SecretAccessKey: test.SecretAccessKey,
				AccountID:       os.Getenv("KORETEST_STS_ACCOUNT_ID"),
			},
			region,
		),
		account: Account{
			AccountEmail:              fmt.Sprintf("%s+%s@%s", username, accountName, domain),
			ManagedOrganizationalUnit: test.GetEnvAndTestLog("KORETEST_ACCOUNT_OU", t),
			NewAccountName:            accountName,
			SSOUserEmail:              ssoEmail,
			SSOUserFirstName:          test.GetEnvAndTestLog("KORETEST_ACCOUNT_SSO_FIRSTNAME", t),
			SSOUserLastName:           test.GetEnvAndTestLog("KORETEST_ACCOUNT_SSO_LASTNAME", t),
		},
	}
}
