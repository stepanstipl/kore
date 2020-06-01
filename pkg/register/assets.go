// Code generated for package register by go-bindata DO NOT EDIT. (@generated)
// sources:
// deploy/crds/accounts.kore.appvia.io_accountmanagement.yaml
// deploy/crds/apps.kore.appvia.io_appdeployments.yaml
// deploy/crds/apps.kore.appvia.io_installplans.yaml
// deploy/crds/aws.compute.kore.appvia.io_eks.yaml
// deploy/crds/aws.compute.kore.appvia.io_ekscredentials.yaml
// deploy/crds/aws.compute.kore.appvia.io_eksnodegroups.yaml
// deploy/crds/aws.compute.kore.appvia.io_eksvpcs.yaml
// deploy/crds/clusters.compute.kore.appvia.io_clusters.yaml
// deploy/crds/clusters.compute.kore.appvia.io_kubernetes.yaml
// deploy/crds/clusters.compute.kore.appvia.io_managedclusterrole.yaml
// deploy/crds/clusters.compute.kore.appvia.io_managedclusterrolebinding.yaml
// deploy/crds/clusters.compute.kore.appvia.io_managedconfig.yaml
// deploy/crds/clusters.compute.kore.appvia.io_managedpodsecuritypoliies.yaml
// deploy/crds/clusters.compute.kore.appvia.io_managedrole.yaml
// deploy/crds/clusters.compute.kore.appvia.io_namespaceclaims.yaml
// deploy/crds/clusters.compute.kore.appvia.io_namespacepolicy.yaml
// deploy/crds/config.kore.appvia.io_allocations.yaml
// deploy/crds/config.kore.appvia.io_planpolicies.yaml
// deploy/crds/config.kore.appvia.io_plans.yaml
// deploy/crds/config.kore.appvia.io_secrets.yaml
// deploy/crds/core.kore.appvia.io_idp.yaml
// deploy/crds/core.kore.appvia.io_oidclient.yaml
// deploy/crds/gcp.compute.kore.appvia.io_organizations.yaml
// deploy/crds/gcp.compute.kore.appvia.io_projectclaims.yaml
// deploy/crds/gcp.compute.kore.appvia.io_projects.yaml
// deploy/crds/gke.compute.kore.appvia.io_gkecredentials.yaml
// deploy/crds/gke.compute.kore.appvia.io_gkes.yaml
// deploy/crds/org.kore.appvia.io_auditevents.yaml
// deploy/crds/org.kore.appvia.io_members.yaml
// deploy/crds/org.kore.appvia.io_teaminvitations.yaml
// deploy/crds/org.kore.appvia.io_teams.yaml
// deploy/crds/org.kore.appvia.io_users.yaml
// deploy/crds/security.kore.appvia.io_securityoverviews.yaml
// deploy/crds/security.kore.appvia.io_securityrules.yaml
// deploy/crds/security.kore.appvia.io_securityscanresults.yaml
// deploy/crds/services.kore.appvia.io_servicecredentials.yaml
// deploy/crds/services.kore.appvia.io_servicekinds.yaml
// deploy/crds/services.kore.appvia.io_serviceplans.yaml
// deploy/crds/services.kore.appvia.io_serviceproviders.yaml
// deploy/crds/services.kore.appvia.io_services.yaml
package register

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _crdsAccountsKoreAppviaIo_accountmanagementYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: accountmanagement.accounts.kore.appvia.io
spec:
  group: accounts.kore.appvia.io
  names:
    kind: AccountManagement
    listKind: AccountManagementList
    plural: accountmanagement
    singular: accountmanagement
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: AccountManagement is the Schema for the accounts API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: AccountManagementSpec defines the desired state of accounting
            for a provider I've a feeling this will probably need provider specific
            attributes are some point
          properties:
            organization:
              description: Organization is the underlying organizational resource
                (only require if more than one)
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            provider:
              description: Provider is the name of provider which maps to the cluster
                kind
              minLength: 1
              type: string
            rules:
              description: Rules is a set of rules for this provider
              items:
                description: AccountsRule defines a rule for the provider
                properties:
                  description:
                    description: Description provides an optional description for
                      the account rule
                    type: string
                  labels:
                    additionalProperties:
                      type: string
                    description: Labels a collection of labels to apply the account
                    type: object
                  name:
                    description: Name is the given name of the rule
                    minLength: 1
                    type: string
                  plans:
                    description: Plans is a list of plans permitted
                    items:
                      type: string
                    minItems: 1
                    type: array
                    x-kubernetes-list-type: set
                  prefix:
                    description: Prefix is a prefix for the account name
                    type: string
                  suffix:
                    description: Suffix is the applied suffix
                    type: string
                required:
                - name
                - plans
                type: object
              type: array
              x-kubernetes-list-type: set
          required:
          - provider
          type: object
        status:
          description: AccountManagementStatus defines the observed state of Allocation
          properties:
            conditions:
              description: Conditions is a collection of potential issues
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is the general status of the resource
              type: string
          type: object
      type: object
  version: v1beta1
  versions:
  - name: v1beta1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsAccountsKoreAppviaIo_accountmanagementYamlBytes() ([]byte, error) {
	return _crdsAccountsKoreAppviaIo_accountmanagementYaml, nil
}

