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
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
)

const (
	serviceCatalogControlTowerPortfolioProviderName = "AWS Control Tower"
	serviceCatalogControlTowerProductName           = "AWS Control Tower Account Factory"

	// KoreAccountAdminStackSetName is the name of the stackset
	KoreAccountAdminStackSetName = "kore-admin-role-for-member-accounts"
	// KoreAccountsAdminRoleName is the iam role name to create in each new account
	KoreAccountsAdminRoleName = "kore-accounts-admin-automation-role"
	// KoreAccountAdminUserName is the username to use when creating credentials in the new accounts
	KoreAccountAdminUserName    = "kore-account-admin"
	adminPolicyARN              = "arn:aws:iam::aws:policy/AdministratorAccess"
	koreAccountsAdminCFTemplate = `{
		"AWSTemplateFormatVersion": "2010-09-09",
		"Description": "A role that grants administrative privileges to the kore-accounts-user to deploy Cloud services into the member accounts.",
		"Resources": {
		   "AutomationRole": {
			  "Type": "AWS::IAM::Role",
			  "Properties": {
				 "RoleName": "%s",
				 "AssumeRolePolicyDocument": {
					"Version": "2012-10-17",
					"Statement": [
					   {
						  "Effect": "Allow",
						  "Principal": {
							 "AWS": "%s"
						  },
						  "Action": [
							 "sts:AssumeRole"
						  ]
					   }
					]
				 },
				 "Path": "/",
				 "ManagedPolicyArns": [
					"arn:aws:iam::aws:policy/AdministratorAccess"
				 ]
			  }
		   }
		},
		"Outputs": {
		   "AutomationRole": {
			  "Description": "AutomationRole",
			  "Value": {
				  "Fn::GetAtt": [
					  "AutomationRole",
					  "Arn"
				  ]
			  }
		   }
		}
	 }`
)

// AccountClienter provides general access to methods for future testing
type AccountClienter interface {
	// Exists will check aws to see if the account exists
	Exists() (bool, error)
	// CreateNewAccount will create a new aws account and return a provisioning record id for checking status
	CreateNewAccount() (string, error)
	// WaitForAccountAvailable is used to wait for the account to be created
	WaitForAccountAvailable(ctx context.Context, provisionRecordID string) error
	// IsAccountReady will determine if the account provisioning is ready
	IsAccountReady(provisionRecordID string) (bool, error)
	// EnsureInitialAccessCreated will ensure access to the account
	EnsureInitialAccessCreated() error
	// IsInitialAccessReady checks if initial access is working
	IsInitialAccessReady() (bool, error)
	// WaitForInitialAccess is used to wait for the account to be created
	WaitForInitialAccess(ctx context.Context) error
	// DoAccountCredentialsExist will create or update any missing accounts
	CreateAccountCredentials() (*Credentials, error)
}

// Account provides the details required to create an account
type Account struct {
	AccountEmail              string
	SSOUserEmail              string
	SSOUserFirstName          string
	SSOUserLastName           string
	ManagedOrganizationalUnit string
	NewAccountName            string
	PrimaryResourceRegion     string
	id                        *string
}

// AccountClient provides access to account management methods
type AccountClient struct {
	session        *session.Session
	svc            *servicecatalog.ServiceCatalog
	roleARN        string
	account        Account
	region         string
	cfSvc          *cloudformation.CloudFormation
	sis            *cloudformation.StackInstanceSummary
	accountSession *session.Session
}

// Ensure we implement the public interface
var _ AccountClienter = (*AccountClient)(nil)

// NewAccountClientFromCredsAndRole will create a client
func NewAccountClientFromCredsAndRole(creds Credentials, roleARN, region string, a Account) *AccountClient {
	return NewAccountClientFromSessionAndRole(getNewSession(creds, region), roleARN, region, a)
}

