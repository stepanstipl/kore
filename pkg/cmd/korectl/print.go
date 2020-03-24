/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore.
 *
 * kore is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore.  If not, see <http://www.gnu.org/licenses/>.
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
