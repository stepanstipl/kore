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

package korectl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/apiserver/types"

	"github.com/appvia/kore/pkg/utils"

	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"gopkg.in/yaml.v2"

	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	createClusterLongDescription = `
Provides the ability to provision a kubernetes cluster in the team. The cluster
itself is provisioned from a predefined plan (a template). You can view the plans
available to you via $ korectl get plans. Once the cluster has been built the
members of your team can gain access via running $ korectl login.

Note: you retrieve a list of all the plans available to you via:
$ korectl get plans
$ korectl get plans <name> -o yaml

Examples:
$ korectl -t <myteam> create cluster dev --plan gke-development --allocation <allocation_name>

# Create a cluster and provision some namespaces on there as well
$ korectl -t <myteam> create cluster dev --plan gke-development -a <name> --namespace=app1,app2

# Check the status of the cluster
$ korectl -t <myteam> get cluster dev -o yaml

Now update your kubeconfig to use your team's provisioned cluster.
$ korectl kubeconfig -t <myteam>

This will modify your ${HOME}/.kube/config. Now you can use 'kubectl' to interact with your team's cluster.
`
)

// GetCreateClusterCommand returns the command to create clusters
// @Note: we probably need to move this cluster provisioning off a plan into the API itself
// and offload it from the CLI - but needs discussion first.
func GetCreateClusterCommand(config *Config) *cli.Command {
	return &cli.Command{
		Name:        "cluster",
		Aliases:     []string{"clusters"},
		Description: formatLongDescription(createClusterLongDescription),
		Usage:       "create a kubernetes cluster within the team",
		ArgsUsage:   "<name> [options]",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "plan",
				Aliases: []string{"p"},
				Usage:   "the plan which this cluster will be templated from `NAME`",
			},
			&cli.StringSliceFlag{
				Name:  "namespace",
				Usage: "you can preprovision a collection namespaces on this cluster as well `NAMES`",
			},
			&cli.StringFlag{
				Name:    "allocation",
				Aliases: []string{"a"},
				Usage:   "the name of the allocated credentials to use for this cluster `NAME`",
			},
			&cli.StringSliceFlag{
				Name:    "plan-param",
				Aliases: []string{"param"},
				Usage:   "used to override the plan parameters",
			},
			&cli.BoolFlag{
				Name:  "show-time",
				Usage: "shows the time it took to successfully provision a new cluster `BOOL`",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "generate the cluster specification but does not apply `BOOL`",
			},
		},

		Before: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("the cluster should have a name")
			}

			return nil
		},

		Action: func(ctx *cli.Context) error {
			name := ctx.Args().First()
			allocation := ctx.String("allocation")
			dry := ctx.Bool("dry-run")
			namespaces := ctx.StringSlice("namespace")
			plan := ctx.String("plan")
			team := ctx.String("team")
			wait := ctx.Bool("no-wait")

			if team == "" {
				return errTeamParameterMissing
			}

			if allocation == "" {
				return fmt.Errorf("no allocation defined, please use $ korectl get allocations -t %s", team)
			}
			if plan == "" {
				return fmt.Errorf("no plan defined, please use: $ korectl get plans")
			}

			found, err := TeamResourceExists(config, team, "clusters", name)
			if err != nil {
				return err
			}
			if found {
				return fmt.Errorf("cluster %q already exists", name)
			}

			whoami, err := GetWhoAmI(config)
			if err != nil {
				return err
			}

			planObj := &configv1.Plan{}
			if err := GetResource(config, "plan", plan, planObj); err != nil {
				return fmt.Errorf("plan %q does not exist", plan)
			}

			credsAlloc := &configv1.Allocation{}
			if err := GetTeamResource(config, team, "allocation", allocation, credsAlloc); err != nil {
				if reqErr, ok := err.(*RequestError); ok {
					if reqErr.statusCode == http.StatusNotFound {
						return fmt.Errorf("allocation %q does not exist", allocation)
					}
				} else {
					return fmt.Errorf("failed to retrieve the allocation from api: %s", err)
				}
			}

			cluster, err := createClusterObject(ctx, name, team, whoami, planObj, credsAlloc)
			if err != nil {
				return err
			}

			if dry {
				out, err := utils.EncodeRuntimeObjectToYAML(cluster)
				if err != nil {
					return fmt.Errorf("failed to parse cluster object: %s", err)
				}
				fmt.Println(string(out))
				return nil
			}

			if err := CreateTeamResource(config, team, "Cluster", name, cluster); err != nil {
				return err
			}

			// @step: create a start time
			now := time.Now()

			// @step: we need to construct the provider type
			if err := WaitForResourceCheck(context.Background(), config, team, "Cluster", name, wait); err != nil {
				return err
			}
			if ctx.Bool("show-time") {
				fmt.Printf("Provisioning took: %s\n", time.Since(now))
			}

			// @step: do we need to provision any namespaces? - note the split and joining
			// allows for --namespace a,b,c
			var list []string
			for _, x := range namespaces {
				list = append(list, strings.Split(x, ",")...)
			}

			for _, x := range list {
				if _, err := CreateClusterNamespace(config, name, team, x, dry); err != nil {
					return fmt.Errorf("trying to provision namespace claim: %s on cluster: %s", x, err)
				}
			}

			// @step: print a the message
			fmt.Printf("\nYou can retrieve your kubeconfig via: $ korectl kubeconfig -t %s\n", team)

			return nil
		},
	}
}

