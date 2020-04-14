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

package render

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/appvia/kore/pkg/utils"
	"github.com/appvia/kore/pkg/utils/jsonpath"
	"github.com/tidwall/gjson"

	"sigs.k8s.io/yaml"
)

// renderer implements the contract
type renderer struct {
	columns     []*column
	disableCase bool
	err         error
	failmissing bool
	foreach     string
	format      string
	resource    string
	showheaders bool
	unknown     string
	writer      io.Writer
}

type column struct {
	name, path string
	formatters []PrinterColumnFormatter
}

// Render creates and returns a renderer
func Render() Interface {
	return &renderer{
		format:      FormatTable,
		showheaders: true,
		writer:      os.Stdout,
		unknown:     "Unknown",
	}
}

// Do performs the actual rendering
func (r *renderer) Do() error {
	// @check if we had an error while formatting - technically you can loose
	// multiple errors on have this is implemented but it's fine for now
	if r.err != nil || r.resource == "" {
		return r.err
	}

	// @check the resource is valid json
	if !jsonpath.IsValid(r.resource) {
		return ErrInvalidResource
	}

	// @step: we handle the direct render first
	switch r.format {
	case FormatJSON:
		return r.RenderJSON()
	case FormatYAML:
		return r.RenderYAML()
	default:
		return r.RenderTable()
	}
}

// IsDocumentPath checks if it's a document path
func IsDocumentPath(path string) bool {
	return path == "."
}

// RenderTable is responsible for rendering the table
func (r *renderer) RenderTable() error {
	// @check if table format and columns empty we return nothing
	if len(r.columns) <= 0 {
		return nil
	}

	// @step: create the tabwriter
	wr := new(tabwriter.Writer)
	wr.Init(r.writer, 0, 10, 4, ' ', tabwriter.DiscardEmptyColumns)
	defer wr.Flush()

	// @step: create the headers and print if required
	if r.showheaders {
		headers := r.ColumnHeadersLine()
		_, _ = wr.Write([]byte(headers))
	}

	// @step: geneate the paths
	paths := r.ColumnPaths()

	columnFunc := func(value gjson.Result) error {
		results := make([]gjson.Result, len(r.columns))

		for i := 0; i < len(r.columns); i++ {
			if IsDocumentPath(paths[0]) {
				results[i] = gjson.Result{Type: gjson.String, Str: value.String()}

				continue
			}

			results[i] = value.Get(r.columns[i].path)
		}

		return r.RenderColumn(wr, results)
	}

	// @step: if we are not splitting by a key
	value := jsonpath.Parse(r.resource)
	if r.foreach != "" {
		value = value.Get(r.foreach)
		if !value.Exists() {
			return nil
		}
	}

	err := func() error {
		if value.IsArray() {
			for _, c := range value.Array() {
				if err := columnFunc(c); err != nil {
					return err
				}
			}
			return nil
		}

		return columnFunc(value)
	}()

	return err
}

// RenderColumn is used to render the resulting column from document
func (r *renderer) RenderColumn(wr io.Writer, list []gjson.Result) error {
	values := make([]string, len(r.columns))

	// @step: iterate the columns and find the values
	for i, x := range list {
		if !x.Exists() {
			if r.failmissing {
				return ErrMissingKey
			}
			values[i] = r.unknown
		} else {
			values[i] = x.String()
		}
		c := r.columns[i]

		//@step: apply the modifiers to the value
		for _, fn := range c.formatters {
			values[i] = fn(values[i])
		}
	}
	// @step: join the values together
	fmt.Fprintf(wr, "%s", strings.Join(values, "\t")+"\n")

	return nil
}

// ShowHeaders indicates if the headers should be shown
func (r *renderer) ShowHeaders(v bool) Interface {
	r.showheaders = v

	return r
}

// RenderJSON is responsible for rendering the json
func (r *renderer) RenderJSON() error {
	fmt.Fprintf(r.writer, r.resource)

	return nil
}

// RenderYAML is responsible for rendering the yaml
func (r *renderer) RenderYAML() error {
	value, err := yaml.JSONToYAML([]byte(r.resource))
	if err != nil {
		return err
	}
	fmt.Fprintf(r.writer, "%s", value)

	return nil
}

// Format defines the output format
func (r *renderer) Format(v string) Interface {
	if v == "" {
		v = DefaultRender
	}

	if !utils.Contains(v, SupportedFormats()) {
		r.err = ErrInvalidFormat
	}
	r.format = v

	return r
}

// Foreach indicates a iteration of items
func (r *renderer) Foreach(v string) Interface {
	r.foreach = v

	return r
}

// Resource defines the resource containing the json
func (r *renderer) Resource(fn ResourceInputFunc) Interface {
	resource, err := fn()
	if err != nil {
		r.err = err

		return r
	}
	r.resource = resource

	return r
}

// Writer defines the io.Writer to render content
func (r *renderer) Writer(v io.Writer) Interface {
	if v == nil {
		r.err = ErrNoWriter
	}
	r.writer = v

	return r
}

// Printer defines a table printer for the resource
func (r *renderer) Printer(list ...PrinterColumnFunc) Interface {
	for _, fn := range list {
		if fn == nil {
			r.err = &ErrInvalidColumn{message: "print column method is nil"}

			return r
		}
		c, err := fn()
		if err != nil {
			r.err = &ErrInvalidColumn{message: err.Error()}

			return r
		}

		r.columns = append(r.columns, c)
	}

	return r
}

// FailOnMissing controls the error on missing keys
func (r *renderer) FailOnMissing() Interface {
	r.failmissing = true

	return r
}

// Unknown defined the value of unknown values
func (r *renderer) Unknown(v string) Interface {
	r.unknown = v

	return r
}

// ColumnPaths returns a list of paths to retrieve
func (r *renderer) ColumnPaths() []string {
	list := make([]string, len(r.columns))

	for i := 0; i < len(r.columns); i++ {
		list[i] = r.columns[i].path
	}

	return list
}

// ColumnHeadersLine returns the tabwriter header
func (r *renderer) ColumnHeadersLine() string {
	return strings.Join(r.ColumnHeaders(), "\t") + "\n"
}

// ColumnHeaders returns a list of column names
func (r *renderer) ColumnHeaders() []string {
	list := make([]string, len(r.columns))

	for i := 0; i < len(r.columns); i++ {
		value := r.columns[i].name
		if !r.disableCase {
			value = strings.ToUpper(value)
		}
		list[i] = value
	}

	return list
}

// DisableUpperCaseHeaders indicates we should disable upper casing the headers
func (r *renderer) DisableUpperCaseHeaders() Interface {
	r.disableCase = true

	return r
}

// SupportedFormats returns a list of formats
func SupportedFormats() []string {
	return []string{FormatJSON, FormatYAML, FormatTable, FormatTemplate}
}
