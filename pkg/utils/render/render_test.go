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
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testone = `apiVersion: test/v1
kind: Test
metadata:
  labels:
    hello: world
  name: test
  namespace: test-namespace
spec:
  hello: world
  one: 2
  something:
    more: complex
    with:
    - 1
    - 2
    - 3
status:
  empty: ""
  status: success
`
	testList = `apiVersion: v1
kind: List
items:
- apiVersion: test/v1
  kind: Test
  metadata:
    name: test
    namespace: test-namespace
  spec:
    hello: world
  status:
    status: success
- apiVersion: test/v1
  kind: Test
  metadata:
    name: test1
    namespace: test-namespace
  spec:
    hello: world
  status:
    status: success
`
	testUserList = `apiVersion: v1
kind: List
items:
- one
- two
- three
`
)

func makeTest(t *testing.T, document string) string {
	d, err := yaml.YAMLToJSON([]byte(document))
	require.NoError(t, err)

	return string(d)
}

func TestRenderJSON(t *testing.T) {
	b := &bytes.Buffer{}
	err := Render().
		Writer(b).
		Format(FormatJSON).
		Resource(FromString(makeTest(t, testone))).
		Do()

	require.NoError(t, err)
	assert.Equal(t, makeTest(t, testone), b.String())
}

func TestRenderJSONBad(t *testing.T) {
	err := Render().
		Format(FormatJSON).
		Resource(FromString(`{`)).
		Do()

	require.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
}

func TestRenderYAML(t *testing.T) {
	b := &bytes.Buffer{}
	err := Render().
		Writer(b).
		Format(FormatYAML).
		Resource(FromString(makeTest(t, testone))).
		Do()

	require.NoError(t, err)
	assert.Equal(t, testone, b.String())
}
func TestRenderYAMLBad(t *testing.T) {
	err := Render().
		Format(FormatYAML).
		Resource(FromString(`{`)).
		Do()

	require.Error(t, err)
	assert.Equal(t, ErrInvalidResource, err)
}

func TestRenderFromStruct(t *testing.T) {
	thing := struct {
		Name        string            `json:"name"`
		NoTag       string            `json:"no_tag"`
		Age         int               `json:"age"`
		Things      []string          `json:"things"`
		MapOfThings map[string]string `json:"map-of-things"`
	}{
		Name:   "test",
		NoTag:  "notag",
		Age:    20,
		Things: []string{"one", "two"},
		MapOfThings: map[string]string{
			"hello": "world",
		},
	}
	b := &bytes.Buffer{}

	err := Render().
		Writer(b).
		Format(FormatTable).
		Resource(FromStruct(&thing)).
		DisableUpperCaseHeaders().
		Printer(
			Column("Name", "name"),
			Column("NoTag", "NoTag"),
		).Do()
	require.NoError(t, err)
	assert.Equal(t, "Name    NoTag\ntest    Unknown\n", b.String())
}

func TestRenderFromStructArray(t *testing.T) {
	things := []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{Name: "test", Age: 20},
		{Name: "test1", Age: 21},
		{Name: "test2", Age: 22},
	}
	b := &bytes.Buffer{}

	err := Render().
		Writer(b).
		Format(FormatTable).
		Resource(FromStruct(&things)).
		DisableUpperCaseHeaders().
		Printer(
			Column("Name", "name"),
			Column("Age", "age"),
		).Do()
	require.NoError(t, err)
	assert.Equal(t, "Name     Age\ntest     20\ntest1    21\ntest2    22\n", b.String())
}

func TestShowWithHeaders(t *testing.T) {
	b := &bytes.Buffer{}
	err := Render().
		Writer(b).
		Resource(FromString(makeTest(t, testone))).
		DisableUpperCaseHeaders().
		Printer(
			Column("Name", "metadata.name"),
		).Do()

	require.NoError(t, err)
	require.Equal(t, "Name\ntest\n", b.String())
}

func TestShowWithNoHeaders(t *testing.T) {
	b := &bytes.Buffer{}
	err := Render().
		Writer(b).
		ShowHeaders(false).
		Resource(FromString(makeTest(t, testone))).
		DisableUpperCaseHeaders().
		Printer(
			Column("Name", "metadata.name"),
		).Do()

	require.NoError(t, err)
	require.Equal(t, "test\n", b.String())
}

func TestFailOnMissing(t *testing.T) {
	err := Render().
		Writer(ioutil.Discard).
		FailOnMissing().
		Resource(FromString(makeTest(t, testone))).
		DisableUpperCaseHeaders().
		Printer(
			Column("Name", "not_there"),
		).Do()

	require.Error(t, err)
	assert.Equal(t, ErrMissingKey, err)
}
func TestManyFailOnMissing(t *testing.T) {
	err := Render().
		Writer(ioutil.Discard).
		FailOnMissing().
		DisableUpperCaseHeaders().
		Resource(FromString(makeTest(t, testone))).
		Printer(
			Column("Name", "metadata.name"),
			Column("NotThere", "not_there"),
		).Do()

	require.Error(t, err)
	assert.Equal(t, ErrMissingKey, err)
}

