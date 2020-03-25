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

package assets_test

import (
	"encoding/json"
	"testing"

	"github.com/appvia/kore/pkg/kore/assets"
	"github.com/appvia/kore/pkg/utils/jsonschema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGKESchemaCompiles(t *testing.T) {
	require.NotPanics(t, func() {
		_ = jsonschema.Validate(assets.GKEPlanSchema, "test", nil)
	})
}

func TestGKESchemaIsValidDraft07Schema(t *testing.T) {
	var schema map[string]interface{}
	err := json.Unmarshal([]byte(assets.GKEPlanSchema), &schema)
	require.NoError(t, err)

	err = jsonschema.Validate(assets.JSONSchemaDraft07, "test", schema)
	assert.NoError(t, err, "GKEPlanSchema is not a valid")
}
