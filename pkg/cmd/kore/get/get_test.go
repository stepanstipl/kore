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

package get

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
	"time"

	orgv1 "github.com/appvia/kore/pkg/apis/org/v1"
	"github.com/appvia/kore/pkg/client/config"
	"github.com/appvia/kore/pkg/client/fake"
	"github.com/appvia/kore/pkg/cmd/errors"
	cmdutil "github.com/appvia/kore/pkg/cmd/utils"
	"github.com/appvia/kore/pkg/kore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeFakeTeam() *orgv1.Team {
	return &orgv1.Team{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Team",
			APIVersion: orgv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test",
			Namespace:         kore.HubNamespace,
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: orgv1.TeamSpec{
			Summary:     "test",
			Description: "test",
		},
	}
}

func makeFakeFactory(t *testing.T, w io.Writer) (fake.Interface, cmdutil.Factory) {
	if w == nil {
		w = ioutil.Discard
	}

	streams := cmdutil.Streams{Stdout: w, Stderr: w}
	client := fake.NewFake(&config.Config{})
	factory, err := cmdutil.NewFactory(client, streams, &config.Config{})
	require.NoError(t, err)
	require.NotNil(t, factory)

	return client, factory
}

func TestGetNotResource(t *testing.T) {
	_, f := makeFakeFactory(t, nil)
	o := &GetOptions{Factory: f}
	assert.Equal(t, errors.ErrMissingResource, cmdutil.ExecuteHandler(o))
}

func TestGet(t *testing.T) {
	b := &bytes.Buffer{}
	c, f := makeFakeFactory(t, b)
	c.SetResult(makeFakeTeam())

	o := &GetOptions{
		Factory:  f,
		Headers:  true,
		Name:     "test",
		Resource: "team",
	}
	assert.NoError(t, cmdutil.ExecuteHandler(o))
	//assert.Equal(t, "", b.String())
}

/*
func TestGetResourceNotFound() {
	svc := cmdtest.NewFakeAPI()
	svc.TrapGet("/teams/test", func(req *http.Request, resp http.ResponseWriter) {



	})
	client, cfg := cmdtest.NewFakeClientFromServer(svc.Endpoint())
	factory := utils.NewFactory(client, streams)

	o := &GetOptions{
		Factory:  f,
		Headers:  true,
		Name:     "not_there",
		Resource: "team",
	}

	assert.Error(t, cmdutil.ExecuteHandler(o))
}
*/
