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

package kore

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/security"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// DefaultConvertor is a default type
	DefaultConvertor Convertor
)

type Convertor struct{}

// ToUserModel converts from api to model
func (c Convertor) ToUserModel(user *orgv1.User) *model.User {
	return &model.User{
		Username: user.Spec.Username,
		Email:    user.Spec.Email,
	}
}

// FromAuditModel converts the model
func (c Convertor) FromAuditModel(i *model.AuditEvent) *orgv1.AuditEvent {
	return &orgv1.AuditEvent{
		TypeMeta: metav1.TypeMeta{
			APIVersion: orgv1.SchemeGroupVersion.String(),
			Kind:       "AuditEvent",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              i.Team,
			CreationTimestamp: metav1.NewTime(i.CreatedAt),
		},
		Spec: orgv1.AuditEventSpec{
			ID:           i.ID,
			CreatedAt:    metav1.NewTime(i.CreatedAt),
			Verb:         i.Verb,
			Team:         i.Team,
			User:         i.User,
			Message:      i.Message,
			Resource:     i.Resource,
			ResourceURI:  i.ResourceURI,
			APIVersion:   i.APIVersion,
			Operation:    i.Operation,
			StartedAt:    metav1.NewTime(i.StartedAt),
			CompletedAt:  metav1.NewTime(i.CompletedAt),
			ResponseCode: i.ResponseCode,
		},
	}
}

// FromUserInvitationModel converts the model to a thing
func (c Convertor) FromUserInvitationModel(i *model.Invitation) *orgv1.TeamInvitation {
	return &orgv1.TeamInvitation{
		TypeMeta: metav1.TypeMeta{
			APIVersion: orgv1.SchemeGroupVersion.String(),
			Kind:       "TeamInvitation",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              i.User.Username,
			CreationTimestamp: metav1.NewTime(i.CreatedAt),
		},
		Spec: orgv1.TeamInvitationSpec{
			Username: i.User.Username,
			Team:     i.Team.Name,
		},
	}
}

// FromUserInvitationModelList convert the list of invitations
func (c Convertor) FromUserInvitationModelList(list []*model.Invitation) *orgv1.TeamInvitationList {
	l := &orgv1.TeamInvitationList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "TeamInvitationList",
		},
		Items: make([]orgv1.TeamInvitation, len(list)),
	}
	length := len(list)
	for i := 0; i < length; i++ {
		l.Items[i] = *c.FromInvitationModel(list[i])
	}

	return l
}

// FromAuditModelList returns a list of
func (c Convertor) FromAuditModelList(list []*model.AuditEvent) *orgv1.AuditEventList {
	l := &orgv1.AuditEventList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "AuditEventList",
		},
		Items: make([]orgv1.AuditEvent, len(list)),
	}
	length := len(list)
	for i := 0; i < length; i++ {
		l.Items[i] = *c.FromAuditModel(list[i])
	}

	return l
}

// FromInvitationModel convert the model
func (c Convertor) FromInvitationModel(i *model.Invitation) *orgv1.TeamInvitation {
	return &orgv1.TeamInvitation{
		TypeMeta: metav1.TypeMeta{
			APIVersion: orgv1.SchemeGroupVersion.String(),
			Kind:       "TeamInvitationList",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              i.User.Username,
			Namespace:         HubNamespace,
			CreationTimestamp: metav1.NewTime(i.CreatedAt),
		},
		Spec: orgv1.TeamInvitationSpec{
			Username: i.User.Username,
			Team:     i.Team.Name,
		},
	}
}

// FromInvitationModelList converts from invitation model
func (c Convertor) FromInvitationModelList(list []*model.Invitation) *orgv1.TeamInvitationList {
	l := &orgv1.TeamInvitationList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "TeamInvitationList",
		},
	}
	for _, x := range list {
		l.Items = append(l.Items, *c.FromInvitationModel(x))
	}

	return l
}

// FromUserModel converts the model user to api user
func (c Convertor) FromUserModel(user *model.User) *orgv1.User {
	return &orgv1.User{
		TypeMeta: metav1.TypeMeta{
			APIVersion: orgv1.SchemeGroupVersion.String(),
			Kind:       "User",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              user.Username,
			Namespace:         HubNamespace,
			CreationTimestamp: metav1.NewTime(user.CreatedAt),
		},
		Spec: orgv1.UserSpec{
			Username: user.Username,
			Email:    user.Email,
		},
		Status: orgv1.UserStatus{
			Status:     corev1.SuccessStatus,
			Conditions: []corev1.Condition{},
		},
	}
}

