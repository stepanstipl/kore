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

package namespaces

import (
	"context"
	"fmt"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	"github.com/appvia/kore/pkg/controllers"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/kubernetes"
	v1 "k8s.io/api/core/v1"

	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	// filterNamespaces is a list of namespaces to ignore
	filterNamespaces = []string{
		"application-system",
		"default",
		"flux",
		"kore",
		"kore-operators",
		"kore-system",
		"kube-public",
		"kube-system",
	}
)

// Reconcile is the entrypoint for the reconciliation logic
func (a *Controller) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	ctx := context.Background()
	cc := a.mgr.GetClient()

	logger := a.logger.WithFields(log.Fields{
		"name":      request.Name,
		"namespace": request.Namespace,
	})
	logger.Debug("attempting to reconcile importing of cluster namespaces")

	// @logic
	// - we iterate this clusters and check they are available
	// - we pull in a listing on the namespaces
	// - we check if a namespaceclaim exists and if not create one
	list := &clustersv1.ClusterList{}
	if err := cc.List(ctx, list, client.InNamespace("")); err != nil {
		logger.WithError(err).Error("trying to retrieve a list of clusters")

		return reconcile.Result{}, err
	}

	var errorCounter int

	for i := 0; i < len(list.Items); i++ {
		// @step: for now we bypass any clusters in the kore namespace
		if list.Items[i].Namespace == kore.HubNamespace {
			continue
		}
		l := logger.WithFields(log.Fields{
			"cluster": list.Items[i].Name,
		})

		// create a client to speak to the cluster
		kc, err := controllers.CreateClient(ctx, cc, list.Items[i].Ownership())
		if err != nil {
			l.WithError(err).Error("trying to create cluster client")
			errorCounter++

			continue
		}

		// retrieve a listing of the namespaces
		nlist := &v1.NamespaceList{}
		if err := kc.List(ctx, nlist, client.InNamespace("")); err != nil {
			l.WithError(err).Error("trying to retrieve a listing of namespaces from cluster")
			errorCounter++

			continue
		}

		// iterate the namespaces and check if we have a claim on them
		for k := 0; k < len(nlist.Items); k++ {
			// filter out anything we are not interesting
			if utils.Contains(nlist.Items[k].Name, filterNamespaces) {
				continue
			}

			// check if a namespaceclaim exists for this
			o := &clustersv1.NamespaceClaim{}
			o.Namespace = list.Items[i].Namespace
			// i.e. the <cluster>-<name>
			o.Name = fmt.Sprintf("%s-%s", list.Items[i].Name, nlist.Items[k].Name)

			found, err := kubernetes.CheckIfExists(ctx, cc, o)
			if err != nil {
				l.WithField(
					"namespace", list.Items[i].Name,
				).WithError(err).Error("trying to check if the namespaceclaim exists")
				errorCounter++

				continue
			}
			if found {
				continue
			}

			// else no namespaceclaim exists for the namespace - lets import it
			o.Annotations = map[string]string{
				kore.Label("import"): "true",
			}
			labels := list.Items[i].GetLabels()
			labels["owned"] = "true"

			o.Spec = clustersv1.NamespaceClaimSpec{
				Annotations: list.Items[i].Annotations,
				Cluster:     list.Items[i].Ownership(),
				Labels:      labels,
				Name:        nlist.Items[k].Name,
			}
			if _, err := kubernetes.CreateOrUpdate(ctx, cc, o); err != nil {
				l.WithField(
					"namespace", list.Items[i].Name,
				).WithError(err).Error("trying to create namespaceclaim")
				errorCounter++

				continue
			}
		}
	}
	// if we have any errors we should retry
	if errorCounter > 0 {
		logger.Warn("encountered errors trying to import the namespaces, will retry")
	}

	return reconcile.Result{}, nil
}
