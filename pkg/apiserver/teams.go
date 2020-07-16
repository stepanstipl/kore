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
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/appvia/kore/pkg/apiserver/params"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&teamHandler{})
}

type teamHandler struct {
	kore.Interface
	// DefaultHandler implements default features
	DefaultHandler
}

// Register is called by the api server to register the service
func (u *teamHandler) Register(i kore.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	u.Interface = i
	path := builder.Path("teams")

	log.WithFields(log.Fields{
		"path": path,
	}).Info("registering the teams webservice with container")

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(path)

	ws.Filter(func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		team := req.PathParameter("team")
		if team == "" {
			chain.ProcessFilter(req, resp)
			return
		}

		// Team resource endpoints do this check themselves
		if strings.HasSuffix(req.Request.RequestURI, fmt.Sprintf("teams/%s", team)) {
			chain.ProcessFilter(req, resp)
			return
		}

		exists, err := u.Teams().Exists(context.Background(), team)
		if err != nil {
			handleError(req, resp, err)
			return
		}
		if !exists {
			writeError(req, resp, fmt.Errorf("team %q does not exist", team), http.StatusNotFound)
			return
		}

		chain.ProcessFilter(req, resp)
	})

	ws.Route(
		ws.GET("").To(u.listTeams).
			Doc("Returns all the teams in the kore").
			Operation("ListTeams").
			Returns(http.StatusOK, "A list of all the teams in the kore", orgv1.TeamList{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}").To(u.findTeam).
			Doc("Return information related to the specific team in the kore").
			Operation("GetTeam").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the team definition from the kore", orgv1.Team{}).
			Returns(http.StatusNotFound, "Team does not exist", nil).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		withAllErrors(ws.PUT("/{team}")).To(u.updateTeam).
			Doc("Used to create or update a team in the kore").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Operation("UpdateTeam").
			Reads(orgv1.Team{}, "Contains the definition for a team in the kore").
			Returns(http.StatusOK, "Contains the team definition from the kore", orgv1.Team{}).
			Returns(http.StatusNotFound, "Team does not exist", nil).
			Returns(http.StatusNotModified, "Indicates the request was processed but no changes applied", orgv1.Team{}),
	)

	ws.Route(
		ws.DELETE("/{team}").To(u.deleteTeam).
			Doc("Used to delete a team from the kore").
			Operation("RemoveTeam").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(params.DeleteCascade()).
			Returns(http.StatusOK, "Contains the former team definition from the kore", orgv1.Team{}).
			Returns(http.StatusNotFound, "Team does not exist", nil).
			Returns(http.StatusNotAcceptable, "Indicates you cannot delete the team for one or more reasons", Error{}).
			Returns(http.StatusInternalServerError, "A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/audits")).To(u.findTeamAudit).
			Doc("Used to return a collection of events against the team").
			Operation("ListTeamAudit").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.QueryParameter("since", "The duration to retrieve from the audit log").DefaultValue("60m")).
			Returns(http.StatusOK, "A collection of audit events against the team", orgv1.AuditEventList{}).
			Returns(http.StatusNotFound, "Team does not exist", nil),
	)

	ws.Route(
		withAllNonValidationErrors(ws.GET("/{team}/plans/{plan}")).To(u.getTeamPlanDetails).
			Operation("GetTeamPlanDetails").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("plan", "Is name the of the plan you're interested in")).
			Doc("Returns the plan, the JSON schema of the plan, and what what parameters are allowed to be edited by this team when using the plan").
			Returns(http.StatusOK, "Contains details of the plan", TeamPlan{}).
			Returns(http.StatusNotFound, "Team or plan doesn't exist", nil),
	)

	u.addTeamMemberRoutes(ws)
	u.addInvitationRoutes(ws)
	u.addAllocationsRoutes(ws)
	u.addNamespaceRoutes(ws)
	u.addTeamSecretRoutes(ws)
	u.addKubernetesRoutes(ws)
	u.addClusterRoutes(ws)

	// GKE/GCP
	u.addGKERoutes(ws)
	u.addGKECredentialsRoutes(ws)
	u.addGCPProjectClaimRoutes(ws)
	u.addGCPOrganizationRoutes(ws)

	// EKS/AWS
	u.addAWSAccountClaimRoutes(ws)
	u.addAWSOrganizationRoutes(ws)
	u.addEKSRoutes(ws)
	u.addEKSNodeGroupRoutes(ws)
	u.addEKSCredentialsRoutes(ws)
	u.addEKSVPCRoutes(ws)

	// AKS/Azure
	u.addAKSRoutes(ws)
	u.addAKSCredentialsRoutes(ws)

	// Team services
	u.addServiceRoutes(ws)
	u.addServiceCredentialRoutes(ws)

	// Team Security Reports
	u.addTeamSecurityRoutes(ws)

	return ws, nil
}

// Name returns the name of the handler
func (u teamHandler) Name() string {
	return "teams"
}

// findTeamAudit returns the audit log for a team
func (u teamHandler) findTeamAudit(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := req.PathParameter("team")

		since := req.QueryParameter("since")
		if since == "" {
			since = "60m"
		}
		tm, err := time.ParseDuration(since)
		if err != nil {
			return err
		}

		list, err := u.Audit().AuditEventsTeam(req.Request.Context(), team, tm)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, list)
	})
}

// Teams Management

// deleteTeam is responsible for deleting a team from the kore
func (u teamHandler) deleteTeam(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		err := u.Teams().Delete(req.Request.Context(), req.PathParameter("team"), parseDeleteOpts(req)...)
		if err != nil {
			return err
		}
		resp.WriteHeader(http.StatusOK)

		return nil
	})
}

// findTeam returns a specific team
func (u teamHandler) findTeam(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team, err := u.Teams().Get(req.Request.Context(), req.PathParameter("team"))
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, team)
	})
}

// listTeams returns all the teams in the kore
func (u teamHandler) listTeams(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		teams, err := u.Teams().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, teams)
	})
}

// updateTeam is responsible for updating for creating a team in the kore
func (u teamHandler) updateTeam(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		team := &orgv1.Team{}
		if err := req.ReadEntity(team); err != nil {
			return err
		}
		team, err := u.Teams().Update(req.Request.Context(), team)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, team)
	})
}

func (u teamHandler) getTeamPlanDetails(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		plan, err := u.Plans().Get(req.Request.Context(), req.PathParameter("plan"))
		if err != nil {
			return err
		}

		clusterProvider, exists := kore.GetClusterProvider(plan.Spec.Kind)
		if !exists {
			writeError(req, resp, fmt.Errorf("unknown cluster provider type %q", plan.Spec.Kind), http.StatusNotFound)
			return nil
		}

		editableParams, err := u.Plans().GetEditablePlanParams(req.Request.Context(), req.PathParameter("team"), plan.Spec.Kind)
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, TeamPlan{
			Plan:           plan.Spec,
			Schema:         clusterProvider.PlanJSONSchema(),
			EditableParams: editableParams,
		})
	})
}