// FromUsersModelList returns a list of users
func (c Convertor) FromUsersModelList(users []*model.User) *orgv1.UserList {
	list := &orgv1.UserList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "UserList",
		},
		Items: make([]orgv1.User, len(users)),
	}
	for i := 0; i < len(users); i++ {
		list.Items[i] = *c.FromUserModel(users[i])
	}

	return list
}

// ToTeamModel converts from api to model
func (c Convertor) ToTeamModel(team *orgv1.Team) *model.Team {
	return &model.Team{
		Name:        team.Name,
		Description: team.Spec.Description,
		Summary:     team.Spec.Summary,
	}
}

// FromTeamModel converts the model team to api team
func (c Convertor) FromTeamModel(team *model.Team) *orgv1.Team {
	return &orgv1.Team{
		TypeMeta: metav1.TypeMeta{
			APIVersion: orgv1.SchemeGroupVersion.String(),
			Kind:       "Team",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              team.Name,
			Namespace:         HubNamespace,
			CreationTimestamp: metav1.NewTime(team.CreatedAt),
		},
		Spec: orgv1.TeamSpec{
			Description: team.Description,
			Summary:     team.Summary,
		},
		Status: orgv1.TeamStatus{
			Status:     corev1.SuccessStatus,
			Conditions: []corev1.Condition{},
		},
	}
}

// FromTeamsModelList returns a list of teams
func (c Convertor) FromTeamsModelList(teams []*model.Team) *orgv1.TeamList {
	list := &orgv1.TeamList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "TeamList",
		},
		Items: make([]orgv1.Team, len(teams)),
	}
	for i := 0; i < len(teams); i++ {
		list.Items[i] = *c.FromTeamModel(teams[i])
	}

	return list
}

// FromMembersToUserList returns a user list from members
func (c Convertor) FromMembersToUserList(members []*model.Member) *orgv1.UserList {
	users := make([]*model.User, len(members))

	for i := 0; i < len(members); i++ {
		users[i] = members[i].User
	}

	return c.FromUsersModelList(users)
}

// FromMembersToTeamList returns a team list
func (c Convertor) FromMembersToTeamList(members []*model.Member) *orgv1.TeamList {
	teams := make([]*model.Team, len(members))

	for i := 0; i < len(members); i++ {
		teams[i] = members[i].Team
	}

	return c.FromTeamsModelList(teams)
}

// ToSecurityScanResult converts to the DB-layer scan result from the Security scan result
func (c Convertor) ToSecurityScanResult(result *securityv1.SecurityScanResult) model.SecurityScanResult {
	res := model.SecurityScanResult{
		ID: result.Spec.ID,
		SecurityResourceReference: model.SecurityResourceReference{
			ResourceGroup:     result.Spec.Resource.Group,
			ResourceVersion:   result.Spec.Resource.Version,
			ResourceKind:      result.Spec.Resource.Kind,
			ResourceNamespace: result.Spec.Resource.Namespace,
			ResourceName:      result.Spec.Resource.Name,
		},
		OwningTeam:    result.Spec.OwningTeam,
		CheckedAt:     result.Spec.CheckedAt.Time,
		ArchivedAt:    result.Spec.ArchivedAt.Time,
		OverallStatus: result.Spec.OverallStatus.String(),
		Results:       make([]model.SecurityRuleResult, len(result.Spec.Results)),
	}
	for i, rr := range result.Spec.Results {
		res.Results[i] = c.ToSecurityRuleResult(rr)
	}
	return res
}

