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
	"errors"
	"fmt"
	"io"
)

const (
	// DefaultRender is the table output if nothing specified
	DefaultRender = FormatTable
	// FormatTable indicates a table format
	FormatTable = "table"
	// FormatJSON indicates json
	FormatJSON = "json"
	// FormatYAML indicates a yaml output
	FormatYAML = "yaml"
	// FormatTemplate defines a go template
	FormatTemplate = "template"
)

var (
	// ErrNoWriter indicates no io.Writer has been specified
	ErrNoWriter = errors.New("no writer specified")
	// ErrInvalidFormat indicates the format selected is not supported
	ErrInvalidFormat = errors.New("unsupported format")
	// ErrInvalidReader indicates the resource io.Reader is invalid
	ErrInvalidReader = errors.New("invalid io.Reader")
	// ErrInvalidResource indicates the resource json is invalid
	ErrInvalidResource = errors.New("invalid resource json")
	// ErrMissingKey indicates the path was missing
	ErrMissingKey = errors.New("missing path")
)

// ErrInvalidColumn indicates a invalid column method
type ErrInvalidColumn struct {
	message string
}

// Error returns the column error
func (e *ErrInvalidColumn) Error() string {
	return fmt.Sprintf("invalid column: %s", e.message)
}

// Interface is the contract to the renderer
type Interface interface {
	// DisableUpperCaseHeaders indicates we should disable upper casing the headers
	DisableUpperCaseHeaders() Interface
	// Do performs the actual rendering
	Do() error
	// FailOnMissing controls the error on missing keys
	FailOnMissing() Interface
	// Foreach indicates a iteration of items
	Foreach(string) Interface
	// Format defines the output format
	Format(string) Interface
	// Printer defines a table printer for the resource
	Printer(...PrinterColumnFunc) Interface
	// Resource defines the resource containing the json
	Resource(ResourceInputFunc) Interface
	// ShowHeaders indicates if the headers should be shown
	ShowHeaders(bool) Interface
	// Writer defines the io.Writer to render content
	Writer(io.Writer) Interface
	// Unknown defined the value of unknown values
	Unknown(string) Interface
}

// ResourceInputFunc defines the input resource i.e. bytes, string, io.Reader
type ResourceInputFunc func() (string, error)

// PrinterColumnFunc defines the printer column method
type PrinterColumnFunc func() (*column, error)

// PrinterColumnFormatter is an optional formatter for the value type
type PrinterColumnFormatter func(string) string
