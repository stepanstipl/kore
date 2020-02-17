// +build !ignore_autogenerated

/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */
// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/appvia/kore/pkg/apis/org/v1.Team":                 schema_pkg_apis_org_v1_Team(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamInvitation":       schema_pkg_apis_org_v1_TeamInvitation(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamInvitationSpec":   schema_pkg_apis_org_v1_TeamInvitationSpec(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamInvitationStatus": schema_pkg_apis_org_v1_TeamInvitationStatus(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamMember":           schema_pkg_apis_org_v1_TeamMember(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamMemberSpec":       schema_pkg_apis_org_v1_TeamMemberSpec(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamMemberStatus":     schema_pkg_apis_org_v1_TeamMemberStatus(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamSpec":             schema_pkg_apis_org_v1_TeamSpec(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.TeamStatus":           schema_pkg_apis_org_v1_TeamStatus(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.User":                 schema_pkg_apis_org_v1_User(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.UserSpec":             schema_pkg_apis_org_v1_UserSpec(ref),
		"github.com/appvia/kore/pkg/apis/org/v1.UserStatus":           schema_pkg_apis_org_v1_UserStatus(ref),
	}
}

func schema_pkg_apis_org_v1_Team(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Team is the Schema for the teams API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.TeamSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.TeamStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/org/v1.TeamSpec", "github.com/appvia/kore/pkg/apis/org/v1.TeamStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_org_v1_TeamInvitation(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamInvitation is the Schema for the teams API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.TeamInvitationSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.TeamInvitationStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/org/v1.TeamInvitationSpec", "github.com/appvia/kore/pkg/apis/org/v1.TeamInvitationStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_org_v1_TeamInvitationSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamInvitationSpec defines the desired state of Team",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"username": {
						SchemaProps: spec.SchemaProps{
							Description: "Username is the user being bound to the team",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"team": {
						SchemaProps: spec.SchemaProps{
							Description: "Team is the name of the team",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"username", "team"},
			},
		},
	}
}

func schema_pkg_apis_org_v1_TeamInvitationStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamInvitationStatus defines the observed state of Team",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a collection of possible errors",
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
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is the status of the resource",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"conditions", "status"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Condition"},
	}
}

func schema_pkg_apis_org_v1_TeamMember(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamMember is the Schema for the teams API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.TeamMemberSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.TeamMemberStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/org/v1.TeamMemberSpec", "github.com/appvia/kore/pkg/apis/org/v1.TeamMemberStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_org_v1_TeamMemberSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamMemberSpec defines the desired state of Team",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"roles": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Role is the role of the user in the team",
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
					"team": {
						SchemaProps: spec.SchemaProps{
							Description: "Team is the name of the team",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"username": {
						SchemaProps: spec.SchemaProps{
							Description: "Username is the user being bound to the team",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"roles", "team", "username"},
			},
		},
	}
}

func schema_pkg_apis_org_v1_TeamMemberStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamMemberStatus defines the observed state of Team",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a collection of possible errors",
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
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is the status of the resource",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"conditions", "status"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Condition"},
	}
}

func schema_pkg_apis_org_v1_TeamSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamSpec defines the desired state of Team",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"summary": {
						SchemaProps: spec.SchemaProps{
							Description: "Summary is a summary name for this team",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"description": {
						SchemaProps: spec.SchemaProps{
							Description: "Description is a description for the team",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"summary", "description"},
			},
		},
	}
}

func schema_pkg_apis_org_v1_TeamStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "TeamStatus defines the observed state of Team",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is a collection of possible errors",
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
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is the status of the resource",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"conditions", "status"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Condition"},
	}
}

func schema_pkg_apis_org_v1_User(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "User is the Schema for the users API",
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
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.UserSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/appvia/kore/pkg/apis/org/v1.UserStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/org/v1.UserSpec", "github.com/appvia/kore/pkg/apis/org/v1.UserStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_org_v1_UserSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "UserSpec defines the desired state of User",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"disabled": {
						SchemaProps: spec.SchemaProps{
							Description: "Disabled indicates if the user is disabled",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"email": {
						SchemaProps: spec.SchemaProps{
							Description: "Email is the email for the user",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"username": {
						SchemaProps: spec.SchemaProps{
							Description: "Username is the userame or identity for this user",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"username"},
			},
		},
	}
}

func schema_pkg_apis_org_v1_UserStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "UserStatus defines the observed state of User",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						VendorExtensible: spec.VendorExtensible{
							Extensions: spec.Extensions{
								"x-kubernetes-list-type": "",
							},
						},
						SchemaProps: spec.SchemaProps{
							Description: "Conditions is collection of potentials error causes",
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
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status provides an overview of the user status",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"conditions", "status"},
			},
		},
		Dependencies: []string{
			"github.com/appvia/kore/pkg/apis/core/v1.Condition"},
	}
}