func createClusterObject(ctx *cli.Context, name, team string, whoAmI *types.WhoAmI, plan *configv1.Plan, credsAlloc *configv1.Allocation) (*clustersv1.Cluster, error) {
	var configuration map[string]interface{}
	if err := json.Unmarshal(plan.Spec.Configuration.Raw, &configuration); err != nil {
		return nil, fmt.Errorf("failed to parse plan configuration: %s", err)
	}

	planParams, err := parsePlanParams(ctx.StringSlice("plan-param"))
	if err != nil {
		return nil, err
	}

	if _, ok := planParams["clusterUsers"]; !ok {
		planParams["clusterUsers"] = []map[string]interface{}{
			{
				"username": whoAmI.Username,
				"roles":    []string{"cluster-admin"},
			},
		}
	}

	for k, v := range planParams {
		configuration[k] = v
	}

	configurationRaw, err := json.Marshal(configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to process cluster configuration: %s", err)
	}

	cluster := &clustersv1.Cluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Cluster",
			APIVersion: clustersv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: team,
		},
		Spec: clustersv1.ClusterSpec{
			Kind:          plan.Spec.Kind,
			Plan:          plan.Name,
			Configuration: v1beta1.JSON{Raw: configurationRaw},
			Credentials:   credsAlloc.Spec.Resource,
		},
	}

	return cluster, nil
}

// CreateClusterNamespace is called to provision a namespace on the cluster
func CreateClusterNamespace(config *Config, clusterName, team, name string, dry bool) (*clustersv1.NamespaceClaim, error) {
	resourceName := fmt.Sprintf("%s-%s", clusterName, name)
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
				Namespace: team,
				Name:      clusterName,
			},
		},
	}
	if dry {
		return nil, yaml.NewEncoder(os.Stdout).Encode(object)
	}

	found, err := TeamResourceExists(config, team, kind, resourceName)
	if err != nil {
		return nil, err
	}
	if found {
		fmt.Printf("--> Namespace: %s already exists, skipping creation\n", name)

		return object, nil
	}
	fmt.Printf("--> Attempting to create namespace: %s\n", name)

	return object, CreateTeamResource(config, team, kind, name, object)
}

func parsePlanParams(params []string) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for _, param := range params {
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
		res[name] = parsed
	}
	return res, nil
}
