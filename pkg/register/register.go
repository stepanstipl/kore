/**
 * Copyright (C) 2020 Appvia Ltd <info@appvia.io>
 *
 * This file is part of kore-apiserver.
 *
 * kore-apiserver is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 2 of the License, or
 * (at your option) any later version.
 *
 * kore-apiserver is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with kore-apiserver.  If not, see <http://www.gnu.org/licenses/>.
 */

package register

import (
	"bytes"
	"encoding/json"

	yaml "github.com/ghodss/yaml"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// GetCustomResourceDefinitions returns all the CRDs for the kore
func GetCustomResourceDefinitions() ([]*apiextensions.CustomResourceDefinition, error) {
	var list []*apiextensions.CustomResourceDefinition

	for _, asset := range AssetNames() {
		definition, err := Asset(asset)
		if err != nil {
			return nil, err
		}

		decoded, err := yaml.YAMLToJSON(definition)
		if err != nil {
			return list, err
		}

		crd := &apiextensions.CustomResourceDefinition{}
		if err := json.NewDecoder(bytes.NewReader(decoded)).Decode(crd); err != nil {
			return list, err
		}

		list = append(list, crd)
	}

	return list, nil
}
