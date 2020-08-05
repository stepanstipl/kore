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
	"fmt"
	"net/http"

	aws "github.com/appvia/kore/pkg/apis/aws/v1alpha1"
	"github.com/appvia/kore/pkg/kore"

	restful "github.com/emicklei/go-restful"
)

const (
	// HeaderAWSSecretAccessKey is the header name used when passing the AWS secret access key credential to the API
	HeaderAWSSecretAccessKey = "x-api-aws-secret-access-key"
	// HeaderAWSAccessKeyID is the header name used when passing AWS access key id credential to the API
	HeaderAWSAccessKeyID = "x-api-aws-access-key-id"
)

func (u teamHandler) addAWSOrganizationRoutes(ws *restful.WebService) {
	ws.Route(
		ws.GET("/{team}/awsorganizations").To(u.findAWSOrganizations).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of aws organizations").
			Operation("ListAWSOrganizations").
			Returns(http.StatusOK, "Contains the former team definition from the kore", aws.AWSOrganizationList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/awsorganizations/{name}").To(u.findAWSOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is the used tor return a specific aws organization").
			Operation("GetAWSOrganization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", aws.AWSOrganization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.PUT("/{team}/awsorganizations/{name}").To(u.updateAWSOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Operation("UpdateAWSOrganization").
			Reads(aws.AWSOrganization{}, "The definition for AWS organization").
			Doc("Is used to provision or update a aws organization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", aws.AWSOrganization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/awsorganizations/{name}").To(u.deleteAWSOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is used to delete a managed gcp organization").
			Operation("DeleteAWSOrganization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", aws.AWSOrganization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// aws OU's (Organizational Units) available for accounts
	ws.Route(
		ws.GET("/{team}/awsorganizations/awsAccountOUs").To(u.getAccountOUs).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.QueryParameter("region", "Is the region where Control Tower and AWS Organizations are enabled")).
			Param(ws.QueryParameter("roleARN", "Is the role in the master account to use for querying the Organization")).
			Param(ws.HeaderParameter(HeaderAWSSecretAccessKey, "Is the AWS Secret Access Key used for authenticating to the AWS API")).
			Param(ws.HeaderParameter(HeaderAWSAccessKeyID, "Is the AWS Access Key ID used for authenticating to the AWS API")).
			Doc("Gets a list of AWS OU's suitable for creating aws accounts within").
			Notes(fmt.Sprintf("Requires the authentication headers %s, and %s", HeaderAWSSecretAccessKey, HeaderAWSAccessKeyID)).
			Operation("ListAWSAccountOUs").
			Returns(http.StatusOK, "Account OUs", []string{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
}

// findOrganizations returns a list of credential claims
func (u teamHandler) findAWSOrganizations(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		list, err := u.Teams().Team(team).
			Cloud().AWS().AWSOrganizations().
			List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// findOrganizations returns a specific credential
func (u teamHandler) findAWSOrganization(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		n, err := u.Teams().Team(team).
			Cloud().AWS().AWSOrganizations().
			Get(req.Request.Context(), name)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// updateOrganization is used to update an credential claim for a team
func (u teamHandler) updateAWSOrganization(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		claim := &aws.AWSOrganization{}
		if err := req.ReadEntity(claim); err != nil {
			return err
		}

		n, err := u.Teams().Team(team).
			Cloud().AWS().AWSOrganizations().
			Update(req.Request.Context(), claim)

		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, n)
	})
}

// deleteOrganization is used to remove a credential from a team cluster
func (u teamHandler) deleteAWSOrganization(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		name := req.PathParameter("name")

		original, err := u.Teams().Team(team).
			Cloud().AWS().AWSOrganizations().
			Delete(req.Request.Context(), name)

		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, original)
	})
}

// getAccountOUs will retrieve a list of OU's that can have account provisioned
func (u teamHandler) getAccountOUs(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")
		params := kore.OUParams{
			AWSAccessKeyID:     req.HeaderParameter(HeaderAWSAccessKeyID),
			AWSSecretAccessKey: req.HeaderParameter(HeaderAWSSecretAccessKey),
			Region:             req.QueryParameter("region"),
			RoleARN:            req.QueryParameter("roleARN"),
		}

		ouList, err := u.Teams().Team(team).Cloud().AWS().AWSOrganizations().GetAccountOUs(
			req.Request.Context(),
			params,
		)

		if err != nil {

			return err
		}
		return resp.WriteHeaderAndEntity(http.StatusOK, ouList)
	})
}
