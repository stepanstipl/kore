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

package utils

import (
	"fmt"
	"strings"
)

// ToPlural converts the type to a plural
func ToPlural(name string) string {
	if strings.HasSuffix(name, "ss") {
		return fmt.Sprintf("%ses", name)
	}
	if strings.HasSuffix(name, "ys") {
		return fmt.Sprintf("%sies", strings.TrimSuffix(name, "ys"))
	}
	if strings.HasSuffix(name, "es") || strings.HasSuffix(name, "s") {
		return name
	}

	return fmt.Sprintf("%ss", name)
}
