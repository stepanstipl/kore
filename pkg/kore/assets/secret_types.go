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

package assets

// SecretTypes is a secret type
type SecretTypes struct {
	// Name is the name of the secret
	Name string
	// Fields is required list of fields for this secret
	Fields []string
}

// @TODO potentially move this to JSONSchema?
var (
	// AWSCredentialsSecret holds kubernetes endpoints
	AWSCredentialsSecret = SecretTypes{
		Name:   "aws-credentials",
		Fields: []string{"access_key_id", "access_secret_key"},
	}

	// GenericSecret is used to generic secrets
	GenericSecret = SecretTypes{
		Name:   "generic",
		Fields: []string{},
	}

	// GKECredentialsSecret holds the gke credentials
	GKECredentialsSecret = SecretTypes{
		Name:   "gke-credentials",
		Fields: []string{"service_account_key"},
	}

	// GCPProjectSecret holds the details related to a gcp project credential
	GCPProjectSecret = SecretTypes{
		Name:   "gcp-project",
		Fields: []string{"expires", "project_id", "project", "key", "key_id"},
	}

	// GCPOrganizationalSecret holds the SA for gcp organization
	GCPOrganizationalSecret = SecretTypes{
		Name:   "gcp-org",
		Fields: []string{"key", "org-id"},
	}

	// KubernetesSecret holds kubernetes endpoints
	KubernetesSecret = SecretTypes{
		Name:   "kubernetes",
		Fields: []string{"ca.crt", "endpoint", "token"},
	}

	// AzureSecret holds Azure API credentials
	AzureSecret = SecretTypes{
		Name:   "azure-credentials",
		Fields: []string{"subscription_id", "tenant_id", "client_id", "client_secret"},
	}
)

// GetSecretTypeOrGeneric returns the secret type
func GetSecretTypeOrGeneric(name string) SecretTypes {
	for _, x := range SupportedSecretTypes() {
		if x.Name == name {
			return x
		}
	}

	return GenericSecret
}

// SupportedSecretTypes returns a list of secret types
func SupportedSecretTypes() []SecretTypes {
	return []SecretTypes{
		AWSCredentialsSecret,
		AzureSecret,
		GCPOrganizationalSecret,
		GCPProjectSecret,
		GenericSecret,
		GKECredentialsSecret,
		KubernetesSecret,
	}
}

// SupportedSecretTypesNames returns a list of secret type names
func SupportedSecretTypesNames() []string {
	list := make([]string, len(SupportedSecretTypes()))
	types := SupportedSecretTypes()

	for i := 0; i < len(list); i++ {
		list[i] = types[i].Name
	}

	return list
}