// NewAccountClientFromSessionAndRole will create a client
func NewAccountClientFromSessionAndRole(s *session.Session, roleARN, region string, a Account) *AccountClient {
	newSession := AssumeRoleFromSession(s, region, roleARN)
	if a.PrimaryResourceRegion == "" {
		a.PrimaryResourceRegion = region
	}
	c := &AccountClient{
		session: newSession,
		roleARN: roleARN,
		region:  region,
		account: a,
	}
	c.svc = servicecatalog.New(c.session)

	return c
}

// Exists will check aws to see if the account exists
func (a *AccountClient) Exists() (bool, error) {
	if a.account.id != nil {

		// return quickly if we've done this already
		return true, nil
	}
	err := a.updateAccountIDIfRequired()
	if err != nil {

		// there's an error
		return false, err
	}
	if a.account.id != nil {

		// we have found an account id
		return true, nil
	}

	// no account id found, no error
	return false, nil
}

// CreateNewAccount will create a new aws account and return a provisioning record id for checking status
func (a *AccountClient) CreateNewAccount() (string, error) {
	parsedARN, err := arn.Parse(a.roleARN)
	if err != nil {
		return "", fmt.Errorf("unable to parse role arn %s for account id", a.roleARN)
	}

	// First ensure the portfolio can be found...
	po, err := a.svc.ListPortfolios(&servicecatalog.ListPortfoliosInput{
		AcceptLanguage: aws.String("en"),
	})
	if err != nil {

		return "", fmt.Errorf("role %s cannot list portfolios - %w", a.roleARN, err)
	}
	var portfolioDetail *servicecatalog.PortfolioDetail
	for _, pd := range po.PortfolioDetails {
		if *pd.ProviderName == serviceCatalogControlTowerPortfolioProviderName {
			if portfolioDetail != nil {

				return "", fmt.Errorf("found more than one portfolio with provider name - %s", serviceCatalogControlTowerPortfolioProviderName)
			}
			portfolioDetail = pd
			// continue searching to make sure we have a unque ID
		}
	}
	if portfolioDetail == nil {

		return "", fmt.Errorf("role %s cannot find portfolios for %s using %v", a.roleARN, serviceCatalogControlTowerPortfolioProviderName, po.PortfolioDetails)
	}
	// Now associate the portfolio with us...
	_, err = a.svc.AssociatePrincipalWithPortfolio(&servicecatalog.AssociatePrincipalWithPortfolioInput{
		PortfolioId:   portfolioDetail.Id,
		PrincipalARN:  &a.roleARN,
		PrincipalType: aws.String("IAM"),
	})
	if err != nil {

		return "", fmt.Errorf("role %s cannot associate portfolios %s to own role %w", serviceCatalogControlTowerPortfolioProviderName, a.roleARN, err)
	}

	// We should be able to find the right product now...
	spo, err := a.svc.SearchProducts(&servicecatalog.SearchProductsInput{
		Filters: map[string][]*string{
			"FullTextSearch": {
				aws.String(serviceCatalogControlTowerProductName),
			},
		},
	})
	if err != nil {

		return "", fmt.Errorf("unable to list products matching %s with error - %w", serviceCatalogControlTowerProductName, err)
	}
	var productID *string
	for _, pvso := range spo.ProductViewSummaries {
		if *pvso.Name == serviceCatalogControlTowerProductName {
			if productID != nil {

				return "", fmt.Errorf("found more than one product with name - %s", serviceCatalogControlTowerPortfolioProviderName)
			}
			productID = pvso.ProductId
			// continue searching to make sure we have a unque ID
		}
	}
	if productID == nil {

		return "", fmt.Errorf("role %s cannot find product %s - %v", serviceCatalogControlTowerProductName, a.roleARN, spo.ProductViewSummaries)
	}
	// Get the provisioning artifact ID
	dpo, err := a.svc.DescribeProduct(&servicecatalog.DescribeProductInput{
		Id: productID,
	})
	if err != nil {

		return "", fmt.Errorf("role %s cannot descibe product %s - %w", a.roleARN, *productID, err)
	}
	// Find the provisioning artifact labeled "DEFAULT"
	var paID *string
	for _, pa := range dpo.ProvisioningArtifacts {
		if aws.StringValue(pa.Guidance) == "DEFAULT" {
			if paID != nil {

				return "", fmt.Errorf("more than one provisioning artifact with Guidance set to DEFAULT for product %s with id %s", serviceCatalogControlTowerProductName, *productID)
			}
			paID = pa.Id
		}
	}
	if paID == nil {

		return "", fmt.Errorf("no provisioning artifacts for with Guidance set to DEFAULT for product %s with id %s", serviceCatalogControlTowerProductName, *productID)
	}
	// Now get the launch paths for the Account Factory product...
	lpo, err := a.svc.ListLaunchPaths(&servicecatalog.ListLaunchPathsInput{
		ProductId: productID,
	})
	if err != nil {
		return "", fmt.Errorf("cannot list launch paths required for product id %s - %w", *productID, err)
	}
	var launchPathID *string
	for _, lps := range lpo.LaunchPathSummaries {
		log.Debugf("launch path name '%s' with launch path id '%s'", *lps.Name, *lps.Id)
		if aws.StringValue(lps.Name) == aws.StringValue(portfolioDetail.DisplayName) {
			launchPathID = lps.Id
		}
	}
	if launchPathID == nil {

		return "", fmt.Errorf("unable to find a launch path ID with display name %s", *portfolioDetail.DisplayName)
	}

	// Not for crypto just for api requests
	st := getAwsStringToken()
	provisionedProduct := "catalog-for-" + a.account.NewAccountName

	log.Debugf("provisioning account %s", a.account.NewAccountName)
	// Now time to provision an account
	ppo, err := a.svc.ProvisionProduct(&servicecatalog.ProvisionProductInput{
		ProductId:              productID,
		ProvisioningArtifactId: paID,
		ProvisionToken:         st,
		PathId:                 launchPathID,
		ProvisionedProductName: &provisionedProduct,
		ProvisioningPreferences: &servicecatalog.ProvisioningPreferences{
			StackSetAccounts: []*string{
				aws.String(parsedARN.AccountID),
			},
			StackSetRegions: []*string{
				aws.String(a.region),
			},
		},
		ProvisioningParameters: []*servicecatalog.ProvisioningParameter{
			{
				Key:   aws.String("AccountName"),
				Value: &a.account.NewAccountName,
			},
			{
				Key:   aws.String("SSOUserEmail"),
				Value: &a.account.SSOUserEmail,
			},
			{
				Key:   aws.String("AccountEmail"),
				Value: &a.account.AccountEmail,
			},
			{
				Key:   aws.String("SSOUserFirstName"),
				Value: &a.account.SSOUserFirstName,
			},
			{
				Key:   aws.String("SSOUserLastName"),
				Value: &a.account.SSOUserLastName,
			},
			{
				Key:   aws.String("ManagedOrganizationalUnit"),
				Value: &a.account.ManagedOrganizationalUnit,
			},
		},
	})
	if err != nil {

		return "", fmt.Errorf("unable to provision new product %s with id %s - %w", serviceCatalogControlTowerProductName, *productID, err)
	}
	log.Debugf("provisioning record - %v", ppo.RecordDetail)

	return aws.StringValue(ppo.RecordDetail.RecordId), nil
}