func crdsAccountsKoreAppviaIo_accountmanagementYaml() (*asset, error) {
	bytes, err := crdsAccountsKoreAppviaIo_accountmanagementYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/accounts.kore.appvia.io_accountmanagement.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsAppsKoreAppviaIo_appdeploymentsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: appdeployments.apps.kore.appvia.io
spec:
  group: apps.kore.appvia.io
  names:
    kind: AppDeployment
    listKind: AppDeploymentList
    plural: appdeployments
    singular: appdeployment
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: AppDeployment is the Schema for the allocations API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: AppDeploymentSpec defines the desired state of Allocation
          properties:
            capabilities:
              description: Capabilities defines the features supported by the package
              items:
                type: string
              minItems: 1
              type: array
              x-kubernetes-list-type: set
            cluster:
              description: Cluster is the cluster the application should be deployed
                on
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            description:
              description: Decription is a longer description of what the application
                provides
              minLength: 1
              type: string
            keywords:
              description: Keywords keywords whuch describe the application
              items:
                type: string
              minItems: 1
              type: array
              x-kubernetes-list-type: set
            official:
              description: Official indicates if the applcation is officially published
                by Appvia
              type: boolean
            package:
              description: Package is the name of the resource being shared
              minLength: 1
              type: string
            replaces:
              description: Replaces indicates the version this replaces
              minLength: 1
              type: string
            source:
              description: Source is the source of the package
              minLength: 1
              type: string
            subscription:
              description: Subscription is the nature of upgrades i.e manual or automatic
              enum:
              - Automatic
              - Manual
              minLength: 1
              type: string
            summary:
              description: Summary is a summary of what the application is
              type: string
            values:
              description: Values are optional values suppilied to the application
                deployment
              x-kubernetes-preserve-unknown-fields: true
            vendor:
              description: Vendor is the entity whom published the package
              minLength: 1
              type: string
            version:
              description: Version is the version of the package to install
              minLength: 1
              type: string
          required:
          - description
          - keywords
          - official
          - package
          - replaces
          - source
          - subscription
          - summary
          - vendor
          - version
          type: object
        status:
          description: AppDeploymentStatus defines the observed state of Allocation
          properties:
            conditions:
              description: Conditions is a collection of potential issues
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            installPlan:
              description: InstallPlan in the name of the installplan which this deployment
                has deployed from
              type: string
            status:
              description: Status is the general status of the resource
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsAppsKoreAppviaIo_appdeploymentsYamlBytes() ([]byte, error) {
	return _crdsAppsKoreAppviaIo_appdeploymentsYaml, nil
}

func crdsAppsKoreAppviaIo_appdeploymentsYaml() (*asset, error) {
	bytes, err := crdsAppsKoreAppviaIo_appdeploymentsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/apps.kore.appvia.io_appdeployments.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsAppsKoreAppviaIo_installplansYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: installplans.apps.kore.appvia.io
spec:
  group: apps.kore.appvia.io
  names:
    kind: InstallPlan
    listKind: InstallPlanList
    plural: installplans
    singular: installplan
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: InstallPlan is the Schema for the allocations API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: InstallPlanSpec defines the desired state of Allocation
          properties:
            approved:
              description: Approved indicates if the update has been approved
              type: boolean
          type: object
        status:
          description: InstallPlanStatus defines the observed state of Allocation
          properties:
            conditions:
              description: Conditions is a collection of potential issues
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            deployed:
              description: Deployed is the applciation deployment parameters
              properties:
                apiVersion:
                  description: 'APIVersion defines the versioned schema of this representation
                    of an object. Servers should convert recognized schemas to the
                    latest internal value, and may reject unrecognized values. More
                    info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                  type: string
                kind:
                  description: 'Kind is a string value representing the REST resource
                    this object represents. Servers may infer this from the endpoint
                    the client submits requests to. Cannot be updated. In CamelCase.
                    More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                  type: string
                metadata:
                  type: object
                spec:
                  description: AppDeploymentSpec defines the desired state of Allocation
                  properties:
                    capabilities:
                      description: Capabilities defines the features supported by
                        the package
                      items:
                        type: string
                      minItems: 1
                      type: array
                      x-kubernetes-list-type: set
                    cluster:
                      description: Cluster is the cluster the application should be
                        deployed on
                      properties:
                        group:
                          description: Group is the api group
                          type: string
                        kind:
                          description: Kind is the name of the resource under the
                            group
                          type: string
                        name:
                          description: Name is name of the resource
                          type: string
                        namespace:
                          description: Namespace is the location of the object
                          type: string
                        version:
                          description: Version is the group version
                          type: string
                      required:
                      - group
                      - kind
                      - name
                      - namespace
                      - version
                      type: object
                    description:
                      description: Decription is a longer description of what the
                        application provides
                      minLength: 1
                      type: string
                    keywords:
                      description: Keywords keywords whuch describe the application
                      items:
                        type: string
                      minItems: 1
                      type: array
                      x-kubernetes-list-type: set
                    official:
                      description: Official indicates if the applcation is officially
                        published by Appvia
                      type: boolean
                    package:
                      description: Package is the name of the resource being shared
                      minLength: 1
                      type: string
                    replaces:
                      description: Replaces indicates the version this replaces
                      minLength: 1
                      type: string
                    source:
                      description: Source is the source of the package
                      minLength: 1
                      type: string
                    subscription:
                      description: Subscription is the nature of upgrades i.e manual
                        or automatic
                      enum:
                      - Automatic
                      - Manual
                      minLength: 1
                      type: string
                    summary:
                      description: Summary is a summary of what the application is
                      type: string
                    values:
                      description: Values are optional values suppilied to the application
                        deployment
                      x-kubernetes-preserve-unknown-fields: true
                    vendor:
                      description: Vendor is the entity whom published the package
                      minLength: 1
                      type: string
                    version:
                      description: Version is the version of the package to install
                      minLength: 1
                      type: string
                  required:
                  - description
                  - keywords
                  - official
                  - package
                  - replaces
                  - source
                  - subscription
                  - summary
                  - vendor
                  - version
                  type: object
                status:
                  description: AppDeploymentStatus defines the observed state of Allocation
                  properties:
                    conditions:
                      description: Conditions is a collection of potential issues
                      items:
                        description: Condition is a reason why something failed
                        properties:
                          detail:
                            description: Detail is a actual error which might contain
                              technical reference
                            type: string
                          message:
                            description: Message is a human readable message
                            type: string
                        required:
                        - detail
                        - message
                        type: object
                      type: array
                      x-kubernetes-list-type: set
                    installPlan:
                      description: InstallPlan in the name of the installplan which
                        this deployment has deployed from
                      type: string
                    status:
                      description: Status is the general status of the resource
                      type: string
                  type: object
              type: object
            status:
              description: Status is the general status of the resource
              type: string
            update:
              description: Update is the incoming deployment is requiring approval
              properties:
                apiVersion:
                  description: 'APIVersion defines the versioned schema of this representation
                    of an object. Servers should convert recognized schemas to the
                    latest internal value, and may reject unrecognized values. More
                    info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                  type: string
                kind:
                  description: 'Kind is a string value representing the REST resource
                    this object represents. Servers may infer this from the endpoint
                    the client submits requests to. Cannot be updated. In CamelCase.
                    More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                  type: string
                metadata:
                  type: object
                spec:
                  description: AppDeploymentSpec defines the desired state of Allocation
                  properties:
                    capabilities:
                      description: Capabilities defines the features supported by
                        the package
                      items:
                        type: string
                      minItems: 1
                      type: array
                      x-kubernetes-list-type: set
                    cluster:
                      description: Cluster is the cluster the application should be
                        deployed on
                      properties:
                        group:
                          description: Group is the api group
                          type: string
                        kind:
                          description: Kind is the name of the resource under the
                            group
                          type: string
                        name:
                          description: Name is name of the resource
                          type: string
                        namespace:
                          description: Namespace is the location of the object
                          type: string
                        version:
                          description: Version is the group version
                          type: string
                      required:
                      - group
                      - kind
                      - name
                      - namespace
                      - version
                      type: object
                    description:
                      description: Decription is a longer description of what the
                        application provides
                      minLength: 1
                      type: string
                    keywords:
                      description: Keywords keywords whuch describe the application
                      items:
                        type: string
                      minItems: 1
                      type: array
                      x-kubernetes-list-type: set
                    official:
                      description: Official indicates if the applcation is officially
                        published by Appvia
                      type: boolean
                    package:
                      description: Package is the name of the resource being shared
                      minLength: 1
                      type: string
                    replaces:
                      description: Replaces indicates the version this replaces
                      minLength: 1
                      type: string
                    source:
                      description: Source is the source of the package
                      minLength: 1
                      type: string
                    subscription:
                      description: Subscription is the nature of upgrades i.e manual
                        or automatic
                      enum:
                      - Automatic
                      - Manual
                      minLength: 1
                      type: string
                    summary:
                      description: Summary is a summary of what the application is
                      type: string
                    values:
                      description: Values are optional values suppilied to the application
                        deployment
                      x-kubernetes-preserve-unknown-fields: true
                    vendor:
                      description: Vendor is the entity whom published the package
                      minLength: 1
                      type: string
                    version:
                      description: Version is the version of the package to install
                      minLength: 1
                      type: string
                  required:
                  - description
                  - keywords
                  - official
                  - package
                  - replaces
                  - source
                  - subscription
                  - summary
                  - vendor
                  - version
                  type: object
                status:
                  description: AppDeploymentStatus defines the observed state of Allocation
                  properties:
                    conditions:
                      description: Conditions is a collection of potential issues
                      items:
                        description: Condition is a reason why something failed
                        properties:
                          detail:
                            description: Detail is a actual error which might contain
                              technical reference
                            type: string
                          message:
                            description: Message is a human readable message
                            type: string
                        required:
                        - detail
                        - message
                        type: object
                      type: array
                      x-kubernetes-list-type: set
                    installPlan:
                      description: InstallPlan in the name of the installplan which
                        this deployment has deployed from
                      type: string
                    status:
                      description: Status is the general status of the resource
                      type: string
                  type: object
              type: object
          required:
          - deployed
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsAppsKoreAppviaIo_installplansYamlBytes() ([]byte, error) {
	return _crdsAppsKoreAppviaIo_installplansYaml, nil
}

func crdsAppsKoreAppviaIo_installplansYaml() (*asset, error) {
	bytes, err := crdsAppsKoreAppviaIo_installplansYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/apps.kore.appvia.io_installplans.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsAwsComputeKoreAppviaIo_eksYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: eks.aws.compute.kore.appvia.io
spec:
  additionalPrinterColumns:
  - JSONPath: .spec.description
    description: A description of the EKS cluster
    name: Description
    type: string
  - JSONPath: .status.endpoint
    description: The endpoint of the eks cluster
    name: Endpoint
    type: string
  - JSONPath: .status.status
    description: The overall status of the cluster
    name: Status
    type: string
  group: aws.compute.kore.appvia.io
  names:
    kind: EKS
    listKind: EKSList
    plural: eks
    singular: eks
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: EKS is the Schema for the eksclusters API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: EKSSpec defines the desired state of EKSCluster
          properties:
            authorizedMasterNetworks:
              description: AuthorizedMasterNetworks is the network ranges which are
                permitted to access the EKS control plane endpoint i.e the managed
                one (not the authentication proxy)
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
            cluster:
              description: Cluster refers to the cluster this object belongs to
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            credentials:
              description: Credentials is a reference to an EKSCredentials object
                to use for authentication
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            region:
              description: Region is the AWS region to launch this cluster within
              type: string
            securityGroupIDs:
              description: SecurityGroupIds is a list of security group IDs
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
            subnetIDs:
              description: SubnetIds is a list of subnet IDs
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
            version:
              description: Version is the Kubernetes version to use
              minLength: 3
              type: string
          required:
          - credentials
          - region
          - subnetIDs
          type: object
        status:
          description: EKSStatus defines the observed state of EKS cluster
          properties:
            caCertificate:
              description: CACertificate is the certificate for this cluster
              type: string
            conditions:
              description: Conditions is the status of the components
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            endpoint:
              description: Endpoint is the endpoint of the cluster
              type: string
            roleARN:
              description: RoleARN is the role ARN which provides permissions to EKS
              type: string
            status:
              description: Status provides a overall status
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsAwsComputeKoreAppviaIo_eksYamlBytes() ([]byte, error) {
	return _crdsAwsComputeKoreAppviaIo_eksYaml, nil
}

func crdsAwsComputeKoreAppviaIo_eksYaml() (*asset, error) {
	bytes, err := crdsAwsComputeKoreAppviaIo_eksYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/aws.compute.kore.appvia.io_eks.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsAwsComputeKoreAppviaIo_ekscredentialsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: ekscredentials.aws.compute.kore.appvia.io
spec:
  group: aws.compute.kore.appvia.io
  names:
    kind: EKSCredentials
    listKind: EKSCredentialsList
    plural: ekscredentials
    singular: ekscredentials
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: EKSCredentials is the Schema for the ekscredentials API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: EKSCredentialsSpec defines the desired state of EKSCredential
          properties:
            accessKeyID:
              description: AccessKeyID is the AWS Access Key ID
              type: string
            accountID:
              description: AccountID is the AWS account these credentials reside within
              minLength: 3
              type: string
            credentialsRef:
              description: CredentialsRef is a reference to the credentials used to
                create clusters
              properties:
                name:
                  description: Name is unique within a namespace to reference a secret
                    resource.
                  type: string
                namespace:
                  description: Namespace defines the space within which the secret
                    name must be unique.
                  type: string
              type: object
            secretAccessKey:
              description: SecretAccessKey is the AWS Secret Access Key
              type: string
          required:
          - accountID
          type: object
        status:
          description: EKSCredentialsStatus defines the observed state of EKSCredential
          properties:
            conditions:
              description: Conditions is a collection of potential issues
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status provides a overall status
              type: string
            verified:
              description: Verified checks that the credentials are ok and valid
              type: boolean
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsAwsComputeKoreAppviaIo_ekscredentialsYamlBytes() ([]byte, error) {
	return _crdsAwsComputeKoreAppviaIo_ekscredentialsYaml, nil
}

func crdsAwsComputeKoreAppviaIo_ekscredentialsYaml() (*asset, error) {
	bytes, err := crdsAwsComputeKoreAppviaIo_ekscredentialsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/aws.compute.kore.appvia.io_ekscredentials.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsAwsComputeKoreAppviaIo_eksnodegroupsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: eksnodegroups.aws.compute.kore.appvia.io
spec:
  additionalPrinterColumns:
  - JSONPath: .spec.description
    description: A description of the EKS cluster nodegroup
    name: Description
    type: string
  - JSONPath: .status.status
    description: The overall status of the cluster nodegroup
    name: Status
    type: string
  group: aws.compute.kore.appvia.io
  names:
    kind: EKSNodeGroup
    listKind: EKSNodeGroupList
    plural: eksnodegroups
    singular: eksnodegroup
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: EKSNodeGroup is the Schema for the eksnodegroups API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: EKSNodeGroupSpec defines the desired state of EKSNodeGroup
          properties:
            amiType:
              description: AMIType is the AWS Machine Image type. We use a sensible
                default.
              type: string
            cluster:
              description: Cluster refers to the cluster this object belongs to
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            credentials:
              description: Credentials is a reference to an AWSCredentials object
                to use for authentication
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            desiredSize:
              description: DesiredSize is the number of nodes to attempt to use
              format: int64
              minimum: 1
              type: integer
            diskSize:
              format: int64
              minimum: 1
              type: integer
            eC2SSHKey:
              description: EC2SSHKey is the Amazon EC2 SSH key that provides access
                for SSH communication with the worker nodes in the managed node group
                https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html
              type: string
            instanceType:
              description: InstanceType is the EC2 machine type
              type: string
            labels:
              additionalProperties:
                type: string
              description: Labels are any custom kubernetes labels to apply to nodes
              type: object
            maxSize:
              description: MaxSize is the most nodes the nodegroups can grow to
              format: int64
              maximum: 100
              type: integer
            minSize:
              description: MinSize is the least nodes the nodegroups can shrink to
              format: int64
              minimum: 1
              type: integer
            region:
              description: Region is the AWS location to launch node group within,
                must match the region of the cluster
              type: string
            releaseVersion:
              description: ReleaseVersion is release version of the managed node ami
              type: string
            sshSourceSecurityGroups:
              description: SSHSourceSecurityGroups is the security groups that are
                allowed SSH access (port 22) to the worker nodes
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
            subnets:
              description: Subnets is the VPC networks to use for the nodes
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
            tags:
              additionalProperties:
                type: string
              description: Tags are the AWS metadata to apply to the node group
              type: object
            version:
              description: Version is the Kubernetes version to run for the kubelet
              type: string
          required:
          - amiType
          - credentials
          - desiredSize
          - diskSize
          - eC2SSHKey
          - maxSize
          - minSize
          - region
          - subnets
          type: object
        status:
          description: EKSNodeGroupStatus defines the observed state of EKSNodeGroup
          properties:
            conditions:
              description: Conditions is the status of the components
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            nodeIAMRole:
              description: NodeIAMRole is the IAM role assumed by the worker nodes
                themselves
              type: string
            status:
              description: Status provides a overall status
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsAwsComputeKoreAppviaIo_eksnodegroupsYamlBytes() ([]byte, error) {
	return _crdsAwsComputeKoreAppviaIo_eksnodegroupsYaml, nil
}

func crdsAwsComputeKoreAppviaIo_eksnodegroupsYaml() (*asset, error) {
	bytes, err := crdsAwsComputeKoreAppviaIo_eksnodegroupsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/aws.compute.kore.appvia.io_eksnodegroups.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsAwsComputeKoreAppviaIo_eksvpcsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: eksvpcs.aws.compute.kore.appvia.io
spec:
  additionalPrinterColumns:
  - JSONPath: .status.status
    description: The overall status of the vpc
    name: Status
    type: string
  group: aws.compute.kore.appvia.io
  names:
    kind: EKSVPC
    listKind: EKSVPCList
    plural: eksvpcs
    singular: eksvpc
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: EKSVPC is the Schema for the eksvpc API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: EKSVPCSpec defines the desired state of EKSVPC
          properties:
            cluster:
              description: Cluster refers to the cluster this object belongs to
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            credentials:
              description: Credentials is a reference to an AWSCredentials object
                to use for authentication
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            privateIPV4Cidr:
              description: PrivateIPV4Cidr is the private range used for the VPC
              type: string
            region:
              description: Region is the AWS region of the VPC and any resources created
              type: string
          required:
          - credentials
          - privateIPV4Cidr
          - region
          type: object
        status:
          description: EKSVPCStatus defines the observed state of a VPC
          properties:
            conditions:
              description: Conditions is the status of the components
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            infra:
              description: Infra provides a cache of values discovered from infrastructure
                k8s:openapi-gen=false
              properties:
                ipv4EgressAddresses:
                  description: PublicIPV4EgressAddresses provides the source addresses
                    for traffic coming from the cluster - can provide input for securing
                    Kube API endpoints in managed clusters
                  items:
                    type: string
                  type: array
                privateSubnetIDs:
                  description: PrivateSubnetIds is a list of subnet IDs to use for
                    the worker nodes
                  items:
                    type: string
                  type: array
                  x-kubernetes-list-type: set
                publicSubnetIDs:
                  description: PublicSubnetIDs is a list of subnet IDs to use for
                    resources that need a public IP (e.g. load balancers)
                  items:
                    type: string
                  type: array
                  x-kubernetes-list-type: set
                securityGroupIDs:
                  description: SecurityGroupIds is a list of security group IDs to
                    use for a cluster
                  items:
                    type: string
                  type: array
                  x-kubernetes-list-type: set
              type: object
            status:
              description: Status provides a overall status
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsAwsComputeKoreAppviaIo_eksvpcsYamlBytes() ([]byte, error) {
	return _crdsAwsComputeKoreAppviaIo_eksvpcsYaml, nil
}

func crdsAwsComputeKoreAppviaIo_eksvpcsYaml() (*asset, error) {
	bytes, err := crdsAwsComputeKoreAppviaIo_eksvpcsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/aws.compute.kore.appvia.io_eksvpcs.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_clustersYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: clusters.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: Cluster
    listKind: ClusterList
    plural: clusters
    singular: cluster
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Cluster is the Schema for the plans API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ClusterSpec defines the desired state of a cluster
          properties:
            configuration:
              description: Configuration are the configuration values for this cluster
                It will contain values from the plan + overrides by the user This
                will provide a simple interface to calculate diffs between plan and
                cluster configuration
              type: object
              x-kubernetes-preserve-unknown-fields: true
            credentials:
              description: Credentials is a reference to the credentials object to
                use
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            kind:
              description: Kind refers to the cluster type (e.g. GKE, EKS)
              minLength: 1
              type: string
            plan:
              description: Plan is the name of the cluster plan which was used to
                create this cluster
              minLength: 1
              type: string
          required:
          - configuration
          - credentials
          - kind
          - plan
          type: object
        status:
          description: ClusterStatus defines the observed state of a cluster
          properties:
            apiEndpoint:
              description: APIEndpoint is the kubernetes API endpoint url
              type: string
            authProxyEndpoint:
              description: AuthProxyEndpoint is the endpoint of the authentication
                proxy for this cluster
              minLength: 1
              type: string
            caCertificate:
              description: CaCertificate is the base64 encoded cluster certificate
              type: string
            components:
              description: Components is a collection of component statuses
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            message:
              description: Message is the description of the current status
              type: string
            status:
              description: Status is the overall status of the cluster
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_clustersYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_clustersYaml, nil
}

func crdsClustersComputeKoreAppviaIo_clustersYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_clustersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_clusters.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_kubernetesYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: kubernetes.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: Kubernetes
    listKind: KubernetesList
    plural: kubernetes
    singular: kubernetes
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Kubernetes is the Schema for the roles API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: KubernetesSpec defines the desired state of Cluster
          properties:
            authProxyAllowedIPs:
              description: AuthProxyAllowedIPs is a list of IP address ranges (using
                CIDR format), which will be allowed to access the proxy
              items:
                type: string
              minItems: 1
              type: array
              x-kubernetes-list-type: set
            authProxyImage:
              description: AuthProxyImage is the kube api proxy used to sso into the
                cluster post provision
              type: string
            cluster:
              description: Cluster refers to the cluster this object belongs to
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            clusterUsers:
              description: ClusterUsers is a collection of users from the team whom
                have permissions across the cluster
              items:
                description: ClusterUser defines a user and their role in the cluster
                properties:
                  roles:
                    description: Roles is the roles the user is permitted access to
                    items:
                      type: string
                    minItems: 1
                    type: array
                    x-kubernetes-list-type: set
                  username:
                    description: Username is the team member the role is being applied
                      to
                    minLength: 1
                    type: string
                required:
                - roles
                - username
                type: object
              type: array
              x-kubernetes-list-type: set
            defaultTeamRole:
              description: DefaultTeamRole is role inherited by all team members
              type: string
            domain:
              description: Domain is the domain of the cluster
              type: string
            enableDefaultTrafficBlock:
              description: EnableDefaultTrafficBlock indicates the cluster should
                default to enabling blocking network policies on all namespaces
              type: boolean
            inheritTeamMembers:
              description: InheritTeamMembers inherits indicates all team members
                are inherited as having access to cluster by default.
              type: boolean
            provider:
              description: Provider is the cloud cluster provider type for this kubernetes
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
          type: object
        status:
          description: KubernetesStatus defines the observed state of Cluster
          properties:
            apiEndpoint:
              description: Endpoint is the kubernetes endpoint url
              type: string
            caCertificate:
              description: CaCertificate is the base64 encoded cluster certificate
              type: string
            components:
              description: Components is a collection of component statuses
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            endpoint:
              description: APIEndpoint is the endpoint of client proxy for this cluster
              minLength: 1
              type: string
            status:
              description: Status is overall status of the workspace
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_kubernetesYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_kubernetesYaml, nil
}

func crdsClustersComputeKoreAppviaIo_kubernetesYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_kubernetesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_kubernetes.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_managedclusterroleYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: managedclusterrole.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: ManagedClusterRole
    listKind: ManagedClusterRoleList
    plural: managedclusterrole
    singular: managedclusterrole
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Kubernetes is the Schema for the roles API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ManagedClusterRoleSpec defines the desired state of Cluster
            role
          properties:
            clusters:
              description: Clusters is used to apply to one of more clusters role
                to a specific cluster
              items:
                description: Ownership indicates the ownership of a resource
                properties:
                  group:
                    description: Group is the api group
                    type: string
                  kind:
                    description: Kind is the name of the resource under the group
                    type: string
                  name:
                    description: Name is name of the resource
                    type: string
                  namespace:
                    description: Namespace is the location of the object
                    type: string
                  version:
                    description: Version is the group version
                    type: string
                required:
                - group
                - kind
                - name
                - namespace
                - version
                type: object
              type: array
              x-kubernetes-list-type: set
            description:
              description: Description provides a short summary of the nature of the
                role
              minLength: 10
              type: string
            enabled:
              description: Enabled indicates if the role is enabled or not
              type: boolean
            rules:
              description: Rules are the permissions on the role
              items:
                description: PolicyRule holds information that describes a policy
                  rule, but does not contain information about who the rule applies
                  to or which namespace the rule applies to.
                properties:
                  apiGroups:
                    description: APIGroups is the name of the APIGroup that contains
                      the resources.  If multiple API groups are specified, any action
                      requested against one of the enumerated resources in any API
                      group will be allowed.
                    items:
                      type: string
                    type: array
                  nonResourceURLs:
                    description: NonResourceURLs is a set of partial urls that a user
                      should have access to.  *s are allowed, but only as the full,
                      final step in the path Since non-resource URLs are not namespaced,
                      this field is only applicable for ClusterRoles referenced from
                      a ClusterRoleBinding. Rules can either apply to API resources
                      (such as "pods" or "secrets") or non-resource URL paths (such
                      as "/api"),  but not both.
                    items:
                      type: string
                    type: array
                  resourceNames:
                    description: ResourceNames is an optional white list of names
                      that the rule applies to.  An empty set means that everything
                      is allowed.
                    items:
                      type: string
                    type: array
                  resources:
                    description: Resources is a list of resources this rule applies
                      to.  ResourceAll represents all resources.
                    items:
                      type: string
                    type: array
                  verbs:
                    description: Verbs is a list of Verbs that apply to ALL the ResourceKinds
                      and AttributeRestrictions contained in this rule.  VerbAll represents
                      all kinds.
                    items:
                      type: string
                    type: array
                required:
                - verbs
                type: object
              type: array
              x-kubernetes-list-type: set
            teams:
              description: Teams is used to filter the clusters to apply by team references
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
          type: object
        status:
          description: ManagedClusterRoleStatus defines the observed state of Cluster
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the workspace
              type: string
          required:
          - conditions
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_managedclusterroleYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_managedclusterroleYaml, nil
}

func crdsClustersComputeKoreAppviaIo_managedclusterroleYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_managedclusterroleYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_managedclusterrole.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_managedclusterrolebindingYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: managedclusterrolebinding.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: ManagedClusterRoleBinding
    listKind: ManagedClusterRoleBindingList
    plural: managedclusterrolebinding
    singular: managedclusterrolebinding
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Kubernetes is the Schema for the roles API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ManagedClusterRoleBindingSpec defines the desired state of
            Cluster role
          properties:
            binding:
              description: Binding is the cluster role binding you wish to propagate
                to the clusters
              properties:
                apiVersion:
                  description: 'APIVersion defines the versioned schema of this representation
                    of an object. Servers should convert recognized schemas to the
                    latest internal value, and may reject unrecognized values. More
                    info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                  type: string
                kind:
                  description: 'Kind is a string value representing the REST resource
                    this object represents. Servers may infer this from the endpoint
                    the client submits requests to. Cannot be updated. In CamelCase.
                    More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                  type: string
                metadata:
                  description: Standard object's metadata.
                  type: object
                roleRef:
                  description: RoleRef can only reference a ClusterRole in the global
                    namespace. If the RoleRef cannot be resolved, the Authorizer must
                    return an error.
                  properties:
                    apiGroup:
                      description: APIGroup is the group for the resource being referenced
                      type: string
                    kind:
                      description: Kind is the type of resource being referenced
                      type: string
                    name:
                      description: Name is the name of resource being referenced
                      type: string
                  required:
                  - apiGroup
                  - kind
                  - name
                  type: object
                subjects:
                  description: Subjects holds references to the objects the role applies
                    to.
                  items:
                    description: Subject contains a reference to the object or user
                      identities a role binding applies to.  This can either hold
                      a direct API object reference, or a value for non-objects such
                      as user and group names.
                    properties:
                      apiGroup:
                        description: APIGroup holds the API group of the referenced
                          subject. Defaults to "" for ServiceAccount subjects. Defaults
                          to "rbac.authorization.k8s.io" for User and Group subjects.
                        type: string
                      kind:
                        description: Kind of object being referenced. Values defined
                          by this API group are "User", "Group", and "ServiceAccount".
                          If the Authorizer does not recognized the kind value, the
                          Authorizer should report an error.
                        type: string
                      name:
                        description: Name of the object being referenced.
                        type: string
                      namespace:
                        description: Namespace of the referenced object.  If the object
                          kind is non-namespace, such as "User" or "Group", and this
                          value is not empty the Authorizer should report an error.
                        type: string
                    required:
                    - kind
                    - name
                    type: object
                  type: array
              required:
              - roleRef
              type: object
            clusters:
              description: Clusters is used to apply the cluster role to a specific
                cluster
              items:
                description: Ownership indicates the ownership of a resource
                properties:
                  group:
                    description: Group is the api group
                    type: string
                  kind:
                    description: Kind is the name of the resource under the group
                    type: string
                  name:
                    description: Name is name of the resource
                    type: string
                  namespace:
                    description: Namespace is the location of the object
                    type: string
                  version:
                    description: Version is the group version
                    type: string
                required:
                - group
                - kind
                - name
                - namespace
                - version
                type: object
              type: array
              x-kubernetes-list-type: set
            teams:
              description: Teams is a filter on the teams
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
          required:
          - binding
          type: object
        status:
          description: ManagedClusterRoleStatus defines the observed state of a cluster
            role binding
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the workspace
              type: string
          required:
          - conditions
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_managedclusterrolebindingYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_managedclusterrolebindingYaml, nil
}

func crdsClustersComputeKoreAppviaIo_managedclusterrolebindingYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_managedclusterrolebindingYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_managedclusterrolebinding.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_managedconfigYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: managedconfig.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: ManagedConfig
    listKind: ManagedConfigList
    plural: managedconfig
    singular: managedconfig
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: ManagedConfig is the Schema for the roles API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ManagedConfigSpec defines the configuration for a cluster
          properties:
            certificateAuthority:
              description: CertificateAuthority is the location of the API certificate
                authority
              properties:
                apiVersion:
                  description: 'APIVersion defines the versioned schema of this representation
                    of an object. Servers should convert recognized schemas to the
                    latest internal value, and may reject unrecognized values. More
                    info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                  type: string
                data:
                  additionalProperties:
                    format: byte
                    type: string
                  description: Data contains the secret data. Each key must consist
                    of alphanumeric characters, '-', '_' or '.'. The serialized form
                    of the secret data is a base64 encoded string, representing the
                    arbitrary (possibly non-string) data value here. Described in
                    https://tools.ietf.org/html/rfc4648#section-4
                  type: object
                immutable:
                  description: Immutable, if set to true, ensures that data stored
                    in the Secret cannot be updated (only object metadata can be modified).
                    If not set to true, the field can be modified at any time. Defaulted
                    to nil. This is an alpha field enabled by ImmutableEphemeralVolumes
                    feature gate.
                  type: boolean
                kind:
                  description: 'Kind is a string value representing the REST resource
                    this object represents. Servers may infer this from the endpoint
                    the client submits requests to. Cannot be updated. In CamelCase.
                    More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                  type: string
                metadata:
                  description: 'Standard object''s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata'
                  type: object
                stringData:
                  additionalProperties:
                    type: string
                  description: stringData allows specifying non-binary secret data
                    in string form. It is provided as a write-only convenience method.
                    All keys and values are merged into the data field on write, overwriting
                    any existing values. It is never output when reading from the
                    API.
                  type: object
                type:
                  description: Used to facilitate programmatic handling of secret
                    data.
                  type: string
              type: object
            clientCertificate:
              description: ClientCertificate is the location of the client certificate
                to speck back to the API
              properties:
                apiVersion:
                  description: 'APIVersion defines the versioned schema of this representation
                    of an object. Servers should convert recognized schemas to the
                    latest internal value, and may reject unrecognized values. More
                    info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                  type: string
                data:
                  additionalProperties:
                    format: byte
                    type: string
                  description: Data contains the secret data. Each key must consist
                    of alphanumeric characters, '-', '_' or '.'. The serialized form
                    of the secret data is a base64 encoded string, representing the
                    arbitrary (possibly non-string) data value here. Described in
                    https://tools.ietf.org/html/rfc4648#section-4
                  type: object
                immutable:
                  description: Immutable, if set to true, ensures that data stored
                    in the Secret cannot be updated (only object metadata can be modified).
                    If not set to true, the field can be modified at any time. Defaulted
                    to nil. This is an alpha field enabled by ImmutableEphemeralVolumes
                    feature gate.
                  type: boolean
                kind:
                  description: 'Kind is a string value representing the REST resource
                    this object represents. Servers may infer this from the endpoint
                    the client submits requests to. Cannot be updated. In CamelCase.
                    More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                  type: string
                metadata:
                  description: 'Standard object''s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata'
                  type: object
                stringData:
                  additionalProperties:
                    type: string
                  description: stringData allows specifying non-binary secret data
                    in string form. It is provided as a write-only convenience method.
                    All keys and values are merged into the data field on write, overwriting
                    any existing values. It is never output when reading from the
                    API.
                  type: object
                type:
                  description: Used to facilitate programmatic handling of secret
                    data.
                  type: string
              type: object
            domain:
              description: Domain is the domain name for this cluster
              minLength: 5
              type: string
          type: object
        status:
          description: ManagedConfigStatus defines the observed state of Cluster
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            phase:
              description: Phase indicates the phase of the cluster
              type: string
            status:
              description: Status is overall status of the workspace
              type: string
          required:
          - phase
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_managedconfigYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_managedconfigYaml, nil
}

func crdsClustersComputeKoreAppviaIo_managedconfigYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_managedconfigYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_managedconfig.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_managedpodsecuritypoliiesYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: managedpodsecuritypoliies.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: ManagedPodSecurityPolicy
    listKind: ManagedPodSecurityPolicyList
    plural: managedpodsecuritypoliies
    singular: managedpodsecuritypolicy
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Kubernetes is the Schema for the roles API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ManagedPodSecurityPolicySpec defines the desired state of Cluster
            role
          properties:
            clusters:
              description: Clusters is used to apply the cluster role to a specific
                cluster
              items:
                description: Ownership indicates the ownership of a resource
                properties:
                  group:
                    description: Group is the api group
                    type: string
                  kind:
                    description: Kind is the name of the resource under the group
                    type: string
                  name:
                    description: Name is name of the resource
                    type: string
                  namespace:
                    description: Namespace is the location of the object
                    type: string
                  version:
                    description: Version is the group version
                    type: string
                required:
                - group
                - kind
                - name
                - namespace
                - version
                type: object
              type: array
              x-kubernetes-list-type: set
            description:
              description: Description describes the nature of this pod security policy
              minLength: 1
              type: string
            policy:
              description: Policy defined a managed pod security policy across the
                clusters
              properties:
                allowPrivilegeEscalation:
                  description: allowPrivilegeEscalation determines if a pod can request
                    to allow privilege escalation. If unspecified, defaults to true.
                  type: boolean
                allowedCSIDrivers:
                  description: AllowedCSIDrivers is a whitelist of inline CSI drivers
                    that must be explicitly set to be embedded within a pod spec.
                    An empty value indicates that any CSI driver can be used for inline
                    ephemeral volumes. This is an alpha field, and is only honored
                    if the API server enables the CSIInlineVolume feature gate.
                  items:
                    description: AllowedCSIDriver represents a single inline CSI Driver
                      that is allowed to be used.
                    properties:
                      name:
                        description: Name is the registered name of the CSI driver
                        type: string
                    required:
                    - name
                    type: object
                  type: array
                allowedCapabilities:
                  description: allowedCapabilities is a list of capabilities that
                    can be requested to add to the container. Capabilities in this
                    field may be added at the pod author's discretion. You must not
                    list a capability in both allowedCapabilities and requiredDropCapabilities.
                  items:
                    description: Capability represent POSIX capabilities type
                    type: string
                  type: array
                allowedFlexVolumes:
                  description: allowedFlexVolumes is a whitelist of allowed Flexvolumes.  Empty
                    or nil indicates that all Flexvolumes may be used.  This parameter
                    is effective only when the usage of the Flexvolumes is allowed
                    in the "volumes" field.
                  items:
                    description: AllowedFlexVolume represents a single Flexvolume
                      that is allowed to be used.
                    properties:
                      driver:
                        description: driver is the name of the Flexvolume driver.
                        type: string
                    required:
                    - driver
                    type: object
                  type: array
                allowedHostPaths:
                  description: allowedHostPaths is a white list of allowed host paths.
                    Empty indicates that all host paths may be used.
                  items:
                    description: AllowedHostPath defines the host volume conditions
                      that will be enabled by a policy for pods to use. It requires
                      the path prefix to be defined.
                    properties:
                      pathPrefix:
                        description: "pathPrefix is the path prefix that the host
                          volume must match. It does not support ` + "`" + `*` + "`" + `. Trailing slashes
                          are trimmed when validating the path prefix with a host
                          path. \n Examples: ` + "`" + `/foo` + "`" + ` would allow ` + "`" + `/foo` + "`" + `, ` + "`" + `/foo/` + "`" + ` and
                          ` + "`" + `/foo/bar` + "`" + ` ` + "`" + `/foo` + "`" + ` would not allow ` + "`" + `/food` + "`" + ` or ` + "`" + `/etc/foo` + "`" + `"
                        type: string
                      readOnly:
                        description: when set to true, will allow host volumes matching
                          the pathPrefix only if all volume mounts are readOnly.
                        type: boolean
                    type: object
                  type: array
                allowedProcMountTypes:
                  description: AllowedProcMountTypes is a whitelist of allowed ProcMountTypes.
                    Empty or nil indicates that only the DefaultProcMountType may
                    be used. This requires the ProcMountType feature flag to be enabled.
                  items:
                    type: string
                  type: array
                allowedUnsafeSysctls:
                  description: "allowedUnsafeSysctls is a list of explicitly allowed
                    unsafe sysctls, defaults to none. Each entry is either a plain
                    sysctl name or ends in \"*\" in which case it is considered as
                    a prefix of allowed sysctls. Single * means all unsafe sysctls
                    are allowed. Kubelet has to whitelist all allowed unsafe sysctls
                    explicitly to avoid rejection. \n Examples: e.g. \"foo/*\" allows
                    \"foo/bar\", \"foo/baz\", etc. e.g. \"foo.*\" allows \"foo.bar\",
                    \"foo.baz\", etc."
                  items:
                    type: string
                  type: array
                defaultAddCapabilities:
                  description: defaultAddCapabilities is the default set of capabilities
                    that will be added to the container unless the pod spec specifically
                    drops the capability.  You may not list a capability in both defaultAddCapabilities
                    and requiredDropCapabilities. Capabilities added here are implicitly
                    allowed, and need not be included in the allowedCapabilities list.
                  items:
                    description: Capability represent POSIX capabilities type
                    type: string
                  type: array
                defaultAllowPrivilegeEscalation:
                  description: defaultAllowPrivilegeEscalation controls the default
                    setting for whether a process can gain more privileges than its
                    parent process.
                  type: boolean
                forbiddenSysctls:
                  description: "forbiddenSysctls is a list of explicitly forbidden
                    sysctls, defaults to none. Each entry is either a plain sysctl
                    name or ends in \"*\" in which case it is considered as a prefix
                    of forbidden sysctls. Single * means all sysctls are forbidden.
                    \n Examples: e.g. \"foo/*\" forbids \"foo/bar\", \"foo/baz\",
                    etc. e.g. \"foo.*\" forbids \"foo.bar\", \"foo.baz\", etc."
                  items:
                    type: string
                  type: array
                fsGroup:
                  description: fsGroup is the strategy that will dictate what fs group
                    is used by the SecurityContext.
                  properties:
                    ranges:
                      description: ranges are the allowed ranges of fs groups.  If
                        you would like to force a single fs group then supply a single
                        range with the same start and end. Required for MustRunAs.
                      items:
                        description: IDRange provides a min/max of an allowed range
                          of IDs.
                        properties:
                          max:
                            description: max is the end of the range, inclusive.
                            format: int64
                            type: integer
                          min:
                            description: min is the start of the range, inclusive.
                            format: int64
                            type: integer
                        required:
                        - max
                        - min
                        type: object
                      type: array
                    rule:
                      description: rule is the strategy that will dictate what FSGroup
                        is used in the SecurityContext.
                      type: string
                  type: object
                hostIPC:
                  description: hostIPC determines if the policy allows the use of
                    HostIPC in the pod spec.
                  type: boolean
                hostNetwork:
                  description: hostNetwork determines if the policy allows the use
                    of HostNetwork in the pod spec.
                  type: boolean
                hostPID:
                  description: hostPID determines if the policy allows the use of
                    HostPID in the pod spec.
                  type: boolean
                hostPorts:
                  description: hostPorts determines which host port ranges are allowed
                    to be exposed.
                  items:
                    description: HostPortRange defines a range of host ports that
                      will be enabled by a policy for pods to use.  It requires both
                      the start and end to be defined.
                    properties:
                      max:
                        description: max is the end of the range, inclusive.
                        format: int32
                        type: integer
                      min:
                        description: min is the start of the range, inclusive.
                        format: int32
                        type: integer
                    required:
                    - max
                    - min
                    type: object
                  type: array
                privileged:
                  description: privileged determines if a pod can request to be run
                    as privileged.
                  type: boolean
                readOnlyRootFilesystem:
                  description: readOnlyRootFilesystem when set to true will force
                    containers to run with a read only root file system.  If the container
                    specifically requests to run with a non-read only root file system
                    the PSP should deny the pod. If set to false the container may
                    run with a read only root file system if it wishes but it will
                    not be forced to.
                  type: boolean
                requiredDropCapabilities:
                  description: requiredDropCapabilities are the capabilities that
                    will be dropped from the container.  These are required to be
                    dropped and cannot be added.
                  items:
                    description: Capability represent POSIX capabilities type
                    type: string
                  type: array
                runAsGroup:
                  description: RunAsGroup is the strategy that will dictate the allowable
                    RunAsGroup values that may be set. If this field is omitted, the
                    pod's RunAsGroup can take any value. This field requires the RunAsGroup
                    feature gate to be enabled.
                  properties:
                    ranges:
                      description: ranges are the allowed ranges of gids that may
                        be used. If you would like to force a single gid then supply
                        a single range with the same start and end. Required for MustRunAs.
                      items:
                        description: IDRange provides a min/max of an allowed range
                          of IDs.
                        properties:
                          max:
                            description: max is the end of the range, inclusive.
                            format: int64
                            type: integer
                          min:
                            description: min is the start of the range, inclusive.
                            format: int64
                            type: integer
                        required:
                        - max
                        - min
                        type: object
                      type: array
                    rule:
                      description: rule is the strategy that will dictate the allowable
                        RunAsGroup values that may be set.
                      type: string
                  required:
                  - rule
                  type: object
                runAsUser:
                  description: runAsUser is the strategy that will dictate the allowable
                    RunAsUser values that may be set.
                  properties:
                    ranges:
                      description: ranges are the allowed ranges of uids that may
                        be used. If you would like to force a single uid then supply
                        a single range with the same start and end. Required for MustRunAs.
                      items:
                        description: IDRange provides a min/max of an allowed range
                          of IDs.
                        properties:
                          max:
                            description: max is the end of the range, inclusive.
                            format: int64
                            type: integer
                          min:
                            description: min is the start of the range, inclusive.
                            format: int64
                            type: integer
                        required:
                        - max
                        - min
                        type: object
                      type: array
                    rule:
                      description: rule is the strategy that will dictate the allowable
                        RunAsUser values that may be set.
                      type: string
                  required:
                  - rule
                  type: object
                runtimeClass:
                  description: runtimeClass is the strategy that will dictate the
                    allowable RuntimeClasses for a pod. If this field is omitted,
                    the pod's runtimeClassName field is unrestricted. Enforcement
                    of this field depends on the RuntimeClass feature gate being enabled.
                  properties:
                    allowedRuntimeClassNames:
                      description: allowedRuntimeClassNames is a whitelist of RuntimeClass
                        names that may be specified on a pod. A value of "*" means
                        that any RuntimeClass name is allowed, and must be the only
                        item in the list. An empty list requires the RuntimeClassName
                        field to be unset.
                      items:
                        type: string
                      type: array
                    defaultRuntimeClassName:
                      description: defaultRuntimeClassName is the default RuntimeClassName
                        to set on the pod. The default MUST be allowed by the allowedRuntimeClassNames
                        list. A value of nil does not mutate the Pod.
                      type: string
                  required:
                  - allowedRuntimeClassNames
                  type: object
                seLinux:
                  description: seLinux is the strategy that will dictate the allowable
                    labels that may be set.
                  properties:
                    rule:
                      description: rule is the strategy that will dictate the allowable
                        labels that may be set.
                      type: string
                    seLinuxOptions:
                      description: 'seLinuxOptions required to run as; required for
                        MustRunAs More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/'
                      properties:
                        level:
                          description: Level is SELinux level label that applies to
                            the container.
                          type: string
                        role:
                          description: Role is a SELinux role label that applies to
                            the container.
                          type: string
                        type:
                          description: Type is a SELinux type label that applies to
                            the container.
                          type: string
                        user:
                          description: User is a SELinux user label that applies to
                            the container.
                          type: string
                      type: object
                  required:
                  - rule
                  type: object
                supplementalGroups:
                  description: supplementalGroups is the strategy that will dictate
                    what supplemental groups are used by the SecurityContext.
                  properties:
                    ranges:
                      description: ranges are the allowed ranges of supplemental groups.  If
                        you would like to force a single supplemental group then supply
                        a single range with the same start and end. Required for MustRunAs.
                      items:
                        description: IDRange provides a min/max of an allowed range
                          of IDs.
                        properties:
                          max:
                            description: max is the end of the range, inclusive.
                            format: int64
                            type: integer
                          min:
                            description: min is the start of the range, inclusive.
                            format: int64
                            type: integer
                        required:
                        - max
                        - min
                        type: object
                      type: array
                    rule:
                      description: rule is the strategy that will dictate what supplemental
                        groups is used in the SecurityContext.
                      type: string
                  type: object
                volumes:
                  description: volumes is a white list of allowed volume plugins.
                    Empty indicates that no volumes may be used. To allow all volumes
                    you may use '*'.
                  items:
                    description: FSType gives strong typing to different file systems
                      that are used by volumes.
                    type: string
                  type: array
              required:
              - fsGroup
              - runAsUser
              - seLinux
              - supplementalGroups
              type: object
            teams:
              description: Teams is a filter on the teams
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
          type: object
        status:
          description: ManagedPodSecurityPolicyStatus defines the observed state of
            Cluster
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the workspace
              type: string
          required:
          - conditions
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_managedpodsecuritypoliiesYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_managedpodsecuritypoliiesYaml, nil
}

func crdsClustersComputeKoreAppviaIo_managedpodsecuritypoliiesYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_managedpodsecuritypoliiesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_managedpodsecuritypoliies.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_managedroleYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: managedrole.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: ManagedRole
    listKind: ManagedRoleList
    plural: managedrole
    singular: managedrole
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Kubernetes is the Schema for the roles API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ManagedRoleSpec defines the desired state of Cluster role
          properties:
            cluster:
              description: Cluster provides a link to the cluster which the role should
                reside
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            description:
              description: Description is a description for the role
              type: string
            role:
              description: Role are the permissions on the role
              items:
                description: PolicyRule holds information that describes a policy
                  rule, but does not contain information about who the rule applies
                  to or which namespace the rule applies to.
                properties:
                  apiGroups:
                    description: APIGroups is the name of the APIGroup that contains
                      the resources.  If multiple API groups are specified, any action
                      requested against one of the enumerated resources in any API
                      group will be allowed.
                    items:
                      type: string
                    type: array
                  nonResourceURLs:
                    description: NonResourceURLs is a set of partial urls that a user
                      should have access to.  *s are allowed, but only as the full,
                      final step in the path Since non-resource URLs are not namespaced,
                      this field is only applicable for ClusterRoles referenced from
                      a ClusterRoleBinding. Rules can either apply to API resources
                      (such as "pods" or "secrets") or non-resource URL paths (such
                      as "/api"),  but not both.
                    items:
                      type: string
                    type: array
                  resourceNames:
                    description: ResourceNames is an optional white list of names
                      that the rule applies to.  An empty set means that everything
                      is allowed.
                    items:
                      type: string
                    type: array
                  resources:
                    description: Resources is a list of resources this rule applies
                      to.  ResourceAll represents all resources.
                    items:
                      type: string
                    type: array
                  verbs:
                    description: Verbs is a list of Verbs that apply to ALL the ResourceKinds
                      and AttributeRestrictions contained in this rule.  VerbAll represents
                      all kinds.
                    items:
                      type: string
                    type: array
                required:
                - verbs
                type: object
              type: array
              x-kubernetes-list-type: set
          required:
          - description
          type: object
        status:
          description: ManagedRoleStatus defines the observed state of Cluster
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the workspace
              type: string
          required:
          - conditions
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_managedroleYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_managedroleYaml, nil
}