// FromSecurityScanResult converts from the DB-layer scan result to the Security scan result
func (c Convertor) FromSecurityScanResult(result *model.SecurityScanResult) securityv1.SecurityScanResult {

	res := securityv1.SecurityScanResult{
		TypeMeta: metav1.TypeMeta{
			APIVersion: securityv1.SchemeGroupVersion.String(),
			Kind:       "SecurityScanResult",
		},
		Spec: securityv1.SecurityScanResultSpec{
			ID: result.ID,
			Resource: corev1.Ownership{
				Group:     result.ResourceGroup,
				Version:   result.ResourceVersion,
				Kind:      result.ResourceKind,
				Namespace: result.ResourceNamespace,
				Name:      result.ResourceName,
			},
			OwningTeam:    result.OwningTeam,
			CheckedAt:     metav1.NewTime(result.CheckedAt),
			ArchivedAt:    metav1.NewTime(result.ArchivedAt),
			OverallStatus: securityv1.RuleStatus(result.OverallStatus),
			Results:       make([]*securityv1.SecurityScanRuleResult, len(result.Results)),
		},
	}
	for i, rr := range result.Results {
		res.Spec.Results[i] = c.FromSecurityRuleResult(rr)
	}
	return res
}

// ToSecurityRuleResult converts to the DB-layer rule result from the Security rule result
func (c Convertor) ToSecurityRuleResult(result *securityv1.SecurityScanRuleResult) model.SecurityRuleResult {
	return model.SecurityRuleResult{
		RuleCode:  result.RuleCode,
		Status:    result.Status.String(),
		Message:   result.Message,
		CheckedAt: result.CheckedAt.Time,
	}
}

// FromSecurityRuleResult converts from the DB-layer rule result to the Security rule result
func (c Convertor) FromSecurityRuleResult(result model.SecurityRuleResult) *securityv1.SecurityScanRuleResult {
	return &securityv1.SecurityScanRuleResult{
		RuleCode:  result.RuleCode,
		Status:    securityv1.RuleStatus(result.Status),
		Message:   result.Message,
		CheckedAt: metav1.NewTime(result.CheckedAt),
	}
}

// FromSecurityRuleList converts from the security API rule slice to the security rule list
func (c Convertor) FromSecurityRuleList(rules []security.Rule) securityv1.SecurityRuleList {
	result := securityv1.SecurityRuleList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: securityv1.SchemeGroupVersion.String(),
			Kind:       "SecurityRuleList",
		},
		Items: make([]securityv1.SecurityRule, len(rules)),
	}

	for i, r := range rules {
		result.Items[i] = c.FromSecurityRule(r)
	}

	return result
}

// FromSecurityRule converts from the security API rule to a security rule
func (c Convertor) FromSecurityRule(rule security.Rule) securityv1.SecurityRule {
	return securityv1.SecurityRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: securityv1.SchemeGroupVersion.String(),
			Kind:       "SecurityRule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: rule.Code(),
		},
		Spec: securityv1.SecurityRuleSpec{
			Code:        rule.Code(),
			Name:        rule.Name(),
			Description: rule.Description(),
			AppliesTo:   security.RuleApplies(rule),
		},
	}
}

func (c Convertor) FromSecurityOverview(overview *model.SecurityOverview) securityv1.SecurityOverview {
	o := securityv1.SecurityOverview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: securityv1.SchemeGroupVersion.String(),
			Kind:       "SecurityOverview",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "overview",
		},
		Spec: securityv1.SecurityOverviewSpec{
			OpenIssueCounts: map[securityv1.RuleStatus]uint64{},
			Resources:       make([]securityv1.SecurityResourceOverview, len(overview.Resources)),
		},
	}

	for k, v := range overview.OpenIssueCounts {
		o.Spec.OpenIssueCounts[securityv1.RuleStatus(k)] = v
	}

	for i, r := range overview.Resources {
		o.Spec.Resources[i] = c.FromSecurityResourceOverview(&r)
	}

	return o
}

func (c Convertor) FromSecurityResourceOverview(resource *model.SecurityResourceOverview) securityv1.SecurityResourceOverview {
	r := securityv1.SecurityResourceOverview{
		Resource: corev1.Ownership{
			Group:     resource.ResourceGroup,
			Version:   resource.ResourceVersion,
			Kind:      resource.ResourceKind,
			Namespace: resource.ResourceNamespace,
			Name:      resource.ResourceName,
		},
		OverallStatus:   securityv1.RuleStatus(resource.OverallStatus),
		LastChecked:     metav1.NewTime(resource.LastChecked),
		OpenIssueCounts: map[securityv1.RuleStatus]uint64{},
	}
	for k, v := range resource.OpenIssueCounts {
		r.OpenIssueCounts[securityv1.RuleStatus(k)] = v
	}
	return r
}
