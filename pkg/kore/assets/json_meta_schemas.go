/**
 * Source: https://github.com/jgonzalezdr/json-schema-spec/blob/draft-07-strict/strict-schema.json
 * Upstream source: http://json-schema.org/draft-07/strict-schema
 *
 * License note from project's README:
 *  > The source material in this repository is licensed under the AFL or BSD license.
 *
 * No license text is included in the repository.
 */

package assets

// JSONSchemaDraft07 is the meta schema for the JSON Schema Draft 07 definition
// We use a strict schema definition instead of the default one.
const JSONSchemaDraft07 = `
{
    "$schema": "http://json-schema.org/draft-07/strict-schema#",
    "$id": "http://json-schema.org/draft-07/strict-schema#",
    "title": "Core schema meta-schema",
    "definitions": {
        "schemaArray": {
            "type": "array",
            "minItems": 1,
            "items": { "$ref": "#" }
        },
        "nonNegativeInteger": {
            "type": "integer",
            "minimum": 0
        },
        "nonNegativeIntegerDefault0": {
            "allOf": [
                { "$ref": "#/definitions/nonNegativeInteger" },
                { "default": 0 }
            ]
        },
        "simpleTypes": {
            "enum": [
                "array",
                "boolean",
                "integer",
                "null",
                "number",
                "object",
                "string"
            ]
        },
        "stringArray": {
            "type": "array",
            "items": { "type": "string" },
            "uniqueItems": true,
            "default": []
        },
        "requiresArrayType": {
            "type": "object",
            "properties": {
                "type": {
                    "anyOf": [ {
                        "const": "array"
                    }, {
                        "type": "array",
                        "minItems": 1,
                        "contains": {
                            "const": "array"
                        }
                    } ]
                }
            }
        },
        "requiresObjectType": {
            "type": "object",
            "properties": {
                "type": {
                    "anyOf": [ {
                        "const": "object"
                    }, {
                        "type": "array",
                        "minItems": 1,
                        "contains": {
                            "const": "object"
                        }
                    } ]
                }
            }
        },
        "requiresStringType": {
            "type": "object",
            "properties": {
                "type": {
                    "anyOf": [ {
                        "const": "string"
                    }, {
                        "type": "array",
                        "minItems": 1,
                        "contains": {
                            "const": "string"
                        }
                    } ]
                }
            }
        },
        "requiresNumberType": {
            "type": "object",
            "properties": {
                "type": {
                    "anyOf": [ {
                        "enum": ["number", "integer"]
                    }, {
                        "type": "array",
                        "minItems": 1,
                        "contains": {
                            "enum": ["number", "integer"]
                        }
                    } ]
                }
            }
        }
    },
    "type": ["object", "boolean"],
    "properties": {
        "$id": {
            "type": "string",
            "format": "uri-reference"
        },
        "$schema": {
            "type": "string",
            "format": "uri"
        },
        "$ref": {
            "type": "string",
            "format": "uri-reference"
        },
        "$comment": {
            "type": "string"
        },
        "title": {
            "type": "string"
        },
        "description": {
            "type": "string"
        },
        "default": true,
        "readOnly": {
            "type": "boolean",
            "default": false
        },
        "examples": {
            "type": "array",
            "items": true
        },
        "multipleOf": {
            "type": "number",
            "exclusiveMinimum": 0
        },
        "maximum": {
            "type": "number"
        },
        "exclusiveMaximum": {
            "type": "number"
        },
        "minimum": {
            "type": "number"
        },
        "exclusiveMinimum": {
            "type": "number"
        },
        "maxLength": { "$ref": "#/definitions/nonNegativeInteger" },
        "minLength": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
        "pattern": {
            "type": "string",
            "format": "regex"
        },
        "additionalItems": { "$ref": "#" },
        "items": {
            "anyOf": [
                { "$ref": "#" },
                { "$ref": "#/definitions/schemaArray" }
            ],
            "default": true
        },
        "maxItems": { "$ref": "#/definitions/nonNegativeInteger" },
        "minItems": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
        "uniqueItems": {
            "type": "boolean",
            "default": false
        },
        "contains": { "$ref": "#" },
        "maxProperties": { "$ref": "#/definitions/nonNegativeInteger" },
        "minProperties": { "$ref": "#/definitions/nonNegativeIntegerDefault0" },
        "required": { "$ref": "#/definitions/stringArray" },
        "additionalProperties": { "$ref": "#" },
        "definitions": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "properties": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "default": {}
        },
        "patternProperties": {
            "type": "object",
            "additionalProperties": { "$ref": "#" },
            "propertyNames": { "format": "regex" },
            "default": {}
        },
        "dependencies": {
            "type": "object",
            "additionalProperties": {
                "anyOf": [
                    { "$ref": "#" },
                    { "$ref": "#/definitions/stringArray" }
                ]
            }
        },
        "propertyNames": { "$ref": "#" },
        "const": true,
        "enum": {
            "type": "array",
            "items": true,
            "minItems": 1,
            "uniqueItems": true
        },
        "type": {
            "anyOf": [
                { "$ref": "#/definitions/simpleTypes" },
                {
                    "type": "array",
                    "items": { "$ref": "#/definitions/simpleTypes" },
                    "minItems": 1,
                    "uniqueItems": true
                }
            ]
        },
        "format": { "type": "string" },
        "contentMediaType": { "type": "string" },
        "contentEncoding": { "type": "string" },
        "if": { "$ref": "#" },
        "then": { "$ref": "#" },
        "else": { "$ref": "#" },
        "allOf": { "$ref": "#/definitions/schemaArray" },
        "anyOf": { "$ref": "#/definitions/schemaArray" },
        "oneOf": { "$ref": "#/definitions/schemaArray" },
        "not": { "$ref": "#" }
    },
    "required": ["type"],
    "additionalProperties": false,
    "dependencies": {
        "multipleOf": { "$ref": "#/definitions/requiresNumberType" },
        "maximum": { "$ref": "#/definitions/requiresNumberType" },
        "minimum": { "$ref": "#/definitions/requiresNumberType" },
        "exclusiveMaximum": { "$ref": "#/definitions/requiresNumberType" },
        "exclusiveMinimum": { "$ref": "#/definitions/requiresNumberType" },
        "maxLength": { "$ref": "#/definitions/requiresStringType" },
        "minLength": { "$ref": "#/definitions/requiresStringType" },
        "pattern": { "$ref": "#/definitions/requiresStringType" },
        "additionalItems": { "$ref": "#/definitions/requiresArrayType" },
        "items": { "$ref": "#/definitions/requiresArrayType" },
        "maxItems": { "$ref": "#/definitions/requiresArrayType" },
        "minItems": { "$ref": "#/definitions/requiresArrayType" },
        "uniqueItems": { "$ref": "#/definitions/requiresArrayType" },
        "contains": { "$ref": "#/definitions/requiresArrayType" },
        "maxProperties": { "$ref": "#/definitions/requiresObjectType" },
        "minProperties": { "$ref": "#/definitions/requiresObjectType" },
        "required": { "$ref": "#/definitions/requiresObjectType" },
        "additionalProperties": { "$ref": "#/definitions/requiresObjectType" },
        "properties": { "$ref": "#/definitions/requiresObjectType" },
        "patternProperties": { "$ref": "#/definitions/requiresObjectType" },
        "dependencies": { "$ref": "#/definitions/requiresObjectType" },
        "propertiesNames": { "$ref": "#/definitions/requiresObjectType" }
    },
    "default": true
}
`