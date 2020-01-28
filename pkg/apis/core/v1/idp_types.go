/*
 * Copyright (C) 2019 Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// IDPConfig represents a configuration required for any Identity Provider available
// Only a single identity provider config should be set
type IDPConfig struct {
	// Google represents a Google IDP config
	// +optional
	Github *GithubIDP `json:"github,omitempty"`
	Google *GoogleIDP `json:"google,omitempty"`
	SAML   *SAMLIDP   `json:"saml,omitempty"`
	OIDC   *OIDCIDP   `json:"oidc,omitempty"`
}

// IDPSpec defines the spec for a configured instance of an IDP
// +k8s:openapi-gen=true
type IDPSpec struct {
	// DisplayName
	// +kubebuilder:validation:Required
	DisplayName string `json:"displayName"`
	// IDPConfig
	// +kubebuilder:validation:Required
	Config IDPConfig `json:"config"`
}

// IDPStatus defines the observed state of an IDP (ID Providers)
// +k8s:openapi-gen=true
type IDPStatus struct {
	// Conditions is a set of condition which has caused an error
	// +kubebuilder:validation:Optional
	// +listType
	Conditions []Condition `json:"conditions"`
	// Status is overall status of the IDP configuration
	Status Status `json:"status"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IDP is the Schema for the class API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=idp,scope=Namespaced
type IDP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IDPSpec   `json:"spec,omitempty"`
	Status IDPStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IDPList contains a list of IDProvider
type IDPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IDP `json:"items"`
}

// GithubIDP provides config for a github OAuth app identity provider
type GithubIDP struct {
	// ClientID is the field name in a Github OAuth app
	ClientID string `json:"clientID"`
	// ClientSecret is the field name in a Github OAuth app
	ClientSecret string `json:"clientSecret"`
	// ClientSecret is the list of possible Organisations in Github the user must be part of
	Orgs []string `json:"orgs"`
}

// GoogleIDP provides config for a Google Identity provider
type GoogleIDP struct {
	// ClientID is the field name in a Google OAuth app
	ClientID string `json:"clientID"`
	// ClientSecret is the field name in a Google OAuth app
	ClientSecret string `json:"clientSecret"`
	// Domains are the google accounts whitelisted for authentication
	Domains []string `json:"domains"`
}

// OIDCIDP config for a generoc Open ID Connect provider
type OIDCIDP struct {
	// ClientID provides the OIDC client ID string
	ClientID string `json:"clientID"`
	// ClientSecret provides the OIDC client secret string
	ClientSecret string `json:"clientSecret"`
	// Issuer provides the IDP URL
	Issuer string `json:"issuer"`
}

// SAMLIDP provides configuration for a generic SAML Identity provider
type SAMLIDP struct {
	// SSOURL provides the SSO URL used for POST value to IDP
	SSOURL string `json:"ssoURL"`
	// CAData is byte array representing the PEM data for the IDP signing CA
	CAData []byte `json:"caData"`
	// UsernameAttr attribute in the returned assertion to map to ID token claims
	UsernameAttr string `json:"usernameAttr"`
	// EmailAttr attribute in the returned assertion to map to ID token claims
	EmailAttr string `json:"emailAttr"`
	// GroupsAttr attribute in the returned assertion to map to ID token claims
	GroupsAttr string `json:"groupsAttr,omitempty"`
	// AllowedGroups provides a list of allowed groups
	AllowedGroups []string `json:"allowedGroups,omitempty"`
	// GroupsDelim characters used to split the single groups field to obtain the user group membership
	GroupsDelim string `json:"groupsDelim,omitempty"`
}
