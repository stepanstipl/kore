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

package create

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/apiserver/types"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	cmdutils "github.com/appvia/kore/pkg/cmd/utils"

	"github.com/spf13/cobra"
	apiexts "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createClusterLongDescription = `
Provides the ability to provision a kubernetes cluster in the team. The cluster
itself is provisioned from a predefined plan (a template). You can view the plans
available to you via $ kore get plans. Once the cluster has been built the
members of your team can gain access via running $ kore login.

Note: you retrieve a list of all the plans available to you via:
$ kore get plans
$ kore get plans <name> -o yaml

Examples:
$ kore -t <myteam> create cluster dev --plan gke-development --allocation <allocation_name>

# Create a cluster and provision some clusters on there as well
$ kore -t <myteam> create cluster dev --plan gke-development -a <name> --cluster=app1,app2

# Check the status of the cluster
$ kore -t <myteam> get cluster dev -o yaml

Now update your kubeconfig to use your team's provisioned cluster.
$ kore kubeconfig -t <myteam>

This will modify your ${HOME}/.kube/config. Now you can use 'kubectl' to interact with your team's cluster.
`
)

// CreateClusterOptions is used to provision a team
type CreateClusterOptions struct {
	cmdutil.Factory
	// Name is the name of the cluster
	Name string
	// Description is a description for the cluster
	Description string
	// Plan is the plan to build the cluster off
	Plan string
	// Team string
	Team string
	// Clusters is collection of clusters to create
	Clusters []string
	// Allocation is the allocation to build the cluster off
	Allocation string
	// PlanParams is a collection of plan overrides
	PlanParams []string
	// Namespaces is a collection of namespaces to provision
	Namespaces []string
	// NoWait indicates if we should wait for a resource to provision
	NoWait bool
	// ShowTime indicate we should show the build time
	ShowTime bool
}

// NewCmdCreateCluster returns the create cluster command
func NewCmdCreateCluster(factory cmdutil.Factory) *cobra.Command {
	o := &CreateClusterOptions{Factory: factory}

	command := &cobra.Command{
		Use:     "cluster",
		Short:   "Create a kubernetes cluster within the team",
		Long:    createClusterLongDescription,
		Example: "kore create cluster -a <allocation> -p <plan> [-t|--team]",
		Run:     cmdutil.DefaultRunFunc(o),
	}

	flags := command.Flags()
	flags.StringVarP(&o.Allocation, "allocation", "a", "", "name of the allocated to use for this cluster `NAME`")
	flags.StringVarP(&o.Plan, "plan", "p", "", "plan which this cluster will be templated from `NAME`")
	flags.StringVarP(&o.Description, "description", "d", "", "a short description for the cluster `DESCRIPTION`")
	flags.StringSliceVar(&o.PlanParams, "param", []string{}, "preprovision a collection namespaces on this cluster as well `NAMES`")
	flags.StringSliceVar(&o.Namespaces, "namespaces", []string{}, "used to override the plan parameters `KEY=VALUE`")
	flags.BoolVarP(&o.ShowTime, "show-time", "T", false, "shows the time it took to successfully provision a new cluster `BOOL`")

	cmdutils.MustMarkFlagRequired(command, "allocation")
	cmdutils.MustMarkFlagRequired(command, "plan")

	// @step: register the autocompletions
	cmdutils.MustRegisterFlagCompletionFunc(command, "allocation", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
		list := &configv1.AllocationList{}
		if err := o.ClientWithTeamResource(cmdutil.GetTeam(cmd), o.Resources().MustLookup("allocation")).Result(list).Get().Error(); err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		var filtered []string
		for _, x := range list.Items {
			switch x.Spec.Resource.Kind {
			case "GKECredentials", "EKSCredentials", "ProjectClaim":
				filtered = append(filtered, x.Name)
			}
		}

		return filtered, cobra.ShellCompDirectiveNoFileComp
	})

	// @TODO would be nice to filter on the allocation here as well - i.e. chosen GKE, only show GKE plans etc
	cmdutils.MustRegisterFlagCompletionFunc(command, "plan", func(cmd *cobra.Command, args []string, complete string) ([]string, cobra.ShellCompDirective) {
		suggestions, err := o.Resources().LookupResourceNames("plan", "")
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return suggestions, cobra.ShellCompDirectiveNoFileComp
	})

	// @TODO add a autogen for the plan parameters? - perhaps when we start doing local caching

	return command
}

// Validate is called to check the options
func (o *CreateClusterOptions) Validate() error {
	if o.Team == "" {
		return errors.ErrTeamMissing
	}
	if o.Name == "" {
		return errors.ErrMissingResourceName
	}

	found, err := o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("allocation")).Name(o.Allocation).Exists()
	if err != nil {
		return err
	}
	if !found {
		return errors.NewResourceNotFoundWithKind(o.Allocation, "allocation")
	}

	found, err = o.ClientWithResource(o.Resources().MustLookup("plan")).Name(o.Plan).Exists()
	if err != nil {
		return err
	}
	if !found {
		return errors.NewResourceNotFoundWithKind(o.Plan, "plan")
	}

	match := regexp.MustCompile("^.*=.*$")

	for _, x := range o.PlanParams {
		if !match.MatchString(x) {
			return errors.NewInvalidParamError("param", x)
		}
	}

	return nil
}

// Run implements the action
func (o *CreateClusterOptions) Run() error {
	// @step: generate the cluster configuration
	config, err := o.CreateClusterConfiguration()
	if err != nil {
		return err
	}

	now := time.Now()
	// @step: provision and wait if required
	err = o.WaitForCreation(
		o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("cluster")).
			Name(o.Name).
			Payload(config),
		o.NoWait,
	)
	if err != nil {
		return err
	}
	if o.ShowTime {
		o.Println("Provisioning took: %s", time.Since(now))
	}

	// @step: we need to provision any namespace on the cluster
	var list []string

	for _, x := range o.Namespaces {
		list = append(list, strings.Split(x, ",")...)
	}

	for _, x := range list {
		if err := o.CreateClusterNamespace(x); err != nil {
			return fmt.Errorf("trying to provision namespace claim: %s on cluster: %s", x, err)
		}
	}
	o.Println("\nYou can retrieve your kubeconfig via: $ kore kubeconfig -t %s", o.Team)

	return nil
}

// CreateClusterConfiguration is responsible for generating the cluster config
func (o *CreateClusterOptions) CreateClusterConfiguration() (*clustersv1.Cluster, error) {
	// @step: retrieve the plan, allocation and user auth
	plan, err := o.GetPlan()
	if err != nil {
		return nil, err
	}
	allocation, err := o.GetAllocation()
	if err != nil {
		return nil, err
	}
	userauth, err := o.GetUserAuth()
	if err != nil {
		return nil, err
	}

	var configuration map[string]interface{}

	if err := json.Unmarshal(plan.Spec.Configuration.Raw, &configuration); err != nil {
		return nil, fmt.Errorf("failed to parse plan configuration: %s", err)
	}

	params, err := o.ParsePlanParams()
	if err != nil {
		return nil, err
	}

	// @step: inject ourself as the cluster admin
	if _, ok := params["clusterUsers"]; !ok {
		params["clusterUsers"] = []map[string]interface{}{
			{
				"username": userauth.Username,
				"roles":    []string{"cluster-admin"},
			},
		}
	}

	// @step: copy the plan parameters into the cluster configuration
	for k, v := range params {
		configuration[k] = v
	}

	// @step: json encode the cluster parameters
	cc := &bytes.Buffer{}
	if err := json.NewEncoder(cc).Encode(configuration); err != nil {
		return nil, fmt.Errorf("failed to process cluster configuration: %s", err)
	}

	cluster := &clustersv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: clustersv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.Name,
			Namespace: o.Team,
		},
		Spec: clustersv1.ClusterSpec{
			Kind:          plan.Spec.Kind,
			Plan:          plan.Name,
			Configuration: apiexts.JSON{Raw: cc.Bytes()},
			Credentials:   allocation.Spec.Resource,
		},
	}

	return cluster, nil
}

// CreateClusterNamespace is called to provision a namespace on the cluster
func (o *CreateClusterOptions) CreateClusterNamespace(name string) error {
	// @note we have have to prefix the name with the cluster here else since these resources
	// are going into a common team namespace. User creates the same namespace name across multiple
	// clusters will conflict. i.e. i've prod and dev and want to create namespace app1 in both

	resourceName := fmt.Sprintf("%s-%s", o.Name, name)
	kind := "namespaceclaims"

	object := &clustersv1.NamespaceClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clustersv1.GroupVersion.String(),
			Kind:       "NamespaceClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
		},
		Spec: clustersv1.NamespaceClaimSpec{
			Name: name,
			Cluster: corev1.Ownership{
				Group:     clustersv1.GroupVersion.Group,
				Version:   clustersv1.GroupVersion.Version,
				Kind:      "Kubernetes",
				Namespace: o.Team,
				Name:      o.Name,
			},
		},
	}

	resource, err := o.Resources().Lookup(kind)
	if err != nil {
		return err
	}

	found, err := o.ClientWithTeamResource(o.Team, resource).Name(resourceName).Exists()
	if err != nil {
		return err
	}
	if found {
		o.Println("--> Namespace: %s already exists, skipping creation", name)

		return nil
	}
	o.Println("--> Attempting to create namespace: %s", name)

	return o.WaitForCreation(
		o.ClientWithTeamResource(o.Team, resource).
			Name(resourceName).
			Payload(object),
		o.NoWait,
	)
}

// ParsePlanParams is responsible for parsing the plan overrides
func (o *CreateClusterOptions) ParsePlanParams() (map[string]interface{}, error) {
	params := map[string]interface{}{}

	for _, param := range o.PlanParams {
		parts := regexp.MustCompile(`\s*=\s*`).Split(strings.TrimSpace(param), 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid plan-param value %q, you must use param=<JSON value> format", param)
		}
		name := parts[0]
		jsonValue := parts[1]

		var parsed interface{}
		if err := json.Unmarshal([]byte(jsonValue), &parsed); err != nil {
			return nil, err
		}
		params[name] = parsed
	}

	return params, nil
}

// GetPlan retrieve the requested cluster plan
func (o *CreateClusterOptions) GetPlan() (*configv1.Plan, error) {
	plan := &configv1.Plan{}

	return plan, o.ClientWithResource(o.Resources().MustLookup("plan")).
		Name(o.Plan).
		Result(plan).
		Get().
		Error()
}

// GetAllocation retrieve the request allocation
func (o *CreateClusterOptions) GetAllocation() (*configv1.Allocation, error) {
	allocation := &configv1.Allocation{}

	return allocation, o.ClientWithTeamResource(o.Team, o.Resources().MustLookup("allocation")).
		Name(o.Allocation).
		Result(allocation).
		Get().
		Error()
}

// GetUserAuth returns the whoami
func (o *CreateClusterOptions) GetUserAuth() (*types.WhoAmI, error) {
	who := &types.WhoAmI{}

	return who, o.ClientWithEndpoint("/whoami").Result(who).Get().Error()
}
