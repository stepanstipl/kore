/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

import (
	"net/http"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/hub"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&teamHandler{})
}

type teamHandler struct {
	hub.Interface
	// DefaultHandlder implements default features
	DefaultHandler

	base string
}

func (u teamHandler) BaseURI() string {
	return u.base
}

// Register is called by the api server to register the service
func (u *teamHandler) Register(i hub.Interface, builder utils.PathBuilder) (*restful.WebService, error) {
	u.Interface = i

	log.WithFields(log.Fields{
		"path": builder.Path("teams"),
	}).Info("registering the teams webservice with container")

	ws := &restful.WebService{}
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path(builder.Path("teams"))

	ws.Route(
		ws.PUT("/invitation/{token}").To(u.invitationSubmit).
			Doc("Used to verify and handle the team invitation generated links").
			Param(ws.PathParameter("token", "The generated base64 invitation token which was provided from the team")).
			Returns(http.StatusOK, "Indicates the generated link is valid and the user has been granted access", nil).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("").To(u.listTeams).
			Doc("Returns all the teams in the hub").
			Returns(http.StatusOK, "A list of all the teams in the hub", orgv1.TeamList{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}").To(u.findTeam).
			Doc("Return information related to the specific team in the hub").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the team definintion from the hub", orgv1.Team{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}").To(u.updateTeam).
			Doc("Used to create or update a team in the hub").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Reads(orgv1.Team{}, "Contains the definition for a team in the hub").
			Returns(http.StatusOK, "Contains the team definintion from the hub", orgv1.Team{}).
			Returns(http.StatusNotModified, "Indicates the request was processed but no changes applied", orgv1.Team{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}").To(u.deleteTeam).
			Doc("Used to delete a team from the hub").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the hub", orgv1.Team{}).
			Returns(http.StatusNotAcceptable, "Indicates you cannot delete the team for one of more reasons", Error{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// Team Members

	ws.Route(
		ws.GET("/{team}/members").To(u.findTeamMembers).
			Doc("Returns a list of user memberships in the team").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains a collection of team memberships for this team", orgv1.TeamMemberList{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/members/{user}").To(u.addTeamMember).
			Doc("Used to add a user to the team via membership").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("user", "Is the user you are adding to the team")).
			Reads(orgv1.TeamMember{}, "The definition for the user in the team").
			Returns(http.StatusOK, "The user has been successfully added to the team", orgv1.TeamMember{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/members/{user}").To(u.removeTeamMember).
			Doc("Used to remove team membership from the team").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("user", "Is the user you are removing to the team")).
			Returns(http.StatusOK, "The user has been successfully removed from the team", orgv1.TeamMember{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// Team Invitations

	ws.Route(
		ws.GET("/{team}/invites/user").To(u.listInvites).
			Doc("Used to return a list of all the users whom have pending invitations").
			Param(ws.PathParameter("team", "The name of the team you are pulling the invitations for")).
			Returns(http.StatusOK, "A list of users whom have an invitation for the team", orgv1.TeamInvitationList{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/invites/user/{user}").To(u.inviteUser).
			Doc("Used to create an invitation for the team").
			Param(ws.PathParameter("team", "The name of the team you are creating an invition")).
			Param(ws.PathParameter("user", "The name of the username of the user the invitation is for")).
			Param(ws.QueryParameter("expire", "The expiration of the generated link").DefaultValue("1h")).
			Returns(http.StatusOK, "Indicates the team invitation for the user has been successful", orgv1.TeamInvitation{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/invites/user/{user}").To(u.removeInvite).
			Doc("Used to remove a user invitation for the team").
			Param(ws.PathParameter("team", "The name of the team you are deleting the invitation")).
			Param(ws.PathParameter("user", "The username of the user whos invitation you are removing")).
			Returns(http.StatusOK, "Indicates the team invitation has been successful removed", nil).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// Invitation Links

	ws.Route(
		ws.GET("/{team}/invites/generate").To(u.inviteLink).
			Doc("Used to generate a link which provides automatic membership of the team").
			Param(ws.PathParameter("team", "The name of the team you are creating an invition link")).
			Param(ws.QueryParameter("expire", "The expiration of the generated link").DefaultValue("1h")).
			Returns(http.StatusOK, "A generated URI which can be used to join a team", "").
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/invites/generate/{user}").To(u.inviteLinkByUser).
			Doc("Used to generate for a specific user to join a team").
			Param(ws.PathParameter("team", "The name of the team you are creating an invition link")).
			Param(ws.PathParameter("user", "The username of the user the link should be limited for")).
			Returns(http.StatusOK, "A generated URI which users can use to join the team", "").
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// Team Allocations

	ws.Route(
		ws.GET("/{team}/allocations").To(u.findAllocations).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.QueryParameter("assigned", "Retrieves all allocations which have been assigned to you")).
			Doc("Used to return a list of all the allocations in the team").
			Returns(http.StatusOK, "Contains the former team definition from the hub", configv1.AllocationList{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.GET("/{team}/allocations/{name}").To(u.findAllocation).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to return")).
			Doc("Used to return all team resources under the team").
			Returns(http.StatusOK, "Contains the former team definition from the hub", configv1.Allocation{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.PUT("/{team}/allocations/{name}").To(u.updateAllocation).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to update")).
			Doc("Used to return all team resources under the team").
			Returns(http.StatusOK, "Contains the former team definition from the hub", configv1.Allocation{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.DELETE("/{team}/allocations/{name}").To(u.deleteAllocation).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to delete")).
			Doc("Used to return all team resources under the team").
			Returns(http.StatusOK, "Contains the former team definition from the hub", configv1.Allocation{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// Team Clusters

	ws.Route(
		ws.GET("/{team}/clusters").To(u.findClusters).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Used to return all team resources under the team").
			Returns(http.StatusOK, "Contains the former team definition from the hub", clustersv1.KubernetesList{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/clusters/{name}").To(u.findCluster).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the kubernetes cluster you are acting upon")).
			Doc("Used to return the cluster definition from the hub").
			Returns(http.StatusOK, "Contains the former team definition from the hub", clustersv1.Kubernetes{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/clusters/{name}").To(u.updateCluster).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the kubernetes cluster you are acting upon")).
			Doc("Used to return all team resources under the team").
			Reads(clustersv1.Kubernetes{}, "The definition for kubernetes cluster").
			Returns(http.StatusOK, "Contains the former team definition from the hub", clustersv1.Kubernetes{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/clusters/{name}").To(u.deleteCluster).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Used to return the cluster definition from the hub").
			Returns(http.StatusOK, "Contains the former team definition from the hub", clustersv1.Kubernetes{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// Team Cloud Providers

	ws.Route(
		ws.GET("/{team}/gkes").To(u.findGKEs).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKEList{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/gkes/{name}").To(u.findGKE).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKE{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/gkes/{name}").To(u.updateGKE).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Doc("Is used to provision or update a GKE cluster in the hub").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKE{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/gkes/{name}").To(u.deleteGKE).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Doc("Is used to delete a managed GKE cluster from the hub").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKE{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	// GKE Credentials - @TODO these all need to be autogenerated

	ws.Route(
		ws.GET("/{team}/gkecredentials").To(u.findGKECredientalss).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKECredentialsList{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/gkecredentials/{name}").To(u.findGKECredientals).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKECredentials{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/gkecredentials/{name}").To(u.updateGKECredientals).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Doc("Is used to provision or update a GKE cluster in the hub").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKECredentials{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/gkecredentials/{name}").To(u.deleteGKECredientals).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Doc("Is used to delete a managed GKE cluster from the hub").
			Returns(http.StatusOK, "Contains the former team definition from the hub", gke.GKECredentials{}).
			DefaultReturns("An generic API error containing the cause of the error", Error{}),
	)

	return ws, nil
}

// Name returns the name of the handler
func (u teamHandler) Name() string {
	return "teams"
}

// Teams Management

// deleteTeam is responsible for deleting a team from the hub
func (u teamHandler) deleteTeam(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		err := u.Teams().Delete(req.Request.Context(), req.PathParameter("team"))
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

// listTeams returns all the teams in the hub
func (u teamHandler) listTeams(req *restful.Request, resp *restful.Response) {
	handleErrors(req, resp, func() error {
		teams, err := u.Teams().List(req.Request.Context())
		if err != nil {
			return err
		}

		return resp.WriteHeaderAndEntity(http.StatusOK, teams)
	})
}

// updateTeam is responsible for updating for creating a team in the hub
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
