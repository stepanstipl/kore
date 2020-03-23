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
	"time"

	"github.com/appvia/kore/pkg/apiserver/types"

	clustersv1 "github.com/appvia/kore/pkg/apis/clusters/v1"
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	eks "github.com/appvia/kore/pkg/apis/eks/v1alpha1"
	gcp "github.com/appvia/kore/pkg/apis/gcp/v1alpha1"
	gke "github.com/appvia/kore/pkg/apis/gke/v1alpha1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/kore"
	"github.com/appvia/kore/pkg/kore/validation"
	"github.com/appvia/kore/pkg/utils"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterHandler(&teamHandler{})
}

type teamHandler struct {
	kore.Interface
	// DefaultHandlder implements default features
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

	ws.Route(
		ws.PUT("/invitation/{token}").To(u.invitationSubmit).
			Doc("Used to verify and handle the team invitation generated links").
			Operation("InvitationSubmit").
			Param(ws.PathParameter("token", "The generated base64 invitation token which was provided from the team")).
			Returns(http.StatusOK, "Indicates the generated link is valid and the user has been granted access", types.TeamInvitationResponse{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("").To(u.listTeams).
			Doc("Returns all the teams in the kore").
			Operation("ListTeams").
			Returns(http.StatusOK, "A list of all the teams in the kore", orgv1.TeamList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}").To(u.findTeam).
			Doc("Return information related to the specific team in the kore").
			Operation("GetTeam").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the team definintion from the kore", orgv1.Team{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}").To(u.updateTeam).
			Doc("Used to create or update a team in the kore").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Operation("UpdateTeam").
			Reads(orgv1.Team{}, "Contains the definition for a team in the kore").
			Returns(http.StatusOK, "Contains the team definintion from the kore", orgv1.Team{}).
			Returns(http.StatusNotModified, "Indicates the request was processed but no changes applied", orgv1.Team{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}").To(u.deleteTeam).
			Doc("Used to delete a team from the kore").
			Operation("RemoveTeam").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", orgv1.Team{}).
			Returns(http.StatusNotAcceptable, "Indicates you cannot delete the team for one or more reasons", Error{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Team Audit Events

	ws.Route(
		ws.GET("/{team}/audit").To(u.findTeamAudit).
			Doc("Used to return a collection of events against the team").
			Operation("GetTeamAudit").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.QueryParameter("since", "The duration to retrieve from the audit log").DefaultValue("60m")).
			Returns(http.StatusOK, "A collection of audit events against the team", orgv1.AuditEventList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Team Members

	ws.Route(
		ws.GET("/{team}/members").To(u.findTeamMembers).
			Doc("Returns a list of user memberships in the team").
			Operation("GetTeamMembers").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains a collection of team memberships for this team", orgv1.TeamMemberList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/members/{user}").To(u.addTeamMember).
			Doc("Used to add a user to the team via membership").
			Operation("AddTeamMember").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("user", "Is the user you are adding to the team")).
			Reads(orgv1.TeamMember{}, "The definition for the user in the team").
			Returns(http.StatusOK, "The user has been successfully added to the team", orgv1.TeamMember{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/members/{user}").To(u.removeTeamMember).
			Doc("Used to remove team membership from the team").
			Operation("RemoveTeamMember").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("user", "Is the user you are removing from the team")).
			Returns(http.StatusOK, "The user has been successfully removed from the team", orgv1.TeamMember{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Team Invitations

	ws.Route(
		ws.GET("/{team}/invites/user").To(u.listInvites).
			Doc("Used to return a list of all the users whom have pending invitations").
			Operation("ListInvites").
			Param(ws.PathParameter("team", "The name of the team you are pulling the invitations for")).
			Returns(http.StatusOK, "A list of users whom have an invitation for the team", orgv1.TeamInvitationList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/invites/user/{user}").To(u.inviteUser).
			Doc("Used to create an invitation for the team").
			Operation("InviteUser").
			Param(ws.PathParameter("team", "The name of the team you are creating an invition")).
			Param(ws.PathParameter("user", "The name of the username of the user the invitation is for")).
			Param(ws.QueryParameter("expire", "The expiration of the generated link").DefaultValue("1h")).
			Returns(http.StatusOK, "Indicates the team invitation for the user has been successful", nil).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/invites/user/{user}").To(u.removeInvite).
			Doc("Used to remove a user invitation for the team").
			Operation("RemoveInvite").
			Param(ws.PathParameter("team", "The name of the team you are deleting the invitation")).
			Param(ws.PathParameter("user", "The username of the user whos invitation you are removing")).
			Returns(http.StatusOK, "Indicates the team invitation has been successful removed", nil).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Invitation Links

	ws.Route(
		ws.GET("/{team}/invites/generate").To(u.inviteLink).
			Doc("Used to generate a link which provides automatic membership of the team").
			Operation("GenerateInviteLink").
			Param(ws.PathParameter("team", "The name of the team you are creating an invition link")).
			Param(ws.QueryParameter("expire", "The expiration of the generated link").DefaultValue("1h")).
			Returns(http.StatusOK, "A generated URI which can be used to join a team", "").
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/invites/generate/{user}").To(u.inviteLinkByUser).
			Doc("Used to generate for a specific user to join a team").
			Operation("GenerateInviteLinkForUser").
			Param(ws.PathParameter("team", "The name of the team you are creating an invition link")).
			Param(ws.PathParameter("user", "The username of the user the link should be limited for")).
			Returns(http.StatusOK, "A generated URI which users can use to join the team", "").
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Team Allocations

	ws.Route(
		ws.GET("/{team}/allocations").To(u.findAllocations).
			Doc("Used to return a list of all the allocations in the team").
			Operation("ListAllocations").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.QueryParameter("assigned", "Retrieves all allocations which have been assigned to you")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.AllocationList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.GET("/{team}/allocations/{name}").To(u.findAllocation).
			Doc("Used to return an allocation within the team").
			Operation("GetAllocation").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to return")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.Allocation{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.PUT("/{team}/allocations/{name}").To(u.updateAllocation).
			Doc("Used to create/update an allocation within the team.").
			Operation("UpdateAllocation").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to update")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.Allocation{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.DELETE("/{team}/allocations/{name}").To(u.deleteAllocation).
			Doc("Remove an allocation from a team").
			Operation("RemoveAllocation").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the allocation you wish to delete")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", configv1.Allocation{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Team Namespaces

	ws.Route(
		ws.GET("/{team}/namespaceclaims").To(u.findNamespaces).
			Doc("Used to return all namespaces for the team").
			Operation("ListNamespaces").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former definition from the kore", clustersv1.NamespaceClaimList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/namespaceclaims/{name}").To(u.findNamespace).
			Doc("Used to return the details of a namespace within a team").
			Operation("GetNamespace").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the namespace claim you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", clustersv1.NamespaceClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/namespaceclaims/{name}").To(u.updateNamespace).
			Doc("Used to create or update the details of a namespace within a team").
			Operation("UpdateNamespace").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the namespace claim you are acting upon")).
			Reads(clustersv1.NamespaceClaim{}, "The definition for namespace claim").
			Returns(http.StatusOK, "Contains the definition from the kore", clustersv1.NamespaceClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/namespaceclaims/{name}").To(u.deleteNamespace).
			Doc("Used to remove a namespace from a team").
			Operation("RemoveNamespace").
			Param(ws.PathParameter("name", "Is name the of the namespace claim you are acting upon")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former definition from the kore", clustersv1.NamespaceClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Secrets is used to provision a secret in the team

	ws.Route(
		ws.GET("/{team}/secrets").To(u.findTeamSecrets).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Used to return all the secrets within the team").
			Returns(http.StatusOK, "Contains the definition for the resource", configv1.Secret{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/secrets/{name}").To(u.findTeamSecret).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the secert in the name")).
			Doc("Used to retrieve the secret from the team").
			Returns(http.StatusOK, "Contains the definition for the resource", configv1.Secret{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/secrets/{name}").To(u.updateTeamSecret).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of secret you are creating / updating")).
			Doc("Used to update the secret in the team").
			Reads(configv1.Secret{}, "The definition for the secret you are creating or updating").
			Returns(http.StatusOK, "Contains updated definition of the secret", configv1.Secret{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/secrets/{name}").To(u.deleteTeamSecret).
			Param(ws.PathParameter("name", "Is name the of the secret you are acting upon")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Used to delete the secret from team").
			Returns(http.StatusOK, "Contains the former definition of the secret", configv1.Secret{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Team Clusters

	ws.Route(
		ws.GET("/{team}/clusters").To(u.findClusters).
			Doc("Lists all clusters available for a team").
			Operation("ListClusters").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", clustersv1.KubernetesList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/clusters/{name}").To(u.findCluster).
			Doc("Used to return the cluster definition from the kore").
			Operation("GetCluster").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the kubernetes cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", clustersv1.Kubernetes{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/clusters/{name}").To(u.updateCluster).
			Doc("Used to create/update a cluster definition for a team").
			Operation("UpdateCluster").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the kubernetes cluster you are acting upon")).
			Reads(clustersv1.Kubernetes{}, "The definition for kubernetes cluster").
			Returns(http.StatusOK, "Contains the former team definition from the kore", clustersv1.Kubernetes{}).
			Returns(http.StatusBadRequest, "Validation error of the provided details", validation.ErrValidation{}). // @TODO: Change this to be a class in the orgv1 space
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/clusters/{name}").To(u.deleteCluster).
			Doc("Used to remove a cluster from a team").
			Operation("RemoveCluster").
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", clustersv1.Kubernetes{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// Team Cloud Providers

	// GKE Clusters

	ws.Route(
		ws.GET("/{team}/gkes").To(u.findGKEs).
			Doc("Returns a list of Google Container Engine clusters which the team has access").
			Operation("ListGKEs").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKEList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/gkes/{name}").To(u.findGKE).
			Doc("Returns a specific Google Container Engine cluster to which the team has access").
			Operation("GetGKE").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKE{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/gkes/{name}").To(u.updateGKE).
			Doc("Is used to provision or update a GKE cluster in the kore").
			Operation("UpdateGKE").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKE{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/gkes/{name}").To(u.deleteGKE).
			Doc("Is used to delete a managed GKE cluster from the kore").
			Operation("RemoveGKE").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKE{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// GKE Credentials - @TODO these all need to be autogenerated

	ws.Route(
		ws.GET("/{team}/gkecredentials").To(u.findGKECredientalss).
			Doc("Returns a list of GKE Credentials to which the team has access").
			Operation("ListGKECredentials").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentialsList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/gkecredentials/{name}").To(u.findGKECredientals).
			Doc("Returns a specific GKE Credential to which the team has access").
			Operation("GetGKECredential").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/gkecredentials/{name}").To(u.updateGKECredientals).
			Doc("Creates or updates a specific GKE Credential to which the team has access").
			Operation("UpdateGKECredential").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/gkecredentials/{name}").To(u.deleteGKECredientals).
			Doc("Deletes a specific GKE Credential from the team").
			Operation("RemoveGKECredential").
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Returns(http.StatusOK, "Contains the former team definition from the kore", gke.GKECredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// GCP Project Claims

	ws.Route(
		ws.GET("/{team}/projectclaims").To(u.findProjectClaims).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.ProjectClaimList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/projectclaims/{name}").To(u.findProjectClaim).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.ProjectClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/projectclaims/{name}").To(u.updateProjectClaim).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is used to provision or update a gcp project claim").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.ProjectClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/projectclaims/{name}").To(u.deleteProjectClaim).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is used to delete a managed gcp project claim").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.ProjectClaim{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// GCP Organization

	ws.Route(
		ws.GET("/{team}/organizations").To(u.findOrganizations).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of gcp organizations").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.OrganizationList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/organizations/{name}").To(u.findOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is the used tor return a specific gcp organization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.Organization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)
	ws.Route(
		ws.PUT("/{team}/organizations/{name}").To(u.updateOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is used to provision or update a gcp organization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.Organization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/organizations/{name}").To(u.deleteOrganization).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the resource you are acting on")).
			Doc("Is used to delete a managed gcp organization").
			Returns(http.StatusOK, "Contains the former team definition from the kore", gcp.Organization{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// EKS clusters

	ws.Route(
		ws.GET("/{team}/ekss").To(u.findEKSs).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Amazon EKS clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/ekss/{name}").To(u.findEKS).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS cluster you are acting upon")).
			Doc("Is the used to return a EKS cluster which the team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKS{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/ekss/{name}").To(u.updateEKS).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS cluster you are acting upon")).
			Doc("Is used to provision or update a EKS cluster in the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKS{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/ekss/{name}").To(u.deleteEKS).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS cluster you are acting upon")).
			Doc("Is used to delete a managed EKS cluster from the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKS{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// EKS Nodegroups
	ws.Route(
		ws.GET("/{team}/eksnodegroups").To(u.findEKSNodeGroups).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Amazon EKS clusters which the team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSNodeGroupList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/eksnodegroups/{name}").To(u.findEKSNodeGroup).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the EKS nodegroup")).
			Doc("Is the used to return a EKS cluster which the team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSNodeGroup{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/eksnodegroups/{name}").To(u.updateEKSNodeGroups).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the EKS nodegroup")).
			Doc("Is used to provision or update a EKS cluster in the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSNodeGroup{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/eksnodegroups/{name}").To(u.deleteEKSNodeGroups).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is the name of the EKS nodegroup")).
			Doc("Is used to delete a managed EKS cluster nodegroup from the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSNodeGroup{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	// EKS Credentials - @TODO these all need to be autogenerated

	ws.Route(
		ws.GET("/{team}/ekscredentials").To(u.findEKSCredentialss).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Doc("Is the used tor return a list of Amazon EKS clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentialsList{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.GET("/{team}/ekscredentials/{name}").To(u.findEKSCredentails).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the GKE cluster you are acting upon")).
			Doc("Is the used tor return a list of Google Container Engine clusters which thhe team has access").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.PUT("/{team}/ekscredentials/{name}").To(u.updateEKSCredentails).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS cluster you are acting upon")).
			Doc("Is used to provision or update a EKS cluster in the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

	ws.Route(
		ws.DELETE("/{team}/ekscredentials/{name}").To(u.deleteEKSCredentails).
			Param(ws.PathParameter("team", "Is the name of the team you are acting within")).
			Param(ws.PathParameter("name", "Is name the of the EKS cluster you are acting upon")).
			Doc("Is used to delete a managed EKS cluster from the kore").
			Returns(http.StatusOK, "Contains the former team definition from the kore", eks.EKSCredentials{}).
			DefaultReturns("A generic API error containing the cause of the error", Error{}),
	)

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

		list, err := u.Teams().Team(team).AuditEvents(req.Request.Context(), tm)
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
