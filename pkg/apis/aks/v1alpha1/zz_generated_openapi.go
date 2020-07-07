// +build !ignore_autogenerated

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

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKS":                  schema_pkg_apis_aks_v1alpha1_AKS(ref),
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSCredentials":       schema_pkg_apis_aks_v1alpha1_AKSCredentials(ref),
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSCredentialsSpec":   schema_pkg_apis_aks_v1alpha1_AKSCredentialsSpec(ref),
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSCredentialsStatus": schema_pkg_apis_aks_v1alpha1_AKSCredentialsStatus(ref),
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSSpec":              schema_pkg_apis_aks_v1alpha1_AKSSpec(ref),
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSStatus":            schema_pkg_apis_aks_v1alpha1_AKSStatus(ref),
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.LinuxProfile":         schema_pkg_apis_aks_v1alpha1_LinuxProfile(ref),
		"github.com/appvia/kore/pkg/apis/aks/v1alpha1.WindowsProfile":       schema_pkg_apis_aks_v1alpha1_WindowsProfile(ref),
	}
}

func schema_pkg_apis_aks_v1alpha1_AKS(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AKS is the schema for an AKS cluster object",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSSpec", "github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_aks_v1alpha1_AKSCredentials(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AKSCredentials are used for storing Azure credentials needed to create AKS clusters",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSCredentialsSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSCredentialsStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSCredentialsSpec", "github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSCredentialsStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_aks_v1alpha1_AKSCredentialsSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AKSCredentialsSpec defines the desired state of AKSCredentials",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"subscriptionID": {
						SchemaProps: spec.SchemaProps{
							Description: "SubscriptionID is the Azure Subscription ID",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"tenantID": {
						SchemaProps: spec.SchemaProps{
							Description: "TenantID is the Azure Tenant ID",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"clientID": {
						SchemaProps: spec.SchemaProps{
							Description: "ClientID is the Azure client ID",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"credentialsRef": {
						SchemaProps: spec.SchemaProps{
							Description: "CredentialsRef is a reference to the credentials used to create clusters",
							Ref:         ref("k8s.io/api/core/v1.SecretReference"),
						},
					},
				},
				Required: []string{"subscriptionID", "tenantID", "clientID"},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.SecretReference"},
	}
}

func schema_pkg_apis_aks_v1alpha1_AKSCredentialsStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AKSCredentialsStatus defines the observed state of AKSCredentials",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a collection of potential issues",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/appvia/kore/pkg/apis/core/v1.Condition"),
									},
								},
							},
						},
					},
					"verified": {
						SchemaProps: spec.SchemaProps{
							Description: "Verified checks that the credentials are ok and valid",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status provides a overall status",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Condition"},
	}
}

func schema_pkg_apis_aks_v1alpha1_AKSSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AKSSpec defines the desired state of an AKS cluster",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"apiServerAuthorizedIPRanges": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "APIServerAuthorizedIPRanges are IP ranges to whitelist for incoming traffic to the API servers",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"cluster": {
						SchemaProps: spec.SchemaProps{
							Description: "Cluster refers to the cluster this object belongs to",
							Ref:         ref("github.com/appvia/kore/pkg/apis/core/v1.Ownership"),
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Description provides a short summary / description of the cluster.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"dnsPrefix": {
						SchemaProps: spec.SchemaProps{
							Description: "DNSPrefix is the DNS prefix for the cluster Must contain between 3 and 45 characters, and can contain only letters, numbers, and hyphens. It must start with a letter and must end with a letter or a number.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"enablePodSecurityPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "EnablePodSecurityPolicy indicates whether Pod Security Policies should be enabled Note that this also requires role based access control to be enabled. This feature is currently in preview and PodSecurityPolicyPreview for namespace Microsoft.ContainerService must be enabled.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "Version is the Kubernetes version",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"linuxProfile": {
						SchemaProps: spec.SchemaProps{
							Description: "LinuxProfile is the configuration for Linux VMs",
							Ref:         ref("github.com/appvia/kore/pkg/apis/aks/v1alpha1.LinuxProfile"),
						},
					},
					"networkPlugin": {
						SchemaProps: spec.SchemaProps{
							Description: "NetworkPlugin is the network plugin to use for networking. \"azure\" or \"kubenet\"",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"networkPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "NetworkPolicy is the network policy to use for networking. \"azure\" or \"calico\"",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"nodePools": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "NodePools is the set of node pools for this cluster. Required unless ALL deprecated properties except subnetwork are set.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSNodePool"),
									},
								},
							},
						},
					},
					"privateClusterEnabled": {
						SchemaProps: spec.SchemaProps{
							Description: "PrivateClusterEnabled controls whether the Kubernetes API is only exposed on the private network",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"region": {
						SchemaProps: spec.SchemaProps{
							Description: "Region is the location where the AKS cluster should be created",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"windowsProfile": {
						SchemaProps: spec.SchemaProps{
							Description: "WindowsProfile is the configuration for Windows VMs",
							Ref:         ref("github.com/appvia/kore/pkg/apis/aks/v1alpha1.WindowsProfile"),
						},
					},
				},
				Required: []string{"description", "dnsPrefix", "networkPlugin", "networkPolicy", "nodePools", "region"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/aks/v1alpha1.AKSNodePool", "github.com/appvia/kore/pkg/apis/aks/v1alpha1.LinuxProfile", "github.com/appvia/kore/pkg/apis/aks/v1alpha1.WindowsProfile", "github.com/appvia/kore/pkg/apis/core/v1.Ownership"},
	}
}

func schema_pkg_apis_aks_v1alpha1_AKSStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "AKSStatus defines the observed state of an AKS cluster",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"components": {
						SchemaProps: spec.SchemaProps{
							Description: "Components is the status of the components",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/appvia/kore/pkg/apis/core/v1.Component"),
									},
								},
							},
						},
					},
					"caCertificate": {
						SchemaProps: spec.SchemaProps{
							Description: "CACertificate is the certificate for this cluster",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"endpoint": {
						SchemaProps: spec.SchemaProps{
							Description: "Endpoint is the endpoint of the cluster",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status provides the overall status",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Message is the status message",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Component"},
	}
}

func schema_pkg_apis_aks_v1alpha1_LinuxProfile(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "LinuxProfile is the configuration for Linux VMs",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"adminUsername": {
						SchemaProps: spec.SchemaProps{
							Description: "AdminUsername is the admin username for Linux VMs",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"sshPublicKeys": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "set",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "SSHPublicKeys is a list of public SSH keys to allow to connect to the Linux VMs",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
				},
				Required: []string{"adminUsername", "sshPublicKeys"},
			},
		},
	}
}

func schema_pkg_apis_aks_v1alpha1_WindowsProfile(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "WindowsProfile is the configuration for Windows VMs",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"adminUsername": {
						SchemaProps: spec.SchemaProps{
							Description: "AdminUsername is the admin username for Windows VMs",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"adminPassword": {
						SchemaProps: spec.SchemaProps{
							Description: "AdminPassword is the admin password for Windows VMs",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"adminUsername", "adminPassword"},
			},
		},
	}
}