func crdsClustersComputeKoreAppviaIo_managedroleYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_managedroleYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_managedrole.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_namespaceclaimsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: namespaceclaims.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: NamespaceClaim
    listKind: NamespaceClaimList
    plural: namespaceclaims
    singular: namespaceclaim
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: NamespaceClaim is the Schema for the namespaceclaims API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: NamespaceClaimSpec defines the desired state of NamespaceClaim
          properties:
            annotations:
              additionalProperties:
                type: string
              description: Annotations is a series of annotations on the namespace
              type: object
            cluster:
              description: Cluster is the cluster the namespace resides
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            labels:
              additionalProperties:
                type: string
              description: Labels is a series of labels for the namespace
              type: object
            name:
              description: Name is the name of the namespace to create
              minLength: 1
              type: string
          required:
          - cluster
          - name
          type: object
        status:
          description: NamespaceClaimStatus defines the observed state of NamespaceClaim
          properties:
            conditions:
              description: Conditions is a series of things that caused the failure
                if any
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is the status of the namespace
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_namespaceclaimsYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_namespaceclaimsYaml, nil
}

func crdsClustersComputeKoreAppviaIo_namespaceclaimsYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_namespaceclaimsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_namespaceclaims.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsClustersComputeKoreAppviaIo_namespacepolicyYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: namespacepolicy.clusters.compute.kore.appvia.io
spec:
  group: clusters.compute.kore.appvia.io
  names:
    kind: NamepacePolicy
    listKind: NamepacePolicyList
    plural: namespacepolicy
    singular: namepacepolicy
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Kubernetes is the Schema for the roles API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: NamepacePolicySpec defines the desired state of Cluster role
          properties:
            defaultAnnotations:
              additionalProperties:
                type: string
              description: DefaultAnnotations are default annotations applied to all
                managed namespaces
              type: object
            defaultLabels:
              additionalProperties:
                type: string
              description: DefaultLabels are the labels applied to all managed namespaces
              type: object
            defaultLimits:
              description: DefaultLimits are the default resource limits applied to
                the namespace
              properties:
                apiVersion:
                  description: 'APIVersion defines the versioned schema of this representation
                    of an object. Servers should convert recognized schemas to the
                    latest internal value, and may reject unrecognized values. More
                    info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
                  type: string
                kind:
                  description: 'Kind is a string value representing the REST resource
                    this object represents. Servers may infer this from the endpoint
                    the client submits requests to. Cannot be updated. In CamelCase.
                    More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                  type: string
                metadata:
                  description: 'Standard object''s metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata'
                  type: object
                spec:
                  description: 'Spec defines the limits enforced. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status'
                  properties:
                    limits:
                      description: Limits is the list of LimitRangeItem objects that
                        are enforced.
                      items:
                        description: LimitRangeItem defines a min/max usage limit
                          for any resource that matches on kind.
                        properties:
                          default:
                            additionalProperties:
                              anyOf:
                              - type: integer
                              - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            description: Default resource requirement limit value
                              by resource name if resource limit is omitted.
                            type: object
                          defaultRequest:
                            additionalProperties:
                              anyOf:
                              - type: integer
                              - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            description: DefaultRequest is the default resource requirement
                              request value by resource name if resource request is
                              omitted.
                            type: object
                          max:
                            additionalProperties:
                              anyOf:
                              - type: integer
                              - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            description: Max usage constraints on this kind by resource
                              name.
                            type: object
                          maxLimitRequestRatio:
                            additionalProperties:
                              anyOf:
                              - type: integer
                              - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            description: MaxLimitRequestRatio if specified, the named
                              resource must have a request and limit that are both
                              non-zero where limit divided by request is less than
                              or equal to the enumerated value; this represents the
                              max burst for the named resource.
                            type: object
                          min:
                            additionalProperties:
                              anyOf:
                              - type: integer
                              - type: string
                              pattern: ^(\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\+|-)?(([0-9]+(\.[0-9]*)?)|(\.[0-9]+))))?$
                              x-kubernetes-int-or-string: true
                            description: Min usage constraints on this kind by resource
                              name.
                            type: object
                          type:
                            description: Type of resource that this limit applies
                              to.
                            type: string
                        required:
                        - type
                        type: object
                      type: array
                  required:
                  - limits
                  type: object
              type: object
          type: object
        status:
          description: NamepacePolicyStatus defines the observed state of Cluster
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the workspace
              type: string
          required:
          - conditions
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsClustersComputeKoreAppviaIo_namespacepolicyYamlBytes() ([]byte, error) {
	return _crdsClustersComputeKoreAppviaIo_namespacepolicyYaml, nil
}

