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

package persistence_test

import (
	"context"
	"testing"

	. "github.com/appvia/kore/pkg/persistence"
	"github.com/appvia/kore/pkg/persistence/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeFakeAlertRule() *model.AlertRule {
	return &model.AlertRule{
		Name:     "test rule",
		Summary:  "test rule",
		Source:   "prometheus",
		Severity: "warning",
		ResourceReference: model.ResourceReference{
			ResourceGroup:     "clusters.kore.appvia.io",
			ResourceKind:      "Cluster",
			ResourceName:      "test",
			ResourceNamespace: "test",
			ResourceVersion:   "v1",
		},
	}
}

func makeAlertRulesList(t *testing.T, filters ...ListFunc) []*model.AlertRule {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.AlertRules().List(context.Background(), filters...)
	require.NoError(t, err)
	require.NotNil(t, v)

	return v
}

func TestAlertRules(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	assert.NotNil(t, store.AlertRules())
}

func TestAlertRulesList(t *testing.T) {
	list := makeAlertRulesList(t)
	require.Equal(t, 2, len(list))
	assert.NotEmpty(t, list[0].Alerts)

	list = makeAlertRulesList(t,
		Filter.WithName("PodDown"),
	)
	require.Equal(t, 1, len(list))
	require.NotEmpty(t, list[0].Alerts)
	require.Equal(t, 2, len(list[0].Alerts))
}

func TestAlertRulesListLatest(t *testing.T) {
	list := makeAlertRulesList(t,
		Filter.WithName("TargetDown"),
		Filter.WithTeam("alert_team"),
		Filter.WithAlertLatest(),
	)
	require.Equal(t, 1, len(list))
	assert.NotEmpty(t, list[0].Alerts)
	assert.Equal(t, 1, len(list[0].Alerts))

	list = makeAlertRulesList(t,
		Filter.WithName("PodDown"),
		Filter.WithTeam("alert_team"),
	)
	require.Equal(t, 1, len(list))
	assert.NotEmpty(t, list[0].Alerts)
	assert.Equal(t, 2, len(list[0].Alerts))
}

func TestAlertRulesListActive(t *testing.T) {
	list := makeAlertRulesList(t,
		Filter.WithTeam("alert_team"),
		Filter.WithAlertStatus([]string{"Active"}),
	)
	require.Equal(t, 1, len(list))
}

func TestAlertRulesListWithHistory(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	list, err := store.AlertRules().List(context.TODO(),
		Filter.WithName("PodDown"),
		Filter.WithTeam("alert_team"),
		Filter.WithAlertHistory(1),
	)
	require.NoError(t, err)

	require.Equal(t, 1, len(list))
	require.NotEmpty(t, list[0].Alerts)

	require.Equal(t, 1, len(list[0].Alerts))

	list, err = store.AlertRules().List(context.TODO(),
		Filter.WithName("PodDown"),
		Filter.WithTeam("alert_team"),
		Filter.WithAlertHistory(2),
	)
	require.NoError(t, err)

	require.Equal(t, 1, len(list))
	require.NotEmpty(t, list[0].Alerts)
	require.Equal(t, 2, len(list[0].Alerts))
}

func TestAlertRulesListByTeam(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v := makeAlertRulesList(t, List.WithTeam("none"))
	require.Empty(t, v)

	v = makeAlertRulesList(t, List.WithTeam("alert_team"))
	assert.NotEmpty(t, v)
}

func TestAlertRulesUpdate(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	list := makeAlertRulesList(t)
	require.Equal(t, 2, len(list))

	rule := makeFakeAlertRule()
	rule.TeamID = 6

	require.NoError(t, store.AlertRules().Update(context.TODO(), rule))
	list = makeAlertRulesList(t)
	require.Equal(t, 3, len(list))

	rule.ID = 0
	require.NoError(t, store.AlertRules().Update(context.TODO(), rule))
	list = makeAlertRulesList(t)
	require.Equal(t, 3, len(list))

	require.NoError(t, store.AlertRules().Delete(context.TODO(), rule))
	list = makeAlertRulesList(t)
	require.Equal(t, 2, len(list))
}

func TestAlertRulesNewRule(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	list := makeAlertRulesList(t)
	require.Equal(t, 2, len(list))

	rule := makeFakeAlertRule()
	rule.TeamID = 6
	rule.Name = "New Rule"
	rule.Summary = "New Rule"

	require.NoError(t, store.AlertRules().Update(context.TODO(), rule))
	list = makeAlertRulesList(t)
	require.Equal(t, 3, len(list))

	rule, err := store.AlertRules().Get(context.TODO(), List.WithID(rule.ID))
	require.NoError(t, err)
	require.Equal(t, 1, len(rule.Alerts))

	store.AlertRules().Delete(context.TODO(), rule)
	list = makeAlertRulesList(t)
	require.Equal(t, 2, len(list))
}