// WaitForAccountAvailable is used to wait for the account to be created
func (a *AccountClient) WaitForAccountAvailable(ctx context.Context, provisionRecordID string) error {
	for {
		// @step: we break out or continue
		select {
		case <-ctx.Done():

			return context.DeadlineExceeded
		default:
		}
		ready, err := a.IsAccountReady(provisionRecordID)
		if err != nil {

			return err
		}
		if ready {

			return nil
		}

		time.Sleep(15 * time.Second)
	}
}

// IsAccountReady will determine if the account provisioning is ready
func (a *AccountClient) IsAccountReady(provisionRecordID string) (bool, error) {
	pro, err := a.svc.DescribeRecord(&servicecatalog.DescribeRecordInput{
		Id: aws.String(provisionRecordID),
	})
	if err != nil {

		return false, err
	}
	log.Debugf("account provisioning status: %s", aws.StringValue(pro.RecordDetail.Status))
	switch aws.StringValue(pro.RecordDetail.Status) {
	case servicecatalog.RecordStatusSucceeded:

		return true, nil
	case servicecatalog.RecordStatusFailed:
	case servicecatalog.RecordStatusInProgressInError:

		return false, fmt.Errorf("account provisioning failed - %v", pro.RecordDetail.RecordErrors)
	default:
		log.Debugf("unknown account provisioning status: %s", aws.StringValue(pro.RecordDetail.Status))
	}

	return false, nil
}

