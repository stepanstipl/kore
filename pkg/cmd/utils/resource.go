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

package utils

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/appvia/kore/pkg/client"
	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonpath"

	log "github.com/sirupsen/logrus"
)

// resourceImpl implements the Resources interface
type resourceImpl struct {
	client client.Interface
}

func newResourceManager(client client.Interface) (Resources, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}

	return &resourceImpl{client: client}, nil
}

// Lookup is used to check if a resource is supported
func (r *resourceImpl) Lookup(name string) (*Resource, error) {
	singular := strings.ToLower(utils.Singularize(name))
	log.WithFields(log.Fields{
		"kind":     name,
		"singular": singular,
	}).Debug("resource lookup discovered the following")

	for _, x := range ResourceList {
		if singular == x.Name || name == x.ShortName {
			log.WithFields(log.Fields{
				"scope":    x.Scope,
				"singular": singular,
			}).Debug("found a matching resource type")

			return &x, nil
		}
	}
	log.Debug("no resource type found, using default")

	return &DefaultResource, nil
}

// Names returns all the names of the resource types
func (r *resourceImpl) Names() ([]string, error) {
	return ResourceNames, nil
}

// List return a full list of resources
func (r *resourceImpl) List() ([]Resource, error) {
	return ResourceList, nil
}

// LookResourceNamesWithFilter returns a list of resource names against a regexp
func (r *resourceImpl) LookResourceNamesWithFilter(kind, team, filter string) ([]string, error) {
	list, err := r.LookupResourceNames(kind, team)
	if err != nil {
		return nil, err
	}

	match, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}

	var filtered []string
	for _, x := range list {
		if match.MatchString(x) {
			filtered = append(filtered, x)
		}
	}

	return filtered, nil
}

// LookupResourceNames returns a list of resources of a specific kind
func (r *resourceImpl) LookupResourceNames(kind, team string) ([]string, error) {
	// @step: first we lookup the resource from the cache
	resource, err := r.Lookup(kind)
	if err != nil {
		return nil, err
	}

	// @step: we then construct a request for the list of that type
	req := r.client.Request().Resource(utils.Pluralize(resource.Name))
	if resource.IsTeamScoped() {
		req.Team(team)
	}
	if err := req.Get().Error(); err != nil {
		return nil, err
	}

	// @step: we read in the response and parse the items.[].metadata.name
	resp, err := ioutil.ReadAll(req.Body())
	if err != nil {
		return nil, err
	}

	var list []string

	for _, x := range jsonpath.GetMany(string(resp), "items.#.metadata.name")[0].Array() {
		list = append(list, x.String())
	}

	return list, nil
}