func crdsClustersComputeKoreAppviaIo_namespacepolicyYaml() (*asset, error) {
	bytes, err := crdsClustersComputeKoreAppviaIo_namespacepolicyYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/clusters.compute.kore.appvia.io_namespacepolicy.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsConfigKoreAppviaIo_allocationsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: allocations.config.kore.appvia.io
spec:
  additionalPrinterColumns:
  - JSONPath: .spec.summary
    description: A summary of what is being shared
    name: Summary
    type: string
  - JSONPath: .spec.resource.group
    description: The API group of the resource being shared
    name: Group
    type: string
  - JSONPath: .spec.resource.namespace
    description: The namespace of the resource being shared
    name: Resource Namespace
    type: string
  - JSONPath: .spec.resource.name
    description: The name of the resource being shared
    name: Resource Name
    type: string
  group: config.kore.appvia.io
  names:
    kind: Allocation
    listKind: AllocationList
    plural: allocations
    singular: allocation
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Allocation is the Schema for the allocations API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: AllocationSpec defines the desired state of Allocation
          properties:
            name:
              description: Name is the name of the resource being shared
              type: string
            resource:
              description: Resource is the resource which is being shared with another
                team
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            summary:
              description: Summary is the summary of the resource being shared
              type: string
            teams:
              description: Teams is a collection of teams the allocation is permitted
                to use
              items:
                type: string
              minItems: 1
              type: array
              x-kubernetes-list-type: set
          required:
          - name
          - resource
          - summary
          - teams
          type: object
        status:
          description: AllocationStatus defines the observed state of Allocation
          properties:
            conditions:
              description: Conditions is a collection of potential issues
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is the general status of the resource
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsConfigKoreAppviaIo_allocationsYamlBytes() ([]byte, error) {
	return _crdsConfigKoreAppviaIo_allocationsYaml, nil
}

func crdsConfigKoreAppviaIo_allocationsYaml() (*asset, error) {
	bytes, err := crdsConfigKoreAppviaIo_allocationsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/config.kore.appvia.io_allocations.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsConfigKoreAppviaIo_planpoliciesYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: planpolicies.config.kore.appvia.io
spec:
  group: config.kore.appvia.io
  names:
    kind: PlanPolicy
    listKind: PlanPolicyList
    plural: planpolicies
    singular: planpolicy
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: PlanPolicy is the Schema for the plan policies API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: PlanPolicySpec defines Plan JSON Schema extensions
          properties:
            description:
              description: Description provides a detailed description of the plan
                policy
              type: string
            kind:
              description: Kind refers to the cluster type this is a plan policy for
              minLength: 1
              type: string
            labels:
              additionalProperties:
                type: string
              description: Labels is a collection of labels for this plan policy
              type: object
            properties:
              description: Properties are the
              items:
                description: PlanPolicyProperty defines a JSON schema for a given
                  property
                properties:
                  allowUpdate:
                    description: AllowUpdate will allow the parameter to be modified
                      by the teams
                    type: boolean
                  disallowUpdate:
                    description: DisallowUpdate will forbid modification of the parameter,
                      even if it was allowed by an other policy
                    type: boolean
                  name:
                    description: Name is the name of the property
                    minLength: 1
                    type: string
                required:
                - allowUpdate
                - disallowUpdate
                - name
                type: object
              minItems: 1
              type: array
              x-kubernetes-list-map-keys:
              - name
              x-kubernetes-list-type: map
            summary:
              description: Summary provides a short title summary for the plan policy
              minLength: 1
              type: string
          required:
          - kind
          - properties
          - summary
          type: object
        status:
          description: PlanPolicyStatus defines the observed state of Plan Policy
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the plan policy
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsConfigKoreAppviaIo_planpoliciesYamlBytes() ([]byte, error) {
	return _crdsConfigKoreAppviaIo_planpoliciesYaml, nil
}

func crdsConfigKoreAppviaIo_planpoliciesYaml() (*asset, error) {
	bytes, err := crdsConfigKoreAppviaIo_planpoliciesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/config.kore.appvia.io_planpolicies.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsConfigKoreAppviaIo_plansYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: plans.config.kore.appvia.io
spec:
  group: config.kore.appvia.io
  names:
    kind: Plan
    listKind: PlanList
    plural: plans
    singular: plan
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Plan is the Schema for the plans API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: PlanSpec defines the desired state of Plan
          properties:
            configuration:
              description: Configuration are the key+value pairs describing a cluster
                configuration
              type: object
              x-kubernetes-preserve-unknown-fields: true
            description:
              description: Description provides a summary of the configuration provided
                by this plan
              minLength: 1
              type: string
            kind:
              description: Resource refers to the resource type this is a plan for
              minLength: 1
              type: string
            labels:
              additionalProperties:
                type: string
              description: Labels is a collection of labels for this plan
              type: object
            summary:
              description: Summary provides a short title summary for the plan
              minLength: 1
              type: string
          required:
          - configuration
          - description
          - kind
          - summary
          type: object
        status:
          description: PlanStatus defines the observed state of Plan
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the workspace
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsConfigKoreAppviaIo_plansYamlBytes() ([]byte, error) {
	return _crdsConfigKoreAppviaIo_plansYaml, nil
}

func crdsConfigKoreAppviaIo_plansYaml() (*asset, error) {
	bytes, err := crdsConfigKoreAppviaIo_plansYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/config.kore.appvia.io_plans.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsConfigKoreAppviaIo_secretsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: secrets.config.kore.appvia.io
spec:
  group: config.kore.appvia.io
  names:
    kind: Secret
    listKind: SecretList
    plural: secrets
    singular: secret
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Secret is the Schema for the plans API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: SecretSpec defines the desired state of Plan
          properties:
            data:
              additionalProperties:
                type: string
              description: Values are the key values to the plan
              type: object
            description:
              description: Description provides a summary of the secret
              minLength: 1
              type: string
            type:
              description: Type refers to the secret type
              minLength: 1
              type: string
          required:
          - description
          - type
          type: object
        status:
          description: SecretStatus defines the observed state of Plan
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the workspace
              type: string
            systemManaged:
              description: SystemManaged indicates the secret is managed by kore and
                cannot be changed
              type: boolean
            verified:
              description: Verified indicates if the secret has been verified as working
              type: boolean
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsConfigKoreAppviaIo_secretsYamlBytes() ([]byte, error) {
	return _crdsConfigKoreAppviaIo_secretsYaml, nil
}

func crdsConfigKoreAppviaIo_secretsYaml() (*asset, error) {
	bytes, err := crdsConfigKoreAppviaIo_secretsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/config.kore.appvia.io_secrets.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsCoreKoreAppviaIo_idpYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: idp.core.kore.appvia.io
spec:
  group: core.kore.appvia.io
  names:
    kind: IDP
    listKind: IDPList
    plural: idp
    singular: idp
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: IDP is the Schema for the class API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: IDPSpec defines the spec for a configured instance of an IDP
          properties:
            config:
              description: IDPConfig
              properties:
                github:
                  description: Google represents a Google IDP config
                  properties:
                    clientID:
                      description: ClientID is the field name in a Github OAuth app
                      type: string
                    clientSecret:
                      description: ClientSecret is the field name in a Github OAuth
                        app
                      type: string
                    orgs:
                      description: Orgs is the list of possible Organisations in Github
                        the user must be part of
                      items:
                        type: string
                      type: array
                  required:
                  - clientID
                  - clientSecret
                  - orgs
                  type: object
                google:
                  description: GoogleIDP provides config for a Google Identity provider
                  properties:
                    clientID:
                      description: ClientID is the field name in a Google OAuth app
                      type: string
                    clientSecret:
                      description: ClientSecret is the field name in a Google OAuth
                        app
                      type: string
                    domains:
                      description: Domains are the google accounts whitelisted for
                        authentication
                      items:
                        type: string
                      type: array
                  required:
                  - clientID
                  - clientSecret
                  - domains
                  type: object
                oidc:
                  description: OIDCIDP config for a generoc Open ID Connect provider
                  properties:
                    clientID:
                      description: ClientID provides the OIDC client ID string
                      type: string
                    clientScopes:
                      description: ClientScopes provides the OIDC client scopes
                      items:
                        type: string
                      type: array
                    clientSecret:
                      description: ClientSecret provides the OIDC client secret string
                      type: string
                    issuer:
                      description: Issuer provides the IDP URL
                      type: string
                    userClaims:
                      description: UserClaims to track the identity field to use
                      items:
                        type: string
                      type: array
                  required:
                  - clientID
                  - clientScopes
                  - clientSecret
                  - issuer
                  - userClaims
                  type: object
                oidcdirect:
                  description: StaticOIDCIDP provides a means to detect when there
                    is no IDP broker It is essetially the same as a generic OIDC type
                  properties:
                    clientID:
                      description: ClientID provides the OIDC client ID string
                      type: string
                    clientScopes:
                      description: ClientScopes provides the OIDC client scopes
                      items:
                        type: string
                      type: array
                    clientSecret:
                      description: ClientSecret provides the OIDC client secret string
                      type: string
                    issuer:
                      description: Issuer provides the IDP URL
                      type: string
                    userClaims:
                      description: UserClaims to track the identity field to use
                      items:
                        type: string
                      type: array
                  required:
                  - clientID
                  - clientScopes
                  - clientSecret
                  - issuer
                  - userClaims
                  type: object
                saml:
                  description: SAMLIDP provides configuration for a generic SAML Identity
                    provider
                  properties:
                    allowedGroups:
                      description: AllowedGroups provides a list of allowed groups
                      items:
                        type: string
                      type: array
                    caData:
                      description: CAData is byte array representing the PEM data
                        for the IDP signing CA
                      format: byte
                      type: string
                    emailAttr:
                      description: EmailAttr attribute in the returned assertion to
                        map to ID token claims
                      type: string
                    groupsAttr:
                      description: GroupsAttr attribute in the returned assertion
                        to map to ID token claims
                      type: string
                    groupsDelim:
                      description: GroupsDelim characters used to split the single
                        groups field to obtain the user group membership
                      type: string
                    ssoURL:
                      description: SSOURL provides the SSO URL used for POST value
                        to IDP
                      type: string
                    usernameAttr:
                      description: UsernameAttr attribute in the returned assertion
                        to map to ID token claims
                      type: string
                  required:
                  - caData
                  - emailAttr
                  - ssoURL
                  - usernameAttr
                  type: object
              type: object
            displayName:
              description: DisplayName
              type: string
          required:
          - config
          - displayName
          type: object
        status:
          description: IDPStatus defines the observed state of an IDP (ID Providers)
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the IDP configuration
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsCoreKoreAppviaIo_idpYamlBytes() ([]byte, error) {
	return _crdsCoreKoreAppviaIo_idpYaml, nil
}

func crdsCoreKoreAppviaIo_idpYaml() (*asset, error) {
	bytes, err := crdsCoreKoreAppviaIo_idpYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/core.kore.appvia.io_idp.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsCoreKoreAppviaIo_oidclientYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: oidclient.core.kore.appvia.io
spec:
  group: core.kore.appvia.io
  names:
    kind: IDPClient
    listKind: IDPClientList
    plural: oidclient
    singular: idpclient
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: IDPClient is the Schema for the class API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: IDPClientSpec defines the spec for a IDP client
          properties:
            displayName:
              description: DisplayName
              type: string
            id:
              description: ID of OIDC client
              type: string
            redirectURIs:
              description: RedirectURIs where to send client after IDP auth
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
            secret:
              description: Secret for OIDC client
              type: string
          required:
          - displayName
          - id
          - redirectURIs
          - secret
          type: object
        status:
          description: IDPClientStatus defines the observed state of an IDP (ID Providers)
          properties:
            conditions:
              description: Conditions is a set of condition which has caused an error
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is overall status of the IDP configuration
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsCoreKoreAppviaIo_oidclientYamlBytes() ([]byte, error) {
	return _crdsCoreKoreAppviaIo_oidclientYaml, nil
}

func crdsCoreKoreAppviaIo_oidclientYaml() (*asset, error) {
	bytes, err := crdsCoreKoreAppviaIo_oidclientYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/core.kore.appvia.io_oidclient.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsGcpComputeKoreAppviaIo_organizationsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: organizations.gcp.compute.kore.appvia.io
spec:
  group: gcp.compute.kore.appvia.io
  names:
    kind: Organization
    listKind: OrganizationList
    plural: organizations
    singular: organization
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Organization is the Schema for the organization API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: OrganizationSpec defines the desired state of Organization
          properties:
            billingAccount:
              description: BillingAccountName is the resource name of the billing
                account associated with the project e.g. '012345-567890-ABCDEF'
              minLength: 1
              type: string
            credentialsRef:
              description: CredentialsRef is a reference to the credentials used to
                provision provision the projects - this is either created by dynamically
                from the oauth token or provided for us
              properties:
                name:
                  description: Name is unique within a namespace to reference a secret
                    resource.
                  type: string
                namespace:
                  description: Namespace defines the space within which the secret
                    name must be unique.
                  type: string
              type: object
            parentID:
              description: ParentID is the type specific ID of the parent this project
                has
              minLength: 1
              type: string
            parentType:
              description: 'ParentType is the type of parent this project has Valid
                types are: "organization", "folder", and "project"'
              enum:
              - organization
              - folder
              - project
              type: string
            serviceAccount:
              description: ServiceAccount is the name used when creating the service
                account e.g. 'hub-admin'
              minLength: 1
              type: string
            tokenRef:
              description: TokenRef is a reference to an ephemeral oauth token used
                provision the admin project
              properties:
                name:
                  description: Name is unique within a namespace to reference a secret
                    resource.
                  type: string
                namespace:
                  description: Namespace defines the space within which the secret
                    name must be unique.
                  type: string
              type: object
          required:
          - billingAccount
          - parentID
          - parentType
          - serviceAccount
          type: object
        status:
          description: OrganizationStatus defines the observed state of Organization
          properties:
            conditions:
              description: Conditions is a set of components conditions
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            projectID:
              description: Project is the GCP project ID
              type: string
            status:
              description: Status provides a overall status
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsGcpComputeKoreAppviaIo_organizationsYamlBytes() ([]byte, error) {
	return _crdsGcpComputeKoreAppviaIo_organizationsYaml, nil
}

func crdsGcpComputeKoreAppviaIo_organizationsYaml() (*asset, error) {
	bytes, err := crdsGcpComputeKoreAppviaIo_organizationsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/gcp.compute.kore.appvia.io_organizations.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsGcpComputeKoreAppviaIo_projectclaimsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: projectclaims.gcp.compute.kore.appvia.io
spec:
  group: gcp.compute.kore.appvia.io
  names:
    kind: ProjectClaim
    listKind: ProjectClaimList
    plural: projectclaims
    singular: projectclaim
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: ProjectClaim is the Schema for the ProjectClaims API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ProjectClaimSpec defines the desired state of ProjectClaim
          properties:
            organization:
              description: Organization isthe GCP organization
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            projectName:
              description: ProjectName is the name of the project to create
              type: string
          required:
          - organization
          - projectName
          type: object
        status:
          description: ProjectClaimStatus defines the observed state of GCP Project
          properties:
            conditions:
              description: Conditions is a set of components conditions
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            credentialRef:
              description: CredentialRef is the reference to the credentials secret
              properties:
                name:
                  description: Name is unique within a namespace to reference a secret
                    resource.
                  type: string
                namespace:
                  description: Namespace defines the space within which the secret
                    name must be unique.
                  type: string
              type: object
            projectID:
              description: ProjectID is the project id
              type: string
            projectRef:
              description: ProjectRef is a reference to the underlying project
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            status:
              description: Status provides a overall status
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsGcpComputeKoreAppviaIo_projectclaimsYamlBytes() ([]byte, error) {
	return _crdsGcpComputeKoreAppviaIo_projectclaimsYaml, nil
}

func crdsGcpComputeKoreAppviaIo_projectclaimsYaml() (*asset, error) {
	bytes, err := crdsGcpComputeKoreAppviaIo_projectclaimsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/gcp.compute.kore.appvia.io_projectclaims.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsGcpComputeKoreAppviaIo_projectsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: projects.gcp.compute.kore.appvia.io
spec:
  group: gcp.compute.kore.appvia.io
  names:
    kind: Project
    listKind: ProjectList
    plural: projects
    singular: project
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Project is the Schema for the ProjectClaims API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ProjectSpec defines the desired state of ProjectClaim
          properties:
            labels:
              additionalProperties:
                type: string
              description: Labels are a set of labels on the project
              type: object
            organization:
              description: Organization is a reference to the gcp admin project to
                use
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            projectName:
              description: ProjectName is the name of the project to create. We do
                this internally so we can easily change the project name without changing
                the resource name
              minLength: 1
              type: string
          required:
          - organization
          - projectName
          type: object
        status:
          description: ProjectStatus defines the observed state of GCP Project
          properties:
            conditions:
              description: Conditions is a set of components conditions
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            credentialRef:
              description: CredentialRef is the reference to the credentials secret
              properties:
                name:
                  description: Name is unique within a namespace to reference a secret
                    resource.
                  type: string
                namespace:
                  description: Namespace defines the space within which the secret
                    name must be unique.
                  type: string
              type: object
            projectID:
              description: ProjectID is the project id
              type: string
            status:
              description: Status provides a overall status
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsGcpComputeKoreAppviaIo_projectsYamlBytes() ([]byte, error) {
	return _crdsGcpComputeKoreAppviaIo_projectsYaml, nil
}

func crdsGcpComputeKoreAppviaIo_projectsYaml() (*asset, error) {
	bytes, err := crdsGcpComputeKoreAppviaIo_projectsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/gcp.compute.kore.appvia.io_projects.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsGkeComputeKoreAppviaIo_gkecredentialsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: gkecredentials.gke.compute.kore.appvia.io
spec:
  additionalPrinterColumns:
  - JSONPath: .spec.region
    description: The name of the GCP region the clusters will reside
    name: Region
    type: string
  - JSONPath: .spec.project
    description: The name of the GCP project
    name: Project
    type: string
  - JSONPath: .status.verified
    description: Indicates is the credentials have been verified
    name: Verified
    type: string
  group: gke.compute.kore.appvia.io
  names:
    kind: GKECredentials
    listKind: GKECredentialsList
    plural: gkecredentials
    singular: gkecredentials
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: GKECredentials is the Schema for the gkecredentials API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: GKECredentialsSpec defines the desired state of GKECredentials
          properties:
            account:
              description: Account is the credentials used to speak the GCP APIs;
                you create a service account under the Cloud IAM within the project,
                adding the permissions 'Compute Admin' role to the service account
                via IAM tab. Once done you can create a key under 'Service Accounts'
                and copy and paste the JSON payload here. This is deprecated, please
                use a Secret and CredentialsRef
              type: string
            credentialsRef:
              description: CredentialsRef is a reference to the credentials used to
                create clusters
              properties:
                name:
                  description: Name is unique within a namespace to reference a secret
                    resource.
                  type: string
                namespace:
                  description: Namespace defines the space within which the secret
                    name must be unique.
                  type: string
              type: object
            project:
              description: Project is the GCP project these credentias pretain to
              minLength: 1
              type: string
            region:
              description: Region is the GCP region you wish to the cluster to reside
                within
              type: string
          required:
          - project
          type: object
        status:
          description: GKECredentialsStatus defines the observed state of GKECredentials
          properties:
            conditions:
              description: Conditions is a collection of potential issues
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status provides a overall status
              type: string
            verified:
              description: Verified checks that the credentials are ok and valid
              type: boolean
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsGkeComputeKoreAppviaIo_gkecredentialsYamlBytes() ([]byte, error) {
	return _crdsGkeComputeKoreAppviaIo_gkecredentialsYaml, nil
}

func crdsGkeComputeKoreAppviaIo_gkecredentialsYaml() (*asset, error) {
	bytes, err := crdsGkeComputeKoreAppviaIo_gkecredentialsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/gke.compute.kore.appvia.io_gkecredentials.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsGkeComputeKoreAppviaIo_gkesYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: gkes.gke.compute.kore.appvia.io
spec:
  additionalPrinterColumns:
  - JSONPath: .spec.description
    description: A description of the GKE cluster
    name: Description
    type: string
  - JSONPath: .status.endpoint
    description: The endpoint of the gke cluster
    name: Endpoint
    type: string
  - JSONPath: .status.status
    description: The overall status of the cluster
    name: Status
    type: string
  group: gke.compute.kore.appvia.io
  names:
    kind: GKE
    listKind: GKEList
    plural: gkes
    singular: gke
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: GKE is the Schema for the gkes API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: GKESpec defines the desired state of GKE
          properties:
            authorizedMasterNetworks:
              description: AuthorizedMasterNetworks is a collection of authorized
                networks which is permitted to speak to the kubernetes API, default
                to all if not provided.
              items:
                description: AuthorizedNetwork provides a definition for the authorized
                  networks
                properties:
                  cidr:
                    description: CIDR is the network range associated to this network
                    type: string
                  name:
                    description: Name provides a descriptive name for this network
                    type: string
                required:
                - cidr
                - name
                type: object
              type: array
              x-kubernetes-list-type: set
            cluster:
              description: Cluster refers to the cluster this object belongs to
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            clusterIPV4Cidr:
              description: ClusterIPV4Cidr is an optional network CIDR which is used
                to place the pod network on
              type: string
            credentials:
              description: Credentials is a reference to the gke credentials object
                to use
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            description:
              description: Description provides a short summary / description of the
                cluster.
              minLength: 1
              type: string
            diskSize:
              description: DiskSize is the size of the disk used by the compute nodes.
              format: int64
              minimum: 100
              type: integer
            enableAutorepair:
              description: EnableAutorepair indicates if the cluster should be configured
                with auto repair is enabled
              type: boolean
            enableAutoscaler:
              description: EnableAutoscaler indicates if the cluster should be configured
                with cluster autoscaling turned on
              type: boolean
            enableAutoupgrade:
              description: EnableAutoUpgrade indicates if the cluster should be configured
                with autograding enabled; meaning both nodes are masters are autoscaled
                scheduled to upgrade during your maintenance window.
              type: boolean
            enableHTTPLoadBalancer:
              description: EnableHTTPLoadBalancer indicates if the cluster should
                be configured with the GKE ingress controller. When enabled GKE will
                autodiscover your ingress resources and provision load balancer on
                your behalf.
              type: boolean
            enableHorizontalPodAutoscaler:
              description: EnableHorizontalPodAutoscaler indicates if the cluster
                is configured with the horizontal pod autoscaler addon. This automatically
                adjusts the cpu and memory resources of pods in accordance with their
                demand. You should ensure you use PodDisruptionBudgets if this is
                enabled.
              type: boolean
            enableIstio:
              description: EnableIstio indicates if the GKE Istio service mesh is
                deployed to the cluster; this provides a more feature rich routing
                and instrumentation.
              type: boolean
            enablePrivateEndpoint:
              description: EnablePrivateEndpoint indicates whether the Kubernetes
                API should only be accessible from internal IP addresses
              type: boolean
            enablePrivateNetwork:
              description: EnablePrivateNetwork indicates if compute nodes should
                have external ip addresses or use private networking and a cloud-nat
                device.
              type: boolean
            enableShieldedNodes:
              description: EnableShieldedNodes indicates we should enable the shielded
                nodes options in GKE. This protects against a variety of attacks by
                hardening the underlying GKE node against rootkits and bootkits.
              type: boolean
            enableStackDriverLogging:
              description: EnableStackDriverLogging indicates if Stackdriver logging
                should be enabled for the cluster
              type: boolean
            enableStackDriverMetrics:
              description: EnableStackDriverMetrics indicates if Stackdriver metrics
                should be enabled for the cluster
              type: boolean
            imageType:
              description: ImageType is the operating image to use for the default
                compute pool.
              minLength: 1
              type: string
            machineType:
              description: MachineType is the machine type which the default nodes
                pool should use.
              minLength: 1
              type: string
            maintenanceWindow:
              description: MaintenanceWindow is the maintenance window provided for
                GKE to perform upgrades if enabled.
              minLength: 1
              type: string
            masterIPV4Cidr:
              description: MasterIPV4Cidr is network range used when private networking
                is enabled. This is the peering subnet used to to GKE master api layer.
                Note, this must be unique within the network.
              type: string
            maxSize:
              description: MaxSize assuming the autoscaler is enabled this is the
                maximum number nodes permitted
              format: int64
              minimum: 2
              type: integer
            network:
              description: Network is the GCP network the cluster reside on, which
                have to be unique within the GCP project and created beforehand.
              minLength: 1
              type: string
            region:
              description: Region is the gcp region you want the cluster to reside
              minLength: 1
              type: string
            servicesIPV4Cidr:
              description: ServicesIPV4Cidr is an optional network cidr configured
                for the cluster services
              type: string
            size:
              description: Size is the number of nodes per zone which should exist
                in the cluster.
              format: int64
              minimum: 1
              type: integer
            subnetwork:
              description: Subnetwork is name of the GCP subnetwork which the cluster
                nodes should reside
              minLength: 1
              type: string
            tags:
              additionalProperties:
                type: string
              description: Tags is a collection of tags related to the cluster type
              type: object
            version:
              description: Version is the initial kubernetes version which the cluster
                should be configured with.
              minLength: 1
              type: string
          required:
          - credentials
          - description
          - diskSize
          - enableShieldedNodes
          - imageType
          - machineType
          - maintenanceWindow
          - maxSize
          - network
          - size
          - subnetwork
          - version
          type: object
        status:
          description: GKEStatus defines the observed state of GKE
          properties:
            caCertificate:
              description: CACertificate is the certificate for this cluster
              type: string
            conditions:
              description: Conditions is the status of the components
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            endpoint:
              description: Endpoint is the endpoint of the cluster
              type: string
            status:
              description: Status provides a overall status
              type: string
          type: object
      type: object
  version: v1alpha1
  versions:
  - name: v1alpha1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsGkeComputeKoreAppviaIo_gkesYamlBytes() ([]byte, error) {
	return _crdsGkeComputeKoreAppviaIo_gkesYaml, nil
}

func crdsGkeComputeKoreAppviaIo_gkesYaml() (*asset, error) {
	bytes, err := crdsGkeComputeKoreAppviaIo_gkesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/gke.compute.kore.appvia.io_gkes.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsOrgKoreAppviaIo_auditeventsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: auditevents.org.kore.appvia.io
spec:
  group: org.kore.appvia.io
  names:
    kind: AuditEvent
    listKind: AuditEventList
    plural: auditevents
    singular: auditevent
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: AuditEvent is the Schema for the audit API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: AuditEventSpec defines the desired state of User
          properties:
            apiVersion:
              description: APIVersion is the version of the API used for this operation.
              type: string
            completedAt:
              description: CompletedAt is the timestamp the operation completed
              format: date-time
              type: string
            createdAt:
              description: CreatedAt is the timestamp of record creation
              format: date-time
              type: string
            id:
              description: ID is the unique identifier of this audit event.
              type: integer
            message:
              description: Message is event message itself
              type: string
            operation:
              description: Operation is the operation performed (e.g. UpdateCluster,
                CreateCluster, etc).
              type: string
            resource:
              description: Resource is the area of the API accessed in this audit
                operation (e.g. teams, ).
              type: string
            resourceURI:
              description: ResourceURI is the identifier of the resource in question.
              type: string
            responseCode:
              description: ResponseCode indicates the HTTP status code of the operation
                (e.g. 200, 404, etc).
              type: integer
            startedAt:
              description: StartedAt is the timestamp the operation was initiated
              format: date-time
              type: string
            team:
              description: Team is the team whom event may be associated to
              type: string
            user:
              description: User is the user which the event is related
              type: string
            verb:
              description: Verb is the type of action performed (e.g. PUT, GET, etc)
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsOrgKoreAppviaIo_auditeventsYamlBytes() ([]byte, error) {
	return _crdsOrgKoreAppviaIo_auditeventsYaml, nil
}

func crdsOrgKoreAppviaIo_auditeventsYaml() (*asset, error) {
	bytes, err := crdsOrgKoreAppviaIo_auditeventsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/org.kore.appvia.io_auditevents.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsOrgKoreAppviaIo_membersYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: members.org.kore.appvia.io
spec:
  group: org.kore.appvia.io
  names:
    kind: TeamMember
    listKind: TeamMemberList
    plural: members
    singular: teammember
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: TeamMember is the Schema for the teams API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: TeamMemberSpec defines the desired state of Team
          properties:
            roles:
              description: Role is the role of the user in the team
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
            team:
              description: Team is the name of the team
              type: string
            username:
              description: Username is the user being bound to the team
              type: string
          required:
          - roles
          - team
          - username
          type: object
        status:
          description: TeamMemberStatus defines the observed state of Team
          properties:
            conditions:
              description: Conditions is a collection of possible errors
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is the status of the resource
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsOrgKoreAppviaIo_membersYamlBytes() ([]byte, error) {
	return _crdsOrgKoreAppviaIo_membersYaml, nil
}

func crdsOrgKoreAppviaIo_membersYaml() (*asset, error) {
	bytes, err := crdsOrgKoreAppviaIo_membersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/org.kore.appvia.io_members.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsOrgKoreAppviaIo_teaminvitationsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: teaminvitations.org.kore.appvia.io
spec:
  group: org.kore.appvia.io
  names:
    kind: TeamInvitation
    listKind: TeamInvitationList
    plural: teaminvitations
    singular: teaminvitation
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: TeamInvitation is the Schema for the teams API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: TeamInvitationSpec defines the desired state of Team
          properties:
            team:
              description: Team is the name of the team
              type: string
            username:
              description: Username is the user being bound to the team
              type: string
          required:
          - team
          - username
          type: object
        status:
          description: TeamInvitationStatus defines the observed state of Team
          properties:
            conditions:
              description: Conditions is a collection of possible errors
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is the status of the resource
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsOrgKoreAppviaIo_teaminvitationsYamlBytes() ([]byte, error) {
	return _crdsOrgKoreAppviaIo_teaminvitationsYaml, nil
}

func crdsOrgKoreAppviaIo_teaminvitationsYaml() (*asset, error) {
	bytes, err := crdsOrgKoreAppviaIo_teaminvitationsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/org.kore.appvia.io_teaminvitations.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsOrgKoreAppviaIo_teamsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: teams.org.kore.appvia.io
spec:
  group: org.kore.appvia.io
  names:
    kind: Team
    listKind: TeamList
    plural: teams
    singular: team
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Team is the Schema for the teams API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: TeamSpec defines the desired state of Team
          properties:
            description:
              description: Description is a description for the team
              type: string
            summary:
              description: Summary is a summary name for this team
              type: string
          required:
          - description
          - summary
          type: object
        status:
          description: TeamStatus defines the observed state of Team
          properties:
            conditions:
              description: Conditions is a collection of possible errors
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status is the status of the resource
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsOrgKoreAppviaIo_teamsYamlBytes() ([]byte, error) {
	return _crdsOrgKoreAppviaIo_teamsYaml, nil
}

func crdsOrgKoreAppviaIo_teamsYaml() (*asset, error) {
	bytes, err := crdsOrgKoreAppviaIo_teamsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/org.kore.appvia.io_teams.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsOrgKoreAppviaIo_usersYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: users.org.kore.appvia.io
spec:
  group: org.kore.appvia.io
  names:
    kind: User
    listKind: UserList
    plural: users
    singular: user
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: User is the Schema for the users API
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: UserSpec defines the desired state of User
          properties:
            disabled:
              description: Disabled indicates if the user is disabled
              type: boolean
            email:
              description: Email is the email for the user
              minLength: 1
              type: string
            username:
              description: Username is the userame or identity for this user
              minLength: 1
              type: string
          required:
          - disabled
          - email
          - username
          type: object
        status:
          description: UserStatus defines the observed state of User
          properties:
            conditions:
              description: Conditions is collection of potentials error causes
              items:
                description: Condition is a reason why something failed
                properties:
                  detail:
                    description: Detail is a actual error which might contain technical
                      reference
                    type: string
                  message:
                    description: Message is a human readable message
                    type: string
                required:
                - detail
                - message
                type: object
              type: array
              x-kubernetes-list-type: set
            status:
              description: Status provides an overview of the user status
              type: string
          required:
          - status
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsOrgKoreAppviaIo_usersYamlBytes() ([]byte, error) {
	return _crdsOrgKoreAppviaIo_usersYaml, nil
}

func crdsOrgKoreAppviaIo_usersYaml() (*asset, error) {
	bytes, err := crdsOrgKoreAppviaIo_usersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/org.kore.appvia.io_users.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsSecurityKoreAppviaIo_securityoverviewsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: securityoverviews.security.kore.appvia.io
spec:
  group: security.kore.appvia.io
  names:
    kind: SecurityOverview
    listKind: SecurityOverviewList
    plural: securityoverviews
    singular: securityoverview
  preserveUnknownFields: false
  scope: Namespaced
  validation:
    openAPIV3Schema:
      description: SecurityOverview contains a report about the current state of Kore
        or a team
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: SecurityOverviewSpec shows the overall current security posture
            of Kore or a team
          properties:
            openIssueCounts:
              additionalProperties:
                format: int64
                type: integer
              description: OpenIssueCounts informs how many issues of each rule status
                exist currently
              type: object
            resources:
              description: Resources contains summaries of the open issues for each
                resource
              items:
                description: SecurityResourceOverview provides an overview of the
                  open issue counts for a resource
                properties:
                  lastChecked:
                    description: LastChecked is the timestamp this resource was last
                      scanned
                    format: date-time
                    type: string
                  openIssueCounts:
                    additionalProperties:
                      format: int64
                      type: integer
                    description: OpenIssueCounts is the summary of open issues for
                      this resource
                    type: object
                  overallStatus:
                    description: OverallStatus is the overall status of this resource
                    type: string
                  resource:
                    description: Resource is a reference to the group/version/kind/namespace/name
                      of the resource scanned by this scan
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                type: object
              type: array
            team:
              description: Team will be populated with the team name if this report
                is about a team, else unpopulated for a report for the whole of Kore
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsSecurityKoreAppviaIo_securityoverviewsYamlBytes() ([]byte, error) {
	return _crdsSecurityKoreAppviaIo_securityoverviewsYaml, nil
}

func crdsSecurityKoreAppviaIo_securityoverviewsYaml() (*asset, error) {
	bytes, err := crdsSecurityKoreAppviaIo_securityoverviewsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/security.kore.appvia.io_securityoverviews.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsSecurityKoreAppviaIo_securityrulesYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: securityrules.security.kore.appvia.io
spec:
  group: security.kore.appvia.io
  names:
    kind: SecurityRule
    listKind: SecurityRuleList
    plural: securityrules
    singular: securityrule
  preserveUnknownFields: false
  scope: Namespaced
  validation:
    openAPIV3Schema:
      description: SecurityRule contains the definition of a security rule
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: SecurityRuleSpec specifies the details of a security rule
          properties:
            appliesTo:
              description: AppliesTo is the list of resource types (e.g. Plan, Cluster)
                that this rule is applicable for
              items:
                type: string
              type: array
            code:
              description: Code is the unique identifier of this rule
              type: string
            description:
              description: Description is the markdown-formatted extended description
                of this rule.
              type: string
            name:
              description: Name is the human-readable name of this rule
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsSecurityKoreAppviaIo_securityrulesYamlBytes() ([]byte, error) {
	return _crdsSecurityKoreAppviaIo_securityrulesYaml, nil
}

func crdsSecurityKoreAppviaIo_securityrulesYaml() (*asset, error) {
	bytes, err := crdsSecurityKoreAppviaIo_securityrulesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/security.kore.appvia.io_securityrules.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsSecurityKoreAppviaIo_securityscanresultsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: securityscanresults.security.kore.appvia.io
spec:
  group: security.kore.appvia.io
  names:
    kind: SecurityScanResult
    listKind: SecurityScanResultList
    plural: securityscanresults
    singular: securityscanresult
  preserveUnknownFields: false
  scope: Namespaced
  validation:
    openAPIV3Schema:
      description: SecurityScanResult contains the result of a scan against all registered
        rules
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: SecurityScanResultSpec shows the overall result of a scan against
            all registered rules
          properties:
            archivedAt:
              description: ArchivedAt is the timestamp this result was superceded
                by a later scan - if ArchivedAt.IsZero() is true this is the most
                recent scan.
              format: date-time
              type: string
            checkedAt:
              description: CheckedAt is the timestamp this result was determined
              format: date-time
              type: string
            id:
              description: ID is the ID of this scan result in the data store
              format: int64
              type: integer
            overallStatus:
              description: OverallStatus indicates the worst-case status of the rules
                checked in this scan
              type: string
            owningTeam:
              description: OwningTeam is the name of the Kore team that owns this
                resource, will be empty if it is a non-team resource.
              type: string
            resource:
              description: Resource is a reference to the group/version/kind/namespace/name
                of the resource scanned by this scan
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            results:
              description: Results are the underlying results of the individual rules
                run as part of this scan
              items:
                description: SecurityScanRuleResult represents the compliance status
                  of a target with respect to a specific security rule.
                properties:
                  checkedAt:
                    description: CheckedAt is the timestamp this result was determined
                    format: date-time
                    type: string
                  message:
                    description: Message provides additional information about the
                      status of this rule on this target, if applicable
                    type: string
                  ruleCode:
                    description: RuleCode indicates the rule that this result relates
                      to
                    type: string
                  status:
                    description: Status indicates the compliance of the target with
                      this rule
                    type: string
                type: object
              type: array
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsSecurityKoreAppviaIo_securityscanresultsYamlBytes() ([]byte, error) {
	return _crdsSecurityKoreAppviaIo_securityscanresultsYaml, nil
}

func crdsSecurityKoreAppviaIo_securityscanresultsYaml() (*asset, error) {
	bytes, err := crdsSecurityKoreAppviaIo_securityscanresultsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/security.kore.appvia.io_securityscanresults.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsServicesKoreAppviaIo_servicecredentialsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: servicecredentials.services.kore.appvia.io
spec:
  group: services.kore.appvia.io
  names:
    kind: ServiceCredentials
    listKind: ServiceCredentialsList
    plural: servicecredentials
    singular: servicecredentials
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: ServiceCredentials is credentials provisioned by a service into
        the target namespace
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ServiceCredentialsSpec defines the the desired status for service
            credentials
          properties:
            cluster:
              description: Cluster contains the reference to the cluster where the
                credentials will be saved as a secret
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            clusterNamespace:
              description: ClusterNamespace is the target namespace in the cluster
                where the secret will be created
              type: string
            configuration:
              description: Configuration are the configuration values for this service
                credentials It will be used by the service provider to provision the
                credentials
              type: object
              x-kubernetes-preserve-unknown-fields: true
            kind:
              description: Kind refers to the service type
              minLength: 1
              type: string
            secretName:
              description: SecretName is the Kubernetes Secret's name that will contain
                the service access information If not set the secret's name will default
                to ` + "`" + `Name` + "`" + `
              type: string
            service:
              description: Service contains the reference to the service object
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
          required:
          - kind
          type: object
        status:
          description: ServiceCredentialsStatus defines the observed state of a service
          properties:
            components:
              description: Components is a collection of component statuses
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            message:
              description: Message is the description of the current status
              type: string
            providerData:
              description: ProviderData is provider specific data
              type: object
              x-kubernetes-preserve-unknown-fields: true
            providerID:
              description: ProviderID is the service credentials identifier in the
                service provider
              type: string
            status:
              description: Status is the overall status of the service
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsServicesKoreAppviaIo_servicecredentialsYamlBytes() ([]byte, error) {
	return _crdsServicesKoreAppviaIo_servicecredentialsYaml, nil
}

func crdsServicesKoreAppviaIo_servicecredentialsYaml() (*asset, error) {
	bytes, err := crdsServicesKoreAppviaIo_servicecredentialsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/services.kore.appvia.io_servicecredentials.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsServicesKoreAppviaIo_servicekindsYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: servicekinds.services.kore.appvia.io
spec:
  group: services.kore.appvia.io
  names:
    kind: ServiceKind
    listKind: ServiceKindList
    plural: servicekinds
    singular: servicekind
  preserveUnknownFields: false
  scope: Namespaced
  validation:
    openAPIV3Schema:
      description: ServiceKind is a service type
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ServiceKindSpec defines the state of a service kind
          properties:
            credentialSchema:
              description: CredentialSchema is the JSON schema for credentials created
                for service using this plan
              type: string
            description:
              description: Description is a detailed description of the service kind
              type: string
            displayName:
              description: DisplayName refers to the display name of the service type
              minLength: 1
              type: string
            documentationURL:
              description: DocumentationURL refers to the documentation page for this
                service
              type: string
            enabled:
              description: Enabled is true if the service kind can be used
              type: boolean
            imageURL:
              description: ImageURL is a thumbnail for the service kind
              type: string
            providerData:
              description: ProviderData is provider specific data
              type: object
              x-kubernetes-preserve-unknown-fields: true
            schema:
              description: Schema is the JSON schema for the plan
              type: string
            summary:
              description: Summary provides a short title summary for the service
                kind
              minLength: 1
              type: string
          required:
          - summary
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsServicesKoreAppviaIo_servicekindsYamlBytes() ([]byte, error) {
	return _crdsServicesKoreAppviaIo_servicekindsYaml, nil
}

func crdsServicesKoreAppviaIo_servicekindsYaml() (*asset, error) {
	bytes, err := crdsServicesKoreAppviaIo_servicekindsYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/services.kore.appvia.io_servicekinds.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsServicesKoreAppviaIo_serviceplansYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: serviceplans.services.kore.appvia.io
spec:
  group: services.kore.appvia.io
  names:
    kind: ServicePlan
    listKind: ServicePlanList
    plural: serviceplans
    singular: serviceplan
  preserveUnknownFields: false
  scope: Namespaced
  validation:
    openAPIV3Schema:
      description: ServicePlan is a template for a service
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ServicePlanSpec defines the desired state of Service plan
          properties:
            configuration:
              description: Configuration are the key+value pairs describing a service
                configuration
              type: object
              x-kubernetes-preserve-unknown-fields: true
            credentialSchema:
              description: CredentialSchema is the JSON schema for credentials created
                for service using this plan
              type: string
            description:
              description: Description is a detailed description of the service plan
              type: string
            displayName:
              description: DisplayName refers to the display name of the service type
              minLength: 1
              type: string
            kind:
              description: Kind refers to the service type this is a plan for
              minLength: 1
              type: string
            labels:
              additionalProperties:
                type: string
              description: Labels is a collection of labels for this plan
              type: object
            providerData:
              description: ProviderData is provider specific data
              type: object
              x-kubernetes-preserve-unknown-fields: true
            schema:
              description: Schema is the JSON schema for the plan
              type: string
            summary:
              description: Summary provides a short title summary for the plan
              minLength: 1
              type: string
          required:
          - kind
          - summary
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsServicesKoreAppviaIo_serviceplansYamlBytes() ([]byte, error) {
	return _crdsServicesKoreAppviaIo_serviceplansYaml, nil
}

func crdsServicesKoreAppviaIo_serviceplansYaml() (*asset, error) {
	bytes, err := crdsServicesKoreAppviaIo_serviceplansYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/services.kore.appvia.io_serviceplans.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsServicesKoreAppviaIo_serviceprovidersYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: serviceproviders.services.kore.appvia.io
spec:
  group: services.kore.appvia.io
  names:
    kind: ServiceProvider
    listKind: ServiceProviderList
    plural: serviceproviders
    singular: serviceprovider
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: ServiceProvider is a template for a service provider
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ServiceProviderSpec defines the desired state of a Service
            provider
          properties:
            configuration:
              description: Configuration are the key+value pairs describing a service
                provider
              type: object
              x-kubernetes-preserve-unknown-fields: true
            configurationFrom:
              description: ConfigurationFrom is a way to load configuration values
                from alternative sources, e.g. from secrets The values from these
                sources will override any existing keys defined in Configuration
              items:
                properties:
                  name:
                    description: Name is the name of the configuration parameter
                    minLength: 1
                    type: string
                  secretKeyRef:
                    description: SecretKeyRef is a reference to a key in a secret
                    properties:
                      key:
                        description: Key is they data key in the secret
                        minLength: 1
                        type: string
                      name:
                        description: Name is the name of the secret
                        minLength: 1
                        type: string
                      optional:
                        description: Optional controls whether the secret with the
                          given key must exist
                        type: boolean
                    required:
                    - name
                    type: object
                required:
                - name
                - secretKeyRef
                type: object
              type: array
              x-kubernetes-list-type: set
            credentials:
              description: Credentials is a reference to the credentials object to
                use
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            description:
              description: Description is a detailed description of the service provider
              type: string
            summary:
              description: Summary provides a short title summary for the provider
              minLength: 1
              type: string
            type:
              description: Type refers to the service provider type
              minLength: 1
              type: string
          required:
          - summary
          - type
          type: object
        status:
          description: ServiceProviderStatus defines the observed state of a service
            provider
          properties:
            components:
              description: Components is a collection of component statuses
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            message:
              description: Message is the description of the current status
              type: string
            status:
              description: Status is the overall status of the service
              type: string
            supportedKinds:
              description: SupportedKinds contains all the supported service kinds
              items:
                type: string
              type: array
              x-kubernetes-list-type: set
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsServicesKoreAppviaIo_serviceprovidersYamlBytes() ([]byte, error) {
	return _crdsServicesKoreAppviaIo_serviceprovidersYaml, nil
}

func crdsServicesKoreAppviaIo_serviceprovidersYaml() (*asset, error) {
	bytes, err := crdsServicesKoreAppviaIo_serviceprovidersYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/services.kore.appvia.io_serviceproviders.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _crdsServicesKoreAppviaIo_servicesYaml = []byte(`
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.2.5
  creationTimestamp: null
  name: services.services.kore.appvia.io
spec:
  group: services.kore.appvia.io
  names:
    kind: Service
    listKind: ServiceList
    plural: services
    singular: service
  preserveUnknownFields: false
  scope: Namespaced
  subresources:
    status: {}
  validation:
    openAPIV3Schema:
      description: Service is a managed service instance
      properties:
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: ServiceSpec defines the desired state of a service
          properties:
            cluster:
              description: Cluster contains the reference to the cluster where the
                service will be created
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            clusterNamespace:
              description: ClusterNamespace is the target namespace in the cluster
                where there the service will be created
              type: string
            configuration:
              description: Configuration are the configuration values for this service
                It will contain values from the plan + overrides by the user This
                will provide a simple interface to calculate diffs between plan and
                service configuration
              type: object
              x-kubernetes-preserve-unknown-fields: true
            configurationFrom:
              description: ConfigurationFrom is a way to load configuration values
                from alternative sources, e.g. from secrets The values from these
                sources will override any existing keys defined in Configuration
              items:
                properties:
                  name:
                    description: Name is the name of the configuration parameter
                    minLength: 1
                    type: string
                  secretKeyRef:
                    description: SecretKeyRef is a reference to a key in a secret
                    properties:
                      key:
                        description: Key is they data key in the secret
                        minLength: 1
                        type: string
                      name:
                        description: Name is the name of the secret
                        minLength: 1
                        type: string
                      optional:
                        description: Optional controls whether the secret with the
                          given key must exist
                        type: boolean
                    required:
                    - name
                    type: object
                required:
                - name
                - secretKeyRef
                type: object
              type: array
              x-kubernetes-list-type: set
            credentials:
              description: Credentials is a reference to the credentials object to
                use
              properties:
                group:
                  description: Group is the api group
                  type: string
                kind:
                  description: Kind is the name of the resource under the group
                  type: string
                name:
                  description: Name is name of the resource
                  type: string
                namespace:
                  description: Namespace is the location of the object
                  type: string
                version:
                  description: Version is the group version
                  type: string
              required:
              - group
              - kind
              - name
              - namespace
              - version
              type: object
            kind:
              description: Kind refers to the service type
              minLength: 1
              type: string
            plan:
              description: Plan is the name of the service plan which was used to
                create this service
              minLength: 1
              type: string
          required:
          - kind
          - plan
          type: object
        status:
          description: ServiceStatus defines the observed state of a service
          properties:
            components:
              description: Components is a collection of component statuses
              items:
                description: Component the state of a component of the resource
                properties:
                  detail:
                    description: Detail is additional details on the error is any
                    type: string
                  message:
                    description: Message is a human readable message on the status
                      of the component
                    type: string
                  name:
                    description: Name is the name of the component
                    type: string
                  resource:
                    description: Resource is a reference to the resource
                    properties:
                      group:
                        description: Group is the api group
                        type: string
                      kind:
                        description: Kind is the name of the resource under the group
                        type: string
                      name:
                        description: Name is name of the resource
                        type: string
                      namespace:
                        description: Namespace is the location of the object
                        type: string
                      version:
                        description: Version is the group version
                        type: string
                    required:
                    - group
                    - kind
                    - name
                    - namespace
                    - version
                    type: object
                  status:
                    description: Status is the status of the component
                    type: string
                type: object
              type: array
            configuration:
              description: Configuration are the applied configuration values for
                this service
              type: object
              x-kubernetes-preserve-unknown-fields: true
            message:
              description: Message is the description of the current status
              type: string
            plan:
              description: Plan is the name of the service plan which was used to
                create this service
              type: string
            providerData:
              description: ProviderData is provider specific data
              type: object
              x-kubernetes-preserve-unknown-fields: true
            providerID:
              description: ProviderID is the service identifier in the service provider
              type: string
            status:
              description: Status is the overall status of the service
              type: string
          type: object
      type: object
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
`)

func crdsServicesKoreAppviaIo_servicesYamlBytes() ([]byte, error) {
	return _crdsServicesKoreAppviaIo_servicesYaml, nil
}

func crdsServicesKoreAppviaIo_servicesYaml() (*asset, error) {
	bytes, err := crdsServicesKoreAppviaIo_servicesYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crds/services.kore.appvia.io_services.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"crds/accounts.kore.appvia.io_accountmanagement.yaml":                 crdsAccountsKoreAppviaIo_accountmanagementYaml,
	"crds/apps.kore.appvia.io_appdeployments.yaml":                        crdsAppsKoreAppviaIo_appdeploymentsYaml,
	"crds/apps.kore.appvia.io_installplans.yaml":                          crdsAppsKoreAppviaIo_installplansYaml,
	"crds/aws.compute.kore.appvia.io_eks.yaml":                            crdsAwsComputeKoreAppviaIo_eksYaml,
	"crds/aws.compute.kore.appvia.io_ekscredentials.yaml":                 crdsAwsComputeKoreAppviaIo_ekscredentialsYaml,
	"crds/aws.compute.kore.appvia.io_eksnodegroups.yaml":                  crdsAwsComputeKoreAppviaIo_eksnodegroupsYaml,
	"crds/aws.compute.kore.appvia.io_eksvpcs.yaml":                        crdsAwsComputeKoreAppviaIo_eksvpcsYaml,
	"crds/clusters.compute.kore.appvia.io_clusters.yaml":                  crdsClustersComputeKoreAppviaIo_clustersYaml,
	"crds/clusters.compute.kore.appvia.io_kubernetes.yaml":                crdsClustersComputeKoreAppviaIo_kubernetesYaml,
	"crds/clusters.compute.kore.appvia.io_managedclusterrole.yaml":        crdsClustersComputeKoreAppviaIo_managedclusterroleYaml,
	"crds/clusters.compute.kore.appvia.io_managedclusterrolebinding.yaml": crdsClustersComputeKoreAppviaIo_managedclusterrolebindingYaml,
	"crds/clusters.compute.kore.appvia.io_managedconfig.yaml":             crdsClustersComputeKoreAppviaIo_managedconfigYaml,
	"crds/clusters.compute.kore.appvia.io_managedpodsecuritypoliies.yaml": crdsClustersComputeKoreAppviaIo_managedpodsecuritypoliiesYaml,
	"crds/clusters.compute.kore.appvia.io_managedrole.yaml":               crdsClustersComputeKoreAppviaIo_managedroleYaml,
	"crds/clusters.compute.kore.appvia.io_namespaceclaims.yaml":           crdsClustersComputeKoreAppviaIo_namespaceclaimsYaml,
	"crds/clusters.compute.kore.appvia.io_namespacepolicy.yaml":           crdsClustersComputeKoreAppviaIo_namespacepolicyYaml,
	"crds/config.kore.appvia.io_allocations.yaml":                         crdsConfigKoreAppviaIo_allocationsYaml,
	"crds/config.kore.appvia.io_planpolicies.yaml":                        crdsConfigKoreAppviaIo_planpoliciesYaml,
	"crds/config.kore.appvia.io_plans.yaml":                               crdsConfigKoreAppviaIo_plansYaml,
	"crds/config.kore.appvia.io_secrets.yaml":                             crdsConfigKoreAppviaIo_secretsYaml,
	"crds/core.kore.appvia.io_idp.yaml":                                   crdsCoreKoreAppviaIo_idpYaml,
	"crds/core.kore.appvia.io_oidclient.yaml":                             crdsCoreKoreAppviaIo_oidclientYaml,
	"crds/gcp.compute.kore.appvia.io_organizations.yaml":                  crdsGcpComputeKoreAppviaIo_organizationsYaml,
	"crds/gcp.compute.kore.appvia.io_projectclaims.yaml":                  crdsGcpComputeKoreAppviaIo_projectclaimsYaml,
	"crds/gcp.compute.kore.appvia.io_projects.yaml":                       crdsGcpComputeKoreAppviaIo_projectsYaml,
	"crds/gke.compute.kore.appvia.io_gkecredentials.yaml":                 crdsGkeComputeKoreAppviaIo_gkecredentialsYaml,
	"crds/gke.compute.kore.appvia.io_gkes.yaml":                           crdsGkeComputeKoreAppviaIo_gkesYaml,
	"crds/org.kore.appvia.io_auditevents.yaml":                            crdsOrgKoreAppviaIo_auditeventsYaml,
	"crds/org.kore.appvia.io_members.yaml":                                crdsOrgKoreAppviaIo_membersYaml,
	"crds/org.kore.appvia.io_teaminvitations.yaml":                        crdsOrgKoreAppviaIo_teaminvitationsYaml,
	"crds/org.kore.appvia.io_teams.yaml":                                  crdsOrgKoreAppviaIo_teamsYaml,
	"crds/org.kore.appvia.io_users.yaml":                                  crdsOrgKoreAppviaIo_usersYaml,
	"crds/security.kore.appvia.io_securityoverviews.yaml":                 crdsSecurityKoreAppviaIo_securityoverviewsYaml,
	"crds/security.kore.appvia.io_securityrules.yaml":                     crdsSecurityKoreAppviaIo_securityrulesYaml,
	"crds/security.kore.appvia.io_securityscanresults.yaml":               crdsSecurityKoreAppviaIo_securityscanresultsYaml,
	"crds/services.kore.appvia.io_servicecredentials.yaml":                crdsServicesKoreAppviaIo_servicecredentialsYaml,
	"crds/services.kore.appvia.io_servicekinds.yaml":                      crdsServicesKoreAppviaIo_servicekindsYaml,
	"crds/services.kore.appvia.io_serviceplans.yaml":                      crdsServicesKoreAppviaIo_serviceplansYaml,
	"crds/services.kore.appvia.io_serviceproviders.yaml":                  crdsServicesKoreAppviaIo_serviceprovidersYaml,
	"crds/services.kore.appvia.io_services.yaml":                          crdsServicesKoreAppviaIo_servicesYaml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"crds": {nil, map[string]*bintree{
		"accounts.kore.appvia.io_accountmanagement.yaml":                 {crdsAccountsKoreAppviaIo_accountmanagementYaml, map[string]*bintree{}},
		"apps.kore.appvia.io_appdeployments.yaml":                        {crdsAppsKoreAppviaIo_appdeploymentsYaml, map[string]*bintree{}},
		"apps.kore.appvia.io_installplans.yaml":                          {crdsAppsKoreAppviaIo_installplansYaml, map[string]*bintree{}},
		"aws.compute.kore.appvia.io_eks.yaml":                            {crdsAwsComputeKoreAppviaIo_eksYaml, map[string]*bintree{}},
		"aws.compute.kore.appvia.io_ekscredentials.yaml":                 {crdsAwsComputeKoreAppviaIo_ekscredentialsYaml, map[string]*bintree{}},
		"aws.compute.kore.appvia.io_eksnodegroups.yaml":                  {crdsAwsComputeKoreAppviaIo_eksnodegroupsYaml, map[string]*bintree{}},
		"aws.compute.kore.appvia.io_eksvpcs.yaml":                        {crdsAwsComputeKoreAppviaIo_eksvpcsYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_clusters.yaml":                  {crdsClustersComputeKoreAppviaIo_clustersYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_kubernetes.yaml":                {crdsClustersComputeKoreAppviaIo_kubernetesYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_managedclusterrole.yaml":        {crdsClustersComputeKoreAppviaIo_managedclusterroleYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_managedclusterrolebinding.yaml": {crdsClustersComputeKoreAppviaIo_managedclusterrolebindingYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_managedconfig.yaml":             {crdsClustersComputeKoreAppviaIo_managedconfigYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_managedpodsecuritypoliies.yaml": {crdsClustersComputeKoreAppviaIo_managedpodsecuritypoliiesYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_managedrole.yaml":               {crdsClustersComputeKoreAppviaIo_managedroleYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_namespaceclaims.yaml":           {crdsClustersComputeKoreAppviaIo_namespaceclaimsYaml, map[string]*bintree{}},
		"clusters.compute.kore.appvia.io_namespacepolicy.yaml":           {crdsClustersComputeKoreAppviaIo_namespacepolicyYaml, map[string]*bintree{}},
		"config.kore.appvia.io_allocations.yaml":                         {crdsConfigKoreAppviaIo_allocationsYaml, map[string]*bintree{}},
		"config.kore.appvia.io_planpolicies.yaml":                        {crdsConfigKoreAppviaIo_planpoliciesYaml, map[string]*bintree{}},
		"config.kore.appvia.io_plans.yaml":                               {crdsConfigKoreAppviaIo_plansYaml, map[string]*bintree{}},
		"config.kore.appvia.io_secrets.yaml":                             {crdsConfigKoreAppviaIo_secretsYaml, map[string]*bintree{}},
		"core.kore.appvia.io_idp.yaml":                                   {crdsCoreKoreAppviaIo_idpYaml, map[string]*bintree{}},
		"core.kore.appvia.io_oidclient.yaml":                             {crdsCoreKoreAppviaIo_oidclientYaml, map[string]*bintree{}},
		"gcp.compute.kore.appvia.io_organizations.yaml":                  {crdsGcpComputeKoreAppviaIo_organizationsYaml, map[string]*bintree{}},
		"gcp.compute.kore.appvia.io_projectclaims.yaml":                  {crdsGcpComputeKoreAppviaIo_projectclaimsYaml, map[string]*bintree{}},
		"gcp.compute.kore.appvia.io_projects.yaml":                       {crdsGcpComputeKoreAppviaIo_projectsYaml, map[string]*bintree{}},
		"gke.compute.kore.appvia.io_gkecredentials.yaml":                 {crdsGkeComputeKoreAppviaIo_gkecredentialsYaml, map[string]*bintree{}},
		"gke.compute.kore.appvia.io_gkes.yaml":                           {crdsGkeComputeKoreAppviaIo_gkesYaml, map[string]*bintree{}},
		"org.kore.appvia.io_auditevents.yaml":                            {crdsOrgKoreAppviaIo_auditeventsYaml, map[string]*bintree{}},
		"org.kore.appvia.io_members.yaml":                                {crdsOrgKoreAppviaIo_membersYaml, map[string]*bintree{}},
		"org.kore.appvia.io_teaminvitations.yaml":                        {crdsOrgKoreAppviaIo_teaminvitationsYaml, map[string]*bintree{}},
		"org.kore.appvia.io_teams.yaml":                                  {crdsOrgKoreAppviaIo_teamsYaml, map[string]*bintree{}},
		"org.kore.appvia.io_users.yaml":                                  {crdsOrgKoreAppviaIo_usersYaml, map[string]*bintree{}},
		"security.kore.appvia.io_securityoverviews.yaml":                 {crdsSecurityKoreAppviaIo_securityoverviewsYaml, map[string]*bintree{}},
		"security.kore.appvia.io_securityrules.yaml":                     {crdsSecurityKoreAppviaIo_securityrulesYaml, map[string]*bintree{}},
		"security.kore.appvia.io_securityscanresults.yaml":               {crdsSecurityKoreAppviaIo_securityscanresultsYaml, map[string]*bintree{}},
		"services.kore.appvia.io_servicecredentials.yaml":                {crdsServicesKoreAppviaIo_servicecredentialsYaml, map[string]*bintree{}},
		"services.kore.appvia.io_servicekinds.yaml":                      {crdsServicesKoreAppviaIo_servicekindsYaml, map[string]*bintree{}},
		"services.kore.appvia.io_serviceplans.yaml":                      {crdsServicesKoreAppviaIo_serviceplansYaml, map[string]*bintree{}},
		"services.kore.appvia.io_serviceproviders.yaml":                  {crdsServicesKoreAppviaIo_serviceprovidersYaml, map[string]*bintree{}},
		"services.kore.appvia.io_services.yaml":                          {crdsServicesKoreAppviaIo_servicesYaml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
