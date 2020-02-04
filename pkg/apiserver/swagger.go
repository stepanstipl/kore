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

package apiserver

import "github.com/go-openapi/spec"

// EnrichSwagger provides the swagger config
func EnrichSwagger(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Appvia Kore API",
			Description: "Kore API provides the frontend API for the Appvia Kore (kore.appvia.io)",
			Contact: &spec.ContactInfo{
				Name:  "Rohith Jayawardene",
				Email: "info@appvia.io",
				URL:   "https://appvia.io",
			},
			License: &spec.License{
				Name: "GPLV2",
				URL:  "http://mit.org",
			},
			Version: "0.0.1",
		},
	}
}
