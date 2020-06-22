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
	"github.com/jinzhu/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeFakeAlert() *model.Alert {
	return &model.Alert{
		Status: model.AlertStatusOK,
		RuleID: 0,
	}
}

func makeAlertList(t *testing.T, filters ...ListFunc) []*model.Alert {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Alerts().List(context.Background(), filters...)
	require.NoError(t, err)
	require.NotNil(t, v)

	return v
}

func TestAlerts(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	assert.NotNil(t, store.Alerts())
}

func TestAlertsList(t *testing.T) {
	v := makeAlertList(t)
	require.NotEmpty(t, v)
	require.Equal(t, 3, len(v))

	v = makeAlertList(t, Filter.WithAlertStatus([]string{"Active"}))
	require.NotEmpty(t, v)
	require.Equal(t, 1, len(v))

	v = makeAlertList(t, Filter.WithStatus("Active"))
	require.NotEmpty(t, v)
	require.Equal(t, 1, len(v))
}

func TestAlertsListByLabels(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v, err := store.Alerts().List(context.Background(),
		Filter.WithAlertLabels([]string{"job=kubelet"}),
	)
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v))

	v, err = store.Alerts().List(context.Background(),
		Filter.WithAlertLabels([]string{"job=none"}),
	)
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 0, len(v))

	v, err = store.Alerts().List(context.Background(),
		Filter.WithAlertLabels([]string{"job=kubelet"}),
		Filter.WithStatus("Active"),
	)
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v))
}

func TestAlertsListByTeam(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v := makeAlertList(t, List.WithTeam("none"))
	require.Empty(t, v)

	v = makeAlertList(t, List.WithTeam("alert_team"))
	assert.NotEmpty(t, v)
}

func TestAlertsListByStatus(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v := makeAlertList(t, List.WithTeam("alert_team"))
	require.NotEmpty(t, v)
	assert.Equal(t, 3, len(v))

	v = makeAlertList(t,
		List.WithTeam("alert_team"),
		List.WithStatus(model.AlertStatusActive),
	)
	assert.NotEmpty(t, v)
	assert.Equal(t, 1, len(v))
}

func TestAlertsListByIdentity(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	v := makeAlertList(t,
		Filter.WithResourceGroup("clusters.kore.appvia.io"),
		Filter.WithResourceVersion("v1"),
		Filter.WithResourceKind("Cluster"),
		Filter.WithNamespace("test"),
		Filter.WithResourceName("test"),
	)
	require.NotEmpty(t, v)
	require.Equal(t, 3, len(v))

	v = makeAlertList(t,
		Filter.WithResourceGroup("clusters.kore.appvia.io"),
		Filter.WithResourceVersion("v1"),
		Filter.WithResourceKind("Cluster"),
		Filter.WithNamespace("test"),
		Filter.WithResourceName("test"),
		Filter.WithStatus("Active"),
	)
	require.NotEmpty(t, v)
	require.Equal(t, 1, len(v))
}

func TestAlertUpdateNoRule(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	require.Error(t, store.Alerts().Update(context.TODO(), makeFakeAlert()))
}

func TestAlertsUpdate(t *testing.T) {
	store := makeTestStore(t)
	defer store.Stop()

	ctx := context.Background()
	name := "NewRule"

	v, err := store.AlertRules().Get(ctx, Filter.WithName(name))
	require.Equal(t, err, gorm.ErrRecordNotFound)
	require.Empty(t, v)

	// 1. create a new rule
	rule := makeFakeAlertRule()
	rule.Name = name
	rule.TeamID = 6
	err = store.AlertRules().Update(ctx, rule)
	require.NoError(t, err)

	// 2. get the rule
	rule, err = store.AlertRules().Get(ctx,
		Filter.WithName(name),
		Filter.WithAlertLatest(),
	)
	require.NoError(t, err)
	require.Equal(t, 1, len(rule.Alerts))

	// 3. add an alert on the alert - we should have OK and Active
	alert := makeFakeAlert()
	alert.Rule = rule
	alert.Status = model.AlertStatusActive
	alert.Summary = "one"
	alert.Labels = []model.AlertLabel{
		{Name: "job", Value: "kubelet"},
		{Name: "namespace", Value: "kore-system"},
	}
	require.NoError(t, store.Alerts().Update(ctx, alert))

	a, _ := store.Alerts().Get(ctx, Filter.WithName(name), Filter.WithAlertLatest())
	require.Equal(t, "Active", a.Status)

	l, _ := store.Alerts().List(ctx, Filter.WithName(name), Filter.WithStatus("Active"))
	require.NotEmpty(t, l)
	require.Equal(t, 1, len(l))

	// 4. add a new instance on the same rule
	alert = makeFakeAlert()
	alert.Rule = rule
	alert.Status = model.AlertStatusActive
	alert.Summary = "two"
	alert.Labels = []model.AlertLabel{
		{Name: "job", Value: "kubelet"},
		{Name: "namespace", Value: "kore-system"},
		{Name: "pod", Value: "new"},
	}
	require.NoError(t, store.Alerts().Update(ctx, alert))

	// we should have two active
	l, _ = store.Alerts().List(ctx,
		Filter.WithName(name),
		Filter.WithStatus("Active"),
		Filter.WithAlertLatest(),
	)
	require.NotEmpty(t, l)
	require.Equal(t, 1, len(l))

	// 5. toggle off one of the alerts
	alert = makeFakeAlert()
	alert.Rule = rule
	alert.Status = model.AlertStatusOK
	alert.Summary = "three"
	alert.Labels = []model.AlertLabel{
		{Name: "job", Value: "kubelet"},
		{Name: "namespace", Value: "kore-system"},
		{Name: "pod", Value: "new"},
	}
	require.NoError(t, store.Alerts().Update(ctx, alert))

	/*
		l, _ = store.Alerts().List(ctx,
			Filter.WithName(name),
			Filter.WithStatus("Active"),
			Filter.WithAlertLatest(),
		)
		require.NotEmpty(t, l)
		require.Equal(t, 1, len(l))
	*/
}
