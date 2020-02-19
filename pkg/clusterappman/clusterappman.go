/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
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

package clusterappman

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	kcore "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/clusterapp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	rc "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// KoreNamespace is the namespace where the clusterappmanager and repo are deployed
	KoreNamespace string = "kore"
	// ParamsConfigMap provides the customisations for deplpoyments carried out from here
	ParamsConfigMap string = "kore-cluster-config"
	// ParamsConfigKey is the configmap key used to store the parameters
	ParamsConfigKey string = "clusterconfig"
	// StatusCongigMap is the name of the configmap object used to store kore cluster status
	StatusCongigMap string = "kore-cluster-status"
	// StatusConfigMapComponentsKey is the key in the configmap used for all conditions
	StatusConfigMapComponentsKey string = "components"
	appControlNamespsace         string = "application-system"
	captainNamespace             string = "captain-system"
)

type clusterappmanImpl struct {
	// client is the kubernetes client to use
	Client        kubernetes.Interface
	RuntimeClient rc.Client
	ClusterApps   []clusterapp.Instance
	cfg           *rest.Config
}

// manifest defines the data types that are required to initialise a clusterapp from
// embebed manifest data.
type manifest struct {
	EmededManifests []string
	Name            string
	Namespace       string
	EnsureNamespace bool
	DeployTimeOut   time.Duration
}

var (
	// cas is the map of clusterapp instances
	cas map[string]*clusterapp.Instance = make(map[string]*clusterapp.Instance)

	// mm defines all the embeded manifests and data required to initialise clusterappman
	mm []manifest = []manifest{
		{
			Name: "Application Controller",
			EmededManifests: []string{
				"application-controller/application-all.yaml",
			},
			Namespace:       appControlNamespsace,
			EnsureNamespace: false,
			DeployTimeOut:   3 * time.Minute,
		},
		{
			Name: "Helm Chart Operator",
			EmededManifests: []string{
				"captain/clusterrolebinding.yaml",
				"captain/captain.yaml",
				"captain/captain-application.yaml",
			},
			Namespace:       captainNamespace,
			EnsureNamespace: true,
			DeployTimeOut:   3 * time.Minute,
		},
		{
			Name: "Kore Helm Repository",
			EmededManifests: []string{
				"kore-helm-repo/application-all.yaml",
			},
			Namespace:       KoreNamespace,
			EnsureNamespace: false,
			DeployTimeOut:   3 * time.Minute,
		},
	}
)

// New is responsible for creating the clusterappman server
func New(config Config) (Interface, error) {
	if err := config.IsValid(); err != nil {
		return nil, err
	}
	var client kubernetes.Interface

	cfg, err := makeKubernetesConfig(config.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes config: %s", err)
	}

	options, err := clusterapp.GetClientOptions()
	if err != nil {
		return nil, fmt.Errorf("failed getting client options (schemes): %s", err)
	}
	cc, err := rc.New(cfg, options)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes runtime client: %s", err)
	}

	client, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed creating kubernetes client: %s", err)
	}
	return &clusterappmanImpl{
		Client:        client,
		RuntimeClient: cc,
	}, nil
}

// Run is responsible for starting the services:
// it should start threads for monitoring cluster apps
// it should only return if "initial" kore deployments don't become ready
func (s clusterappmanImpl) Run(ctx context.Context) error {
	// Maybe this whole thing runs in a loop - ensuring no manual change?
	// deploy / reconcile the Application kind controller
	// deploy / reconcile the Helm operator
	// deploy / reconcile the Appvia helm repo

	// initialise clusterapps and parse all the manifests
	// fail early if we can't even load the manifests
	// we should add a test for this!
	logger := log.WithFields(log.Fields{
		"service": "clusterappman",
	})
	logger.Info("attempting to reconcile the applications incluster")

	var wg sync.WaitGroup
	// componant updates channel
	ch := make(chan *kcore.Component, len(mm))

	// make this tesable
	logger.Info("loading manifests")
	if err := LoadAllManifests(s.RuntimeClient); err != nil {
		return fmt.Errorf("failed loading manifests - %s", err)
	}
	for _, m := range mm {
		// Get the pre-loaded cluster application
		ca, ok := cas[m.Name]
		if !ok {
			return fmt.Errorf("failed creating all manifests as could not find manifest %s", m.Name)
		}
		if m.EnsureNamespace {
			logger.Infof("ensuring namespace %s exists", m.Namespace)
			if err := ensureNamespace(ctx, s.RuntimeClient, m.Namespace); err != nil {
				return fmt.Errorf("failed creating namespace %s: %s", m.Namespace, err)
			}
		}
		// Write all objects to the API on a seperate thread...
		if err := ca.CreateOrUpdate(ctx, m.Namespace); err != nil {
			return fmt.Errorf("failed to create or update '%s' deployment: %s", m.Name, err)
		}
		// TODO: write a status monitor entry point
		deployCtx, cancel := context.WithTimeout(ctx, m.DeployTimeOut)
		// make sure we run this cancel whatevs
		defer cancel()
		wg.Add(1)
		// start a deployment thread:
		logger.Infof("starting to wait for '%s' to become ready", ca.Component.Name)
		go ca.WaitForReadyOrTimeout(deployCtx, ch, &wg)
	}
	wg.Wait()
	close(ch)

	// now gather up all the component slices as they complete...
	logger.Infof("waiting for %d cluster apps", len(mm))
	var cs = make([]*kcore.Component, len(mm))
	for i := range cs {
		// get the first component "reason"
		c := <-ch
		// get the corresponding app
		ca, ok := cas[c.Name]
		// overwrite the component with the recieved updated "reason"
		ca.Component = c
		if !ok {
			return fmt.Errorf("missing application from component -  %s", c.Name)
		}
		if ca.Component.Status == kcore.SuccessStatus {
			logger.Infof("cluster app '%s' is ready", ca.Component.Name)
		} else {
			logger.Errorf(
				"cluster app '%s' has timed out and has status [%s] due to [%s]",
				ca.Component.Name,
				ca.Component.Status,
				ca.Component.Detail,
			)
		}
		cs[i] = ca.Component
	}
	var components kcore.Components = cs
	logger.Logger.Infof("saving status")
	if err := createStatusConfig(ctx, s.RuntimeClient, components); err != nil {
		return fmt.Errorf("error reporting status: %s", err)
	}
	// TODO:
	// 3. Update check at other side...

	// Further TODO:
	// 2. start kore-helm-repo reconciler
	// 3. watch for all application types from operators in our repo to monitor readiness across reszt of cluster estate...
	// 3. de-serialize the parameters and deploy all the CRD's (probably using templates)
	// 4. Ensure we have application types or status known for all CRD's and add to components list
	// 4. re-start loop

	// we shouldn't get here!
	return nil
}

// Stop is responsible for trying to stop services
func (s clusterappmanImpl) Stop(context.Context) error {
	return nil
}

func (s clusterappmanImpl) UpgradeClient(options rc.Options) (err error) {
	// TODO may need to make this thread safe!
	s.RuntimeClient, err = rc.New(s.cfg, options)
	if err != nil {
		return err
	}
	return nil
}