func TestTableDoOk(t *testing.T) {
	cases := []struct {
		Columns  []PrinterColumnFunc
		Data     string
		Expected string
	}{
		{},
		{Data: makeTest(t, testone)},
		{
			Data: makeTest(t, testone),
			Columns: []PrinterColumnFunc{
				Column("Name", "metadata.name"),
			},
			Expected: "Name\ntest\n",
		},
		{
			Data: makeTest(t, testone),
			Columns: []PrinterColumnFunc{
				Column("Name", "metadata.name"),
				Column("Namespace", "metadata.namespace"),
			},
			Expected: "Name\tNamespace\ntest\ttest-namespace\n",
		},
		{
			Data: makeTest(t, testone),
			Columns: []PrinterColumnFunc{
				Column("One", "spec.one"),
			},
			Expected: "One\n2\n",
		},
		{
			Data: makeTest(t, testone),
			Columns: []PrinterColumnFunc{
				Column("Name", "spec.something.with|@sjoin"),
			},
			Expected: "Name\n1,2,3\n",
		},
		{
			Data: makeTest(t, testone),
			Columns: []PrinterColumnFunc{
				Column("Name", "not_there"),
			},
			Expected: "Name\nUnknown\n",
		},
	}
	for i, c := range cases {
		b := &bytes.Buffer{}
		err := Render().
			Writer(b).
			Printer(c.Columns...).
			DisableUpperCaseHeaders().
			Resource(FromString(c.Data)).
			Do()

		formatted := strings.ReplaceAll(c.Expected, "\t", "    ")

		require.NoError(t, err, "case %d, did not expect error", i)
		assert.Equal(t, formatted, b.String(), "case %d not as expected", i)
	}
}

func TestTableDoListOk(t *testing.T) {
	cases := []struct {
		Columns  []PrinterColumnFunc
		Data     string
		Expected string
	}{
		{},
		{
			Data: makeTest(t, testList),
			Columns: []PrinterColumnFunc{
				Column("Name", "metadata.name"),
			},
			Expected: "Name\ntest\ntest1\n",
		},
		{
			Data: makeTest(t, testList),
			Columns: []PrinterColumnFunc{
				Column("Name", "metadata.name"),
				Column("Namespace", "metadata.name"),
				Column("Not_there", "not_there"),
			},
			Expected: "Name\tNamespace\tNot_there\ntest\ttest\tUnknown\ntest1\ttest1\tUnknown\n",
		},
		{
			Data: makeTest(t, testUserList),
			Columns: []PrinterColumnFunc{
				Column("Name", "."),
			},
			Expected: "Name\none\ntwo\nthree\n",
		},
	}
	for i, c := range cases {
		b := &bytes.Buffer{}
		err := Render().
			Writer(b).
			DisableUpperCaseHeaders().
			Foreach("items").
			Printer(c.Columns...).
			Resource(FromString(c.Data)).
			Do()

		clean := strings.ReplaceAll(c.Expected, "\t", "")
		whitespace := strings.ReplaceAll(b.String(), " ", "")

		require.NoError(t, err, "case %d, did not expect error", i)
		assert.Equal(t, clean, whitespace, "case %d not as expected", i)
	}
}

func TestRenderAge(t *testing.T) {
	data := `{ "created": "2020-04-03T20:27:37Z" }`
	b := &bytes.Buffer{}
	err := Render().
		Writer(b).
		DisableUpperCaseHeaders().
		Resource(FromString(data)).
		Printer(
			Column("Age", "created", Age()),
		).Do()
	require.NoError(t, err)
	assert.NotEqual(t, "Invalid", b.String())
}

func TestRender(t *testing.T) {
	require.NotNil(t, Render())
}

func TestRenderNothing(t *testing.T) {
	require.NoError(t, Render().Do())
}

func TestInvalidWriter(t *testing.T) {
	err := Render().Writer(nil).Do()
	require.Error(t, err)
	require.Equal(t, err, ErrNoWriter)
}

func TestInvalidFormat(t *testing.T) {
	err := Render().Format("bad").Do()
	require.Error(t, err)
	require.Equal(t, err, ErrInvalidFormat)
}

func TestResourceFromString(t *testing.T) {
	data := `{ "hello": "world" }`
	r := Render().Resource(FromString(data))
	v := r.(*renderer)

	assert.NoError(t, v.err)
	assert.Equal(t, data, v.resource)
}

func TestPrinterColumnAdded(t *testing.T) {
	r := Render().Printer(
		Column("test", "test"),
	)
	require.NotNil(t, r)
	assert.NoError(t, r.(*renderer).err)
	assert.Equal(t, 1, len(r.(*renderer).columns))
}

func TestPrinterColumnNil(t *testing.T) {
	err := Render().Printer(
		Column("test", "test", nil),
	).Do()
	require.Error(t, err)
	assert.Equal(t, &ErrInvalidColumn{message: "invalid column: formatter method is nil"}, err)
}
func TestUnknown(t *testing.T) {
	r := Render().Unknown("test")
	require.NotNil(t, r)
	assert.NoError(t, r.(*renderer).err)
	assert.Equal(t, "test", r.(*renderer).unknown)
}

func TestPrinterColumnBad(t *testing.T) {
	r := Render().Printer(
		Column("", "test"),
	)
	require.NotNil(t, r)
	assert.Error(t, r.(*renderer).err)
	assert.Equal(t, &ErrInvalidColumn{message: "invalid column: no name"}, r.(*renderer).err)
}

func TestRenderWithNoResource(t *testing.T) {
	assert.NoError(t, Render().Printer(Column("test", "test")).Do())
}
