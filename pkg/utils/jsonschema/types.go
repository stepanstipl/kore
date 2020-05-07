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

package jsonschema

import (
	"fmt"
	"strconv"
)

// Schema is the JSON schema object
type Schema struct {
	Properties map[string]Property `json:"properties"`
}

// Property is an object property
type Property struct {
	Default     interface{}   `json:"default"`
	Description string        `json:"description"`
	Title       string        `json:"title"`
	Type        string        `json:"type"`
	Enum        []interface{} `json:"enum"`
}

func (p Property) ParseDefault() (interface{}, error) {
	if p.Default == nil {
		return nil, nil
	}

	switch p.Type {
	case "string":
		if val, ok := p.Default.(string); ok {
			return val, nil
		}
		return fmt.Sprintf("%v", p.Default), nil
	case "boolean":
		if val, ok := p.Default.(bool); ok {
			return val, nil
		}
		val, err := strconv.ParseBool(fmt.Sprintf("%v", p.Default))
		return val, err
	case "integer":
		if val, ok := p.Default.(int64); ok {
			return val, nil
		}
		if val, ok := p.Default.(float64); ok {
			return int64(val), nil
		}
		val, err := strconv.ParseInt(fmt.Sprintf("%v", p.Default), 10, 64)
		return val, err
	case "number":
		if val, ok := p.Default.(int64); ok {
			return val, nil
		}
		if val, ok := p.Default.(float64); ok {
			return val, nil
		}
		val, err := strconv.ParseFloat(fmt.Sprintf("%v", p.Default), 64)
		return val, err
	default:
		return p.Default, nil
	}
}
