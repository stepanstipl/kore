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
	configv1 "github.com/appvia/kore/pkg/apis/config/v1"
	corev1 "github.com/appvia/kore/pkg/apis/core/v1"
	monitoring "github.com/appvia/kore/pkg/apis/monitoring/v1beta1"
	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	securityv1 "github.com/appvia/kore/pkg/apis/security/v1"
	"github.com/appvia/kore/pkg/persistence/model"
	"github.com/appvia/kore/pkg/security"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	// DefaultConvertor is a default type
	DefaultConvertor Convertor
)

// Convertor is used to convert the models
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

// FromAlertModel converts the rule model
func (c Convertor) FromAlertModel(alert *model.Alert) *monitoring.Alert {
	o := &monitoring.Alert{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoring.SchemeGroupVersion.String(),
			Kind:       "Alert",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              alert.Rule.Name,
			Namespace:         alert.Rule.Team.Name,
			CreationTimestamp: metav1.NewTime(alert.CreatedAt),
			UID:               types.UID(alert.UID),
			Annotations: map[string]string{
				"fingerprint": alert.Fingerprint,
			},
		},
		Spec: monitoring.AlertSpec{
			Event:   alert.RawAlert,
			Summary: alert.Summary,
		},
		Status: monitoring.AlertStatus{
			Status: alert.Status,
			Detail: alert.StatusMessage,
		},
	}

	if len(alert.Labels) > 0 {
		o.Spec.Labels = make(map[string]string)
		for _, x := range alert.Labels {
			o.Spec.Labels[x.Name] = x.Value
		}
	}

	if alert.Rule != nil {
		o.Status.Rule = c.FromAlertRuleModel(alert.Rule)
	}
	if alert.ArchivedAt != nil {
		o.Status.ArchivedAt = metav1.NewTime(*alert.ArchivedAt)
	}
	if alert.Expiration != nil {
		o.Status.SilencedUntil = metav1.NewTime(*alert.Expiration)
	}

	return o
}

// FromAlertRuleModel converts the rule model
func (c Convertor) FromAlertRuleModel(rule *model.AlertRule) *monitoring.AlertRule {
	o := &monitoring.AlertRule{
		TypeMeta: metav1.TypeMeta{
			APIVersion: monitoring.SchemeGroupVersion.String(),
			Kind:       "Rule",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              rule.Name,
			Namespace:         rule.Team.Name,
			CreationTimestamp: metav1.NewTime(rule.CreatedAt),
			Labels:            make(map[string]string),
		},
		Spec: monitoring.AlertRuleSpec{
			Severity: rule.Severity,
			Source:   rule.Source,
			Summary:  rule.Summary,
			RawRule:  rule.RawRule,
			Resource: corev1.Ownership{
				Group:     rule.ResourceGroup,
				Version:   rule.ResourceVersion,
				Kind:      rule.ResourceKind,
				Namespace: rule.ResourceNamespace,
				Name:      rule.ResourceName,
			},
		},
	}
	for _, label := range rule.Labels {
		o.Labels[label.Name] = label.Value
	}

	return o
}

// ToAlert converts the alert model
func (c Convertor) ToAlert(m *monitoring.Alert) *model.Alert {
	o := &model.Alert{
		Fingerprint:   m.GetAnnotations()["fingerprint"],
		RawAlert:      m.Spec.Event,
		Status:        m.Status.Status,
		StatusMessage: m.Status.Detail,
		Summary:       m.Spec.Summary,
	}
	if !m.Status.SilencedUntil.IsZero() {
		o.Expiration = &m.Status.SilencedUntil.Time
	}
	for k, v := range m.Spec.Labels {
		o.Labels = append(o.Labels, model.AlertLabel{Name: k, Value: v})
	}

	return o
}

