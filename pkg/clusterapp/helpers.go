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
	"time"

	log "github.com/sirupsen/logrus"

	korev1 "github.com/appvia/kore/pkg/apis/core/v1"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	applicationv1beta "sigs.k8s.io/application/api/v1beta1"
)

func setMissingNamespace(namespace string, obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		log.Debugf("no setting namespace here ->%v<- - %s", obj, err)

		return err
	}
	if accessor.GetNamespace() == "" {
		accessor.SetNamespace(namespace)
		log.Debugf(
			"updated namespace to %s on %s",
			accessor.GetNamespace(),
			accessor.GetName(),
		)
	}

	return nil
}

// waitOnStatus manages a timeout context when getting application status
func waitOnApplicationStatus(ctx context.Context, ca *Instance) error {
	for {
		select {
		case <-ctx.Done():
			log.Debugf("context waiting for '%s' timed out", ca.Component.Name)
			// we just accept the last status - it's not an error

			return nil
		default:
		}
		err := func() error {
			if err := getStatus(ctx, ca); err != nil {
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
func getStatus(ctx context.Context, ca *Instance) (err error) {
	//TODO - implement watcher
	// 1. watcher (separate thread) will watch the kube api for changes to specific "Application" CRD instance
	// update channel with nil (if status success) or with error (including timeout)

	if ca.ApplicationObject == nil {
		// TODO - have to support the application operator itself so need to do something here
		// we have to check the existence of something other...
		// maybe we just have to look for presence of the statefulset???
		ca.Component.Detail = "We have to assume ok as we do not have an application to track"
		ca.Component.Message = "System component not checked"
		ca.Component.Status = korev1.SuccessStatus
		log.Debugf(
			"no application object should only be the case for application controller - component is %s",
			ca.Component.Name,
		)
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

func getObjMeta(obj runtime.Object) (metav1.ObjectMeta, error) {
	metaObj := metav1.ObjectMeta{}
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return metaObj, err
	}
	// TODO: this should be a pointer to the origonal data?
	metaObj.Name = accessor.GetName()
	metaObj.Namespace = accessor.GetNamespace()
	metaObj.Labels = accessor.GetLabels()
	return metaObj, nil
}

func toUnstructuredObj(obj runtime.Object) (*unstructured.Unstructured, error) {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Version: obj.GetObjectKind().GroupVersionKind().Version,
		Kind:    obj.GetObjectKind().GroupVersionKind().Kind,
		Group:   obj.GetObjectKind().GroupVersionKind().Group,
	})
	objMeta, _ := getObjMeta(obj)
	u.SetName(objMeta.Name)
	u.SetNamespace(objMeta.Namespace)
	u.SetLabels(objMeta.Labels)
	return u, nil
}

func fromUnstructuredApplication(us *unstructured.Unstructured) (*applicationv1beta.Application, error) {
	app := &applicationv1beta.Application{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(us.Object, app); err != nil {
		return nil, err
	}
	return app, nil
}