// EnsureInitialAccessCreated will ensure access to the account
func (a *AccountClient) EnsureInitialAccessCreated() error {
	cfSvc := a.getCfSvc()
	accountExists, err := a.Exists()
	if err != nil {

		return fmt.Errorf("cannot determine if account exists - %w", err)
	}
	if !accountExists {

		return fmt.Errorf("cannot check access until account exists")
	}
	ouID, err := GetOUID(a.session, a.account.ManagedOrganizationalUnit)
	if err != nil {
		return fmt.Errorf("unable to oibtain ou id from ou %s - %w", a.account.ManagedOrganizationalUnit, err)
	}

	_, err = cfSvc.DescribeStackSet(&cloudformation.DescribeStackSetInput{
		StackSetName: aws.String(KoreAccountAdminStackSetName),
	})
	stackSetExists := true
	if err != nil {
		if !isAWSErr(err, cloudformation.ErrCodeStackSetNotFoundException, "") {

			return fmt.Errorf("cannot query for stackset - %w", err)
		}
		stackSetExists = false
	}
	if !stackSetExists {

		// First we create a "service managed" stackset to deploy an admin role we can assume
		// any kore managed identities (roles) will be created using this account
		template := fmt.Sprintf(koreAccountsAdminCFTemplate, KoreAccountsAdminRoleName, a.roleARN)
		_, err := cfSvc.CreateStackSet(&cloudformation.CreateStackSetInput{
			StackSetName: aws.String(KoreAccountAdminStackSetName),
			Description:  aws.String("Kore managed stackset to enable priovision an admin role to managed accounts"),
			AutoDeployment: &cloudformation.AutoDeployment{
				Enabled:                      aws.Bool(true),
				RetainStacksOnAccountRemoval: aws.Bool(true),
			},
			PermissionModel: aws.String(cloudformation.PermissionModelsServiceManaged),
			TemplateBody:    aws.String(template),
			Capabilities: aws.StringSlice([]string{
				cloudformation.CapabilityCapabilityNamedIam,
			}),
			// Retries are enabled with the default session for us
			// AWS wants a unique token so it knowns we are not creating different stacksets
			ClientRequestToken: getAwsStringToken(),
		})
		if err != nil {

			return fmt.Errorf("error creating stackset %s from template %s - %w", KoreAccountAdminStackSetName, template, err)
		}
	}
	// TODO: add else detect drift here (does stackset match latest definition)

	// Now ensure we have a stack instance for our account
	err = a.updateStackSetInstanceSummary()
	if err != nil {
		return err
	}
	if a.sis == nil {
		// create a stackset instance
		_, err := cfSvc.CreateStackInstances(&cloudformation.CreateStackInstancesInput{
			StackSetName: aws.String(KoreAccountAdminStackSetName),
			OperationId:  getAwsStringToken(),
			Regions: []*string{
				&a.account.PrimaryResourceRegion,
			},
			DeploymentTargets: &cloudformation.DeploymentTargets{
				OrganizationalUnitIds: []*string{
					ouID,
				},
			},
		})
		if err != nil {

			return fmt.Errorf("unable to create stackset instance for %s with account %s - %w", KoreAccountAdminStackSetName, *a.account.id, err)
		}
	}

	return nil
}

