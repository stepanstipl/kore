/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package kore

import (
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/services/users/model"

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

// FromAuditModel convets the model
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
			CreatedAt:   metav1.NewTime(i.CreatedAt),
			Type:        i.Type,
			Team:        i.Team,
			User:        i.User,
			Message:     i.Message,
			Resource:    i.Resource,
			ResourceUID: i.ResourceUID,
		},
	}
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
