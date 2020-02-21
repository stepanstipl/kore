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

package clusterapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	kcore "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Instance provides access to a Cluster Application
// a cluster app is a facility, running in a cluster
// it is provided directly or indirectly by cluster manager
// it provides facilities to manage a cluster app lifecycle
type Instance struct {
	// client provides access to a kubernetes api
	client client.Client
	// Resources are the K8 objects created for deployment / update
	Resources []runtime.Object
	// ApplicationObject provides access to the application kind
	ApplicationObject runtime.Object
	// Kore standard component and status information
	Component *kcore.Component
}

// NewAppFromManifestFiles creates a new cluster application
func NewAppFromManifestFiles(client client.Client, name string, manifestfiles []http.File) (Instance, error) {
	ca := Instance{
		client:    client,
		Resources: make([]runtime.Object, 0),
		Component: &kcore.Component{
			Name:    name,
			Status:  kcore.Unknown,
			Detail:  "Deployment is defined but not started",
			Message: "Not yet deplopyed",
		},
	}
	// for all the embedded paths specified...
	for _, file := range manifestfiles {
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return ca, err
		}
		// Pass the yaml ([]bytes) and function (addAllToScheme) to create the runtim.Objects
		apiObjs, err := kubernetes.ParseK8sYaml(fileBytes, addAllToScheme)
		if err != nil {
			return ca, err
		}
		for _, obj := range apiObjs {
			if obj.GetObjectKind().GroupVersionKind().Kind == "Application" {
				if ca.ApplicationObject != nil {
					return ca, fmt.Errorf("only one application kind per kore cluster app is supported in cluster app %s", name)
				}
				ca.ApplicationObject = obj
			}
			ca.Resources = append(ca.Resources, obj)
		}
	}
	return ca, nil
}

// CreateOrUpdate will deploy or update all the manifets
// the deafultNamespace is used if not otherwise specified.
func (ca Instance) CreateOrUpdate(ctx context.Context, defaultNamepsace string) error {
	for _, res := range ca.Resources {
		objMeta, err := getObjMeta(res)
		if err != nil {
			log.Errorf("error getting metadata for %v - %s", res, err)
		}
		if err := setMissingNamespace(defaultNamepsace, res); err != nil {
			log.Debugf("error setting namespace for %v - %s", res, err)
		}
		log.Debugf(
			"deploying %s/%s",
			objMeta.Namespace,
			objMeta.Name,
		)
		// Create / update / replace resources as required
		if _, err := kubernetes.CreateOrUpdate(ctx, ca.client, res); err != nil {
			return err
		}
	}
	ca.Component.Status = kcore.PendingStatus
	ca.Component.Detail = "The deployment of all manifests have been accepted by the API"
	ca.Component.Message = "Deployment started"

	return nil
}

// WaitForReadyOrTimeout will wait a reasonable time (defined in context) until a resource is ready
// if the resource become ready, it will update the channel with the component (and status)
// if there is any error or a timeout, it will update the channel with the details on the component
func (ca Instance) WaitForReadyOrTimeout(ctx context.Context, respond chan<- *kcore.Component, wg *sync.WaitGroup) {
	defer wg.Done()

	// here we handle channels and wait groups not errors so pass the timeout context on:
	if err := waitOnApplicationStatus(ctx, &ca); err != nil {
		log.Errorf("error with %s", ca.Component.Name)
		ca.Component.Status = kcore.Unknown
		ca.Component.Message = fmt.Sprintf("An error occured when checking for the status of %s", ca.Component.Name)
		ca.Component.Detail = fmt.Sprintf("An error occured waiting for status %s", err)
	}
	respond <- ca.Component
}

// GetApplicationObjectName will inspect the metadata and return the object name
// will return error if not defined
func (ca Instance) GetApplicationObjectName() string {
	if ca.ApplicationObject == nil {
		return ""
	}
	metaObj, err := getObjMeta(ca.ApplicationObject)
	if err != nil {
		return ""
	}
	return metaObj.Name
}
