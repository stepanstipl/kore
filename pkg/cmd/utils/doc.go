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

// Package utils provides a client factory to the cli utilities
package utils

/*
	factory := NewFactory(stdout, stdin, stderr)

	resp, err := factory.Client().Context(ctx).
		Team(team).
		Name(name).
		Resource("team").
		Payload(object).
		Result(object).
		Get().
		Do()

	// Raw query
	resp, err := factory.Raw.Context(ctx).
		Endpoint("/teams/{team}/{resource}/{name}").
		Parameters(
			QueryParam("allocation", true),
			PathParam("team", team),
			PathParam("resource", resource),
			PathParam("name", team),
		).
		SetResult(object)
		Get().
		Do()

	// Response handlers
	resp, err := factory.Client().Context(ctx).
		Team(team).
		Resource(resource).
		SubResource("member").
		Name(name).
		Result(object).
		Get().
		Do()

*/
