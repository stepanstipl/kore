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

package korectl

import (
	"fmt"
	"io"
	"strings"

	"github.com/juju/ansiterm/tabwriter"
)

type printer interface {
	Print(output io.Writer) error
}

type table struct {
	header  string
	columns []string
	lines   [][]string
}

func (t table) Print(output io.Writer) error {
	if _, err := fmt.Printf("%s\n", t.header); err != nil {
		return err
	}

	w := new(tabwriter.Writer)
	w.Init(output, 0, 10, 4, ' ', 0)

	if _, err := fmt.Fprintf(w, "%s\n", strings.Join(t.columns, "\t")); err != nil {
		return err
	}

	for _, line := range t.lines {
		if _, err := fmt.Fprintf(w, "%s\n", strings.Join(line, "\t")); err != nil {
			return err
		}
	}

	return w.Flush()
}