// ToAlertRule convert the api to the model
func (c Convertor) ToAlertRule(o *monitoring.AlertRule) *model.AlertRule {
	m := &model.AlertRule{
		Name:     o.Name,
		RawRule:  o.Spec.RawRule,
		Source:   o.Spec.Source,
		Severity: o.Spec.Severity,
		Summary:  o.Spec.Summary,
		ResourceReference: model.ResourceReference{
			ResourceGroup:     o.Spec.Resource.Group,
			ResourceVersion:   o.Spec.Resource.Version,
			ResourceKind:      o.Spec.Resource.Kind,
			ResourceNamespace: o.Spec.Resource.Namespace,
			ResourceName:      o.Spec.Resource.Name,
		},
	}

	var labels []model.RuleLabel
	for k, v := range o.GetLabels() {
		labels = append(labels, model.RuleLabel{Name: k, Value: v})
	}
	m.Labels = labels

	return m
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

// FromAlertsRuleModelList return a list of rules
func (c Convertor) FromAlertsRuleModelList(rules []*model.AlertRule) (*monitoring.AlertRuleList, error) {
	list := &monitoring.AlertRuleList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "RuleList",
		},
		Items: make([]monitoring.AlertRule, len(rules)),
	}

	for i, x := range rules {
		list.Items[i] = *c.FromAlertRuleModel(x)
	}

	return list, nil
}

// FromAlertsModelList return a list of alerts
func (c Convertor) FromAlertsModelList(alerts []*model.Alert) (*monitoring.AlertList, error) {
	list := &monitoring.AlertList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "AlertList",
		},
		Items: make([]monitoring.Alert, len(alerts)),
	}

	for i, x := range alerts {
		list.Items[i] = *c.FromAlertModel(x)
	}

	return list, nil
}

// ToTeamModel converts from api to model
func (c Convertor) ToTeamModel(team *orgv1.Team) *model.Team {
	return &model.Team{
		Identifier:  team.Labels[LabelTeamIdentifier],
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
			Labels: map[string]string{
				LabelTeamIdentifier: team.Identifier,
			},
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

// FromSecurityOverview converts the security overview model
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

// FromSecurityResourceOverview converts from the model to api
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

// ToConfigModel converts from api to model
func (c Convertor) ToConfigModel(config *configv1.Config) *model.Config {
	confs := []model.ConfigItems{}
	for k, v := range config.Spec.Values {
		conf := model.ConfigItems{
			Key:   k,
			Value: v,
		}
		confs = append(confs, conf)
	}

	res := &model.Config{
		Name:  config.Name,
		Items: confs,
	}

	return res
}

// FromConfigModelList converts the list of configs
func (c Convertor) FromConfigModelList(config []*model.Config) *configv1.ConfigList {
	list := &configv1.ConfigList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigList",
		},
		Items: make([]configv1.Config, len(config)),
	}
	for i := 0; i < len(config); i++ {
		list.Items[i] = *c.FromConfigModel(config[i])
	}

	return list
}

// FromConfigModel converts the config user to api config
func (c Convertor) FromConfigModel(config *model.Config) *configv1.Config {
	values := make(map[string]string)
	for _, ea := range config.Items {
		values[ea.Key] = ea.Value
	}
	return &configv1.Config{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Config",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Name,
		},
		Spec: configv1.ConfigSpec{
			Values: values,
		},
	}
}

// ToIdentityModel convert from the api to db model
func (c Convertor) ToIdentityModel(o *orgv1.Identity) *model.Identity {
	model := &model.Identity{}
	if o.Spec.User != nil {
		model.User = c.ToUserModel(o.Spec.User)
	}
	switch {
	case o.Spec.BasicAuth != nil:
		model.ProviderToken = o.Spec.BasicAuth.Password
	case o.Spec.IDPUser != nil:
		model.ProviderEmail = o.Spec.IDPUser.Email
		model.ProviderUID = o.Spec.IDPUser.UUID
	}

	return model
}

// FromIdentityModelList converts the db model list into a api list
func (c Convertor) FromIdentityModelList(models []*model.Identity) *orgv1.IdentityList {
	list := &orgv1.IdentityList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "IdentityList",
		},
		Items: make([]orgv1.Identity, len(models)),
	}
	for i := 0; i < len(models); i++ {
		list.Items[i] = *c.FromIdentityModel(models[i])
	}

	return list
}

// FromIdentityModel convert the db model to the api model
func (c Convertor) FromIdentityModel(model *model.Identity) *orgv1.Identity {
	i := &orgv1.Identity{
		Spec: orgv1.IdentitySpec{
			User: c.FromUserModel(model.User),
		},
	}

	switch model.Provider {
	case "basicauth":
		i.Spec.BasicAuth = &orgv1.BasicAuth{}
		i.Spec.AccountType = AccountLocal
	case "sso":
		i.Spec.IDPUser = &orgv1.IDPUser{
			Email: model.ProviderEmail,
			UUID:  model.ProviderUID,
		}
		i.Spec.AccountType = AccountSSO
	}

	return i
}