// WaitForInitialAccess is used to wait for the account to be created
func (a *AccountClient) WaitForInitialAccess(ctx context.Context) error {
	err := a.updateAccountIDIfRequired()
	if err != nil {
		return err
	}
	for {
		// @step: we break out or continue
		select {
		case <-ctx.Done():

			return context.DeadlineExceeded
		default:
		}
		ready, err := a.isInitialAccessReady()
		if err != nil {

			return err
		}
		if ready {

			return nil
		}

		time.Sleep(15 * time.Second)
	}
}

// IsInitialAccessReady will discover if the stacksets have deployed initial access
func (a *AccountClient) IsInitialAccessReady() (bool, error) {
	err := a.updateAccountIDIfRequired()
	if err != nil {
		return false, err
	}

	return a.isInitialAccessReady()
}

// DoAccountCredentialsExist will create or update any missing accounts
func (a *AccountClient) DoAccountCredentialsExist() (bool, error) {
	s, err := a.assumeAccountRole()
	if err != nil {

		return false, err
	}
	i := iam.New(s)
	_, err = i.GetUser(&iam.GetUserInput{
		UserName: aws.String(KoreAccountAdminUserName),
	})
	if err != nil {
		if isAWSErr(err, iam.ErrCodeNoSuchEntityException, "") {
			return false, nil
		}
		return false, fmt.Errorf("error checking if user %s exists in account %s - %w", KoreAccountAdminUserName, a.account.NewAccountName, err)
	}
	return true, nil
}

// CreateAccountCredentials will create an aws user
func (a *AccountClient) CreateAccountCredentials() (*Credentials, error) {
	s, err := a.assumeAccountRole()
	if err != nil {

		return nil, err
	}
	i := iam.New(s)
	_, err = i.CreateUser(&iam.CreateUserInput{
		UserName: aws.String(KoreAccountAdminUserName),
	})
	if err != nil {

		return nil, fmt.Errorf("error creating user %s in account %s - %w", KoreAccountAdminUserName, a.account.NewAccountName, err)
	}

	// Attach admin policy for account scoped creds:
	_, err = i.AttachUserPolicy(&iam.AttachUserPolicyInput{
		PolicyArn: aws.String(adminPolicyARN),
		UserName:  aws.String(KoreAccountAdminUserName),
	})
	if err != nil {

		return nil, fmt.Errorf("unable to attach admin policy to user %s in account %s - %w", KoreAccountAdminUserName, a.account.NewAccountName, err)
	}

	cao, err := i.CreateAccessKey(&iam.CreateAccessKeyInput{
		UserName: aws.String(KoreAccountAdminUserName),
	})
	if err != nil {

		return nil, fmt.Errorf("error creating user credentials for user %s in account %s - %w", KoreAccountAdminUserName, a.account.NewAccountName, err)
	}
	err = a.updateAccountIDIfRequired()
	if err != nil {

		return nil, err
	}

	return &Credentials{
		AccessKeyID:     aws.StringValue(cao.AccessKey.AccessKeyId),
		SecretAccessKey: aws.StringValue(cao.AccessKey.SecretAccessKey),
		AccountID:       aws.StringValue(a.account.id),
	}, nil
}

// isInitialAccessReady will determine if we are ready to assume role in new account
func (a *AccountClient) isInitialAccessReady() (bool, error) {
	ready, err := a.isStackSetInstanceReady(true)
	if err != nil {

		return false, fmt.Errorf("unbale to check if instance for %s with account %s is ready - %w", KoreAccountAdminStackSetName, *a.account.id, err)
	}
	if !ready {

		return false, nil
	}
	// create a new session
	s, err := a.assumeAccountRole()
	if err != nil {

		return false, err
	}

	// try and assume to acccount:
	_, err = s.Config.Credentials.Get()
	if err != nil {

		return false, err
	}

	return true, nil
}

