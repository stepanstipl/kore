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
	"bytes"
	"encoding/json"

	"github.com/ghodss/yaml"
	"k8s.io/apimachinery/pkg/runtime"
)

// EncodeRuntimeObjectToYAML is used to encode the object to a yaml document
func EncodeRuntimeObjectToYAML(object runtime.Object) ([]byte, error) {
	b := &bytes.Buffer{}

	// @step: encode to json first of all
	if err := json.NewEncoder(b).Encode(object); err != nil {
		return []byte{}, err
	}

	return yaml.JSONToYAML(b.Bytes())
}
