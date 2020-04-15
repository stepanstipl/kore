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

// Package client provides a client factory
package client

/*
	client := New(config)

	resp, err := client.Request().
		Context(ctx).
		Team(team).
		Name(name).
		Resource("team").
		Payload(object).
		Result(object).
		Get()

	// Raw query
	resp, err := client.Request().Context(ctx).
		Endpoint("/teams/{team}/{resource}/{name}").
		Parameters(
			QueryParameter("allocation", true),
			PathParameter("team", team),
			PathParameter("resource", resource),
			PathParameter("name", team),
		).
		Result(object)
		Get()

	// Response handlers
	resp, err := client.Request().
		Context(ctx).
		Team(team).
		Resource(resource).
		SubResource("member").
		Name(name).
		Result(object).
		Get()

*/
