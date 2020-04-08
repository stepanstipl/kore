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

package apiserver

import (
	"fmt"

	"github.com/go-openapi/spec"
)

// EnrichSwagger provides the swagger config
func EnrichSwagger(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Appvia Kore API",
			Description: "Kore API provides the frontend API for the Appvia Kore (kore.appvia.io)",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "Appvia Ltd",
					Email: "info@appvia.io",
					URL:   "https://appvia.io",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "Apache 2.0",
					URL:  "http://www.apache.org/licenses/LICENSE-2.0",
				},
			},
			Version: "0.0.1",
		},
	}
	swo.SecurityDefinitions = spec.SecurityDefinitions{
		"OAuth2": &spec.SecurityScheme{
			// @TODO: Set these correctly for the currently-running system:
			SecuritySchemeProps: spec.SecuritySchemeProps{
				Type:             "oauth2",
				Flow:             "accessCode",
				AuthorizationURL: "http://localhost:10080/auth",
				TokenURL:         "http://localhost:10080/token",
				Scopes: map[string]string{
					"admin": "Admin scope",
					"team":  "Team scope",
				},
			},
		},
	}
	swo.Security = []map[string][]string{
		{
			"OAuth2": {"admin", "team"},
		},
	}

	// These are horrible hacks to override the type for apiextv1.JSON properties, which is handled as "string",
	// but it should be an "object". ModelTypeNameHandler didn't work in restfulspec.Config.
	enrichSwaggerFixRawJSON(swo, "v1.PlanSpec", "configuration")
	enrichSwaggerFixRawJSON(swo, "v1.ClusterSpec", "configuration")

}

func enrichSwaggerFixRawJSON(swo *spec.Swagger, typeName, propertyName string) {
	def, ok := swo.Definitions[typeName]
	if !ok {
		panic(fmt.Errorf("%q doesn't exist, you may have to amend apiserver.EnrichSwagger", typeName))
	}

	property, ok := def.Properties[propertyName]
	if !ok {
		panic(fmt.Errorf("%q property doesn't exist in %q, you may have to amend apiserver.EnrichSwagger", propertyName, typeName))
	}

	property.Type = []string{"object"}
	def.Properties[propertyName] = property
	swo.Definitions[typeName] = def
}
