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

package clusterapp

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	kcore "github.com/appvia/kore/pkg/apis/core/v1"
	korev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"

	"k8s.io/apimachinery/pkg/runtime"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ClusterAppControllerComponentName is the reserved name for the Application controller
	ClusterAppControllerComponentName string = "Application Controller"
)

var (
	appControllerStatus bool
	mu                  = &sync.Mutex{}
)

func getAppControllerStatus() bool {
	mu.Lock()
	defer mu.Unlock()
	return appControllerStatus
}

func setAppControllerStatus(s bool) {
	mu.Lock()
	defer mu.Unlock()
	appControllerStatus = s
}

// Instance provides access to a Cluster Application
// a cluster app is a facility, running in a cluster
// it is provided directly or indirectly by cluster manager
// it provides facilities to manage a cluster app lifecycle
type Instance struct {
	// client provides access to a kubernetes api
	client client.Client
	// PreDeleteResources are the K8 objects created for deletion
	PreDeleteResources []runtime.Object
	// Resources are the K8 objects created for deployment / update
	Resources []runtime.Object
	// ApplicationObject provides access to the application kind
	ApplicationObject runtime.Object
	// Kore standard component and status information
	Component *kcore.Component
}

// NewAppFromManifestFiles creates a new cluster application
func NewAppFromManifestFiles(client client.Client, name string, manifestfiles, deleteResfiles []http.File) (Instance, error) {
	ca := Instance{
		client:    client,
		Resources: make([]runtime.Object, 0),
		Component: &kcore.Component{
			Name:    name,
			Status:  kcore.Unknown,
			Message: "Component is not yet deployed",
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
	for _, file := range deleteResfiles {
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return ca, err
		}
		// Pass the yaml ([]bytes) and function (addAllToScheme) to create the runtim.Objects
		apiObjs, err := kubernetes.ParseK8sYaml(fileBytes, addAllToScheme)
		if err != nil {
			return ca, err
		}
		ca.PreDeleteResources = append(ca.PreDeleteResources, apiObjs...)
	}

	return ca, nil
}

// CreateOrUpdate will deploy or update all the manifets
// the deafultNamespace is used if not otherwise specified.
func (ca Instance) CreateOrUpdate(ctx context.Context, defaultNamepsace string) error {
	for _, res := range ca.PreDeleteResources {
		getObjMetaAndSetDefaultNamespace(res, defaultNamepsace)
		// Create / update / replace resources as required
		if err := kubernetes.DeleteIfExists(ctx, ca.client, res); err != nil {
			return err
		}
	}
	for _, res := range ca.Resources {
		objMeta := getObjMetaAndSetDefaultNamespace(res, defaultNamepsace)
		// do not affect original object (we don't want to affect redeploy in loop)
		resCopy := res.DeepCopyObject()
		if err := waitOnKindDeploy(ctx, ca.client, resCopy); err != nil {
			return fmt.Errorf(
				"can not deploy %s of kind %s to namespace %s - %s",
				objMeta.Name,
				res.GetObjectKind().GroupVersionKind().Kind,
				objMeta.Namespace,
				err)
		}
	}
	ca.Component.Status = kcore.PendingStatus
	ca.Component.Message = "Deployment started"

	return nil
}

// WaitForReadyOrTimeout will wait a reasonable time (defined in context) until a resource is ready
// if the resource become ready, it will update the channel with the component (and status)
// if there is any error or a timeout, it will update the channel with the details on the component
func (ca Instance) WaitForReadyOrTimeout(ctx context.Context, respond chan<- *kcore.Component, wg *sync.WaitGroup) {
	defer wg.Done()

	// here we handle channels and wait groups not errors so pass the timeout context on:
	if err := ca.waitOnApplicationStatus(ctx); err != nil {
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
	metaObj := getObjMeta(ca.ApplicationObject)
	return metaObj.Name
}

// waitOnStatus manages a timeout context when getting application status
func (ca Instance) waitOnApplicationStatus(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Debugf("context waiting for '%s' timed out", ca.Component.Name)
			// we just accept the last status - it's not an error

			return nil
		default:
		}
		err := func() error {
			if err := ca.getStatus(ctx); err != nil {
				log.Debugf("error getting status for %s - %s", ca.Component.Name, err)

				return err
			}

			return nil
		}()
		if err == nil {
			if ca.Component.Status == korev1.SuccessStatus {
				return nil
			}
			// keep waiting
		}
		log.Debugf("not ready so waiting for %s", ca.Component.Name)
		time.Sleep(10 * time.Second)
	}
}

//getStatus will update the ca.component.status from the ca.ApplicationObject conditions
func (ca Instance) getStatus(ctx context.Context) (err error) {
	if ca.ApplicationObject == nil {
		if ca.Component.Name == ClusterAppControllerComponentName {
			if getAppControllerStatus() {
				ca.Component.Message = "Application controller is operational"
				ca.Component.Status = korev1.SuccessStatus
			} else {
				ca.Component.Message = "Status unknown until another application has been deployed"
				ca.Component.Status = korev1.Unknown
			}
		} else {
			ca.Component.Detail = "no application kind created for this component"
			ca.Component.Message = "Component is not checked directly"
			ca.Component.Status = korev1.Unknown
		}
	} else {
		// we need to check if the application CRD exists and get it's status'
		// TODO uses kubebuilder client to get application type and resolve data...
		// First pass just return if object exists?
		us, err := toUnstructuredObj(ca.ApplicationObject)
		if err != nil {

			return fmt.Errorf("error trying to create an unstructured object from the application kind - %s", err)

		}
		log.Debugf("attempting to get status for %s", ca.GetApplicationObjectName())
		exists, err := kubernetes.GetIfExists(ctx, ca.client, us)
		if err != nil {

			return fmt.Errorf(
				"error trying to get %s - %s",
				ca.Component.Name,
				err,
			)

		}
		if !exists {
			log.Debugf("attempting to get status for %s", ca.ApplicationObject)
			ca.Component.Status = korev1.Unknown
			ca.Component.Message = "Application status has not been created"
			ca.Component.Detail = "The application kind"

			return nil
		}
		// Marshall unstructure object back to application kind
		app, err := fromUnstructuredApplication(us)
		if err != nil {

			return fmt.Errorf("error when trying to create an application crd object from an untrustured type - %s", err)

		}
		for _, condition := range app.Status.Conditions {
			setAppControllerStatus(true)
			if condition.Type == applicationv1beta.Ready {
				if condition.Status == "True" {
					ca.Component.Status = korev1.SuccessStatus
					ca.Component.Message = condition.Message

					// All good
					return nil
				}
			}
			if condition.Type == applicationv1beta.Error {
				if condition.Status == "True" {
					// Overright any possible good status
					ca.Component.Status = korev1.FailureStatus
					ca.Component.Message = condition.Message
				}
			}
		}
	}

	return nil
}
