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

// Package render is used to render the result to the screen
package render

/*
	// render some pretty json
	render.Render().
		Format("json").
		Resource(json).
		Do()

	// Render a table to stdout
	render.Render().
		Writer(io.Stdout). // default is stdout
		Resource(FromBytes()).
		Format("table").   // default is table
		Foreach("items").  // used to iterate
		Printer(           // without a print the table format should error
			Column("Name", "metadata.name"),
			// the third method is an optional formatter
			Column("Age", "metadata.created", render.Age()),
		).Do()


	// Render by gotemplate
	render.Render().
		Writer(io.Stdout).
		Resource(FromString(`{ "hello": "world" }`)).
		Format("template").
		Template("Hello {{ .hello }}").
		Do()
*/
