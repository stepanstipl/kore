package credentials

import (
	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
)

// MaxChunkSize is the largest number of permissions that can be checked in one request
const MaxChunkSize = 100

type awsClient struct {
	credentials *aws.AWSCredential
}

// NewClient creates and returns a permissions verifier
func NewClient(credentials *aws.AWSCredential) (*awsClient, error) {
	awsClient := &awsClient{}

	// Example from GKE credentials...
	/*
		options := []option.ClientOption{option.WithCredentialsJSON([]byte(credentials.Spec.Account))}

		crm, err := resourcemanager.NewService(context.Background(), options...)
		if err != nil {
			return nil, err
		}
	*/
	return awsClient, nil
}

// HasRequiredPermissions tests whether a serviceaccount has the required permissions for cluster manager
func (c *awsClient) HasRequiredPermissions() (bool, error) {
	// TODO work out AWS equivalent of IAM API verification
	return true, nil
}
