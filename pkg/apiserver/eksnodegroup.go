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

package apiserver

import (
	"net/http"

	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"

	restful "github.com/emicklei/go-restful"
)

// findEKSNodegroups returns all the nodegroups for a EKS cluster for a team
func (u teamHandler) findEKSNodeGroups(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		list, err := u.Teams().Team(team).Cloud().EKSNodeGroup().List(req.Request.Context())
		if err != nil {
			return err
		}
		// filter list by just the relevant cluster
		newList := list.DeepCopy()
		newItems := make([]eks.EKSNodeGroup, 0)
		for _, ng := range list.Items {
			if ng.Name == name {
				newItems = append(newItems, ng)
			}
		}
		newList.Items = newItems

		return resp.WriteHeaderAndEntity(http.StatusOK, newList)
	})
}

// findEKS returns a specific nodegroup for a cluster under the team
func (u teamHandler) findEKSNodeGroup(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		ng, err := u.Teams().Team(team).Cloud().EKSNodeGroup().Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, ng)
	})
}

// deleteEKS is responsible for deleting an eksnodegroup resource
func (u teamHandler) deleteEKSNodeGroups(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		ctx := req.Request.Context()
		name := req.PathParameter("name")
		team := req.PathParameter("team")

		object, err := u.Teams().Team(team).Cloud().EKSNodeGroup().Get(ctx, name)
		if err != nil {
			return err
		}

		err = u.Teams().Team(team).Cloud().EKSNodeGroup().Delete(ctx, name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}

// updateEKS is responsible for putting an resource into a team
func (u teamHandler) updateEKSNodeGroups(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		object := &eks.EKSNodeGroup{}
		if err := req.ReadEntity(object); err != nil {
			return err
		}

		if _, err := u.Teams().Team(team).Cloud().EKSNodeGroup().Update(req.Request.Context(), object); err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, object)
	})
}
