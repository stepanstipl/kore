/*
 * Copyright (C) 2019  Rohith Jayawardene <gambol99@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"bytes"
	"encoding/json"
)

// ConvertToMap converts a struct to a map - note the fields must be
// exported for refection to work
func ConvertToMap(v interface{}) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	if v == nil {
		return values, nil
	}

	encoded := &bytes.Buffer{}
	if err := json.NewEncoder(encoded).Encode(v); err != nil {
		return nil, err
	}

	return values, json.NewDecoder(encoded).Decode(&values)
}