func (a *AccountClient) isStackSetInstanceReady(checkNow bool) (bool, error) {
	if checkNow {
		a.sis = nil
	}
	if a.sis == nil {
		err := a.updateStackSetInstanceSummary()
		if err != nil {

			return false, err
		}
		if a.sis == nil {

			return false, nil
		}
	}
	switch *a.sis.Status {
	case cloudformation.StackInstanceStatusCurrent:

		return true, nil
	case cloudformation.StackInstanceStatusInoperable:

		return false, fmt.Errorf("stacksetid %s is inoperable - %s", *a.sis.StackId, *a.sis.StatusReason)
	default:
	}

	return false, nil
}

func (a *AccountClient) updateStackSetInstanceSummary() error {
	sio, err := a.getCfSvc().ListStackInstances(&cloudformation.ListStackInstancesInput{
		StackSetName:         aws.String(KoreAccountAdminStackSetName),
		StackInstanceAccount: a.account.id,
		StackInstanceRegion:  &a.account.PrimaryResourceRegion,
	})
	if err != nil {
		if !isAWSErr(err, cloudformation.ErrCodeStackInstanceNotFoundException, "") {
			return fmt.Errorf("unable to list stackset instances for %s - %w", KoreAccountAdminStackSetName, err)
		}
	}
	// determine if we have an instance for this account
	for _, si := range sio.Summaries {
		if *si.Account == *a.account.id {
			a.sis = si
		}
	}
	return nil
}

func (a *AccountClient) getCfSvc() *cloudformation.CloudFormation {
	if a.cfSvc != nil {

		return a.cfSvc
	}
	a.cfSvc = cloudformation.New(a.session)

	return a.cfSvc
}

// updateAccountIDIfRequired will populate the account id if not already set
func (a *AccountClient) updateAccountIDIfRequired() error {
	if a.account.id != nil {

		return nil
	}

	orgSvc := organizations.New(a.session)

	var ao *organizations.ListAccountsOutput
	var ai *organizations.ListAccountsInput

	// Handle paginated aws api results
	for {
		if ao != nil {
			ai = &organizations.ListAccountsInput{
				NextToken: ao.NextToken,
			}
		} else {
			ai = &organizations.ListAccountsInput{}
		}
		var err error
		// Do not create new vars here (we have to use THIS ao on next iteration)
		ao, err = orgSvc.ListAccounts(ai)
		if err != nil {

			return fmt.Errorf("unable to list accounts - %w", err)
		}
		for _, acc := range ao.Accounts {
			if *acc.Name == a.account.NewAccountName {
				if a.account.id != nil {

					return fmt.Errorf("more than one account found with account name %s - %s and %s", a.account.NewAccountName, *a.account.id, *acc.Id)
				}
				a.account.id = acc.Id
			}
		}
		if aws.StringValue(ao.NextToken) == "" {

			break
		}
	}

	return nil
}

func (a *AccountClient) assumeAccountRole() (*session.Session, error) {
	if a.accountSession != nil {

		return a.accountSession, nil
	}
	err := a.updateAccountIDIfRequired()
	if err != nil {
		return nil, err
	}

	// create a new session
	newAccountARN := fmt.Sprintf("arn:aws:iam::%s:role/%s", *a.account.id, KoreAccountsAdminRoleName)
	a.accountSession = AssumeRoleFromSession(a.session, a.account.PrimaryResourceRegion, newAccountARN)

	return a.accountSession, nil
}

func getAwsStringToken() *string {
	// Not for crypto just for api requests
	rand.Seed(time.Now().UnixNano())

	return aws.String(fmt.Sprintf("%d", rand.Intn(99999999999)))
}
