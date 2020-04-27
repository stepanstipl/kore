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

package validation

import (
	"fmt"
	"strings"
)

// FieldRoot is used to reference the root object
const FieldRoot = "(root)"

// Error is a specific error returned when the input provided by
// the user has failed validation somehow.
type Error struct {
	// Code is an optional machine readable code used to describe the code
	Code int `json:"code"`
	// Message is a human readable message related to the error
	Message string `json:"message"`
	// FieldErrors are the individual validation errors found against the submitted data.
	FieldErrors []FieldError `json:"fieldErrors"`
}

// NewError returns a new validation error
func NewError(format string, args ...interface{}) *Error {
	return &Error{
		Code:    400,
		Message: fmt.Sprintf(strings.TrimRight(format, ":\n"), args...),
	}
}

// Error returns the details of the validation error.
func (e Error) Error() string {
	if len(e.FieldErrors) == 0 {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteString(e.Message)
	sb.WriteString(":\n")
	for _, fe := range e.FieldErrors {
		if fe.Field == FieldRoot {
			sb.WriteString(fmt.Sprintf(" * %s\n", fe.Message))
		} else {
			sb.WriteString(fmt.Sprintf(" * %s: %s\n", fe.Field, fe.Message))
		}
	}
	return sb.String()
}

// HasErrors returns true if any field errors have been added to this validation error.
func (e *Error) HasErrors() bool {
	return len(e.FieldErrors) > 0
}

// WithFieldError adds a field error to the validation error and returns it for fluent loveliness.
func (e *Error) WithFieldError(field string, errCode ErrorCode, message string) *Error {
	e.AddFieldError(field, errCode, message)
	return e
}

// WithFieldErrorf adds an error for a specific field to a validation error.
func (e *Error) WithFieldErrorf(field string, errCode ErrorCode, format string, args ...interface{}) *Error {
	e.AddFieldErrorf(field, errCode, format, args...)
	return e
}

// AddFieldError adds an error for a specific field to a validation error.
func (e *Error) AddFieldError(field string, errCode ErrorCode, message string) {
	e.FieldErrors = append(e.FieldErrors, FieldError{
		Field:   field,
		ErrCode: errCode,
		Message: message,
	})
}

// AddFieldErrorf adds an error for a specific field to a validation error.
func (e *Error) AddFieldErrorf(field string, errCode ErrorCode, format string, args ...interface{}) {
	e.FieldErrors = append(e.FieldErrors, FieldError{
		Field:   field,
		ErrCode: errCode,
		Message: fmt.Sprintf(format, args...),
	})
}

// FieldError provides information about a validation error on a specific field.
type FieldError struct {
	// Field causing the error, in format x.y.z
	Field string `json:"field"`
	// ErrCode is the type of constraint which has been broken.
	ErrCode ErrorCode `json:"errCode"`
	// Message is a human-readable description of the validation error.
	Message string `json:"message"`
}

// ErrorCode is the type of validation error detected.
type ErrorCode string

// The error codes should match the validator names from JSON Schema
const (
	// MinLength error indicates the supplied value is shorted than the allowed minimum.
	MinLength ErrorCode = "minLength"
	// MaxLength error indicates the supplied value is longer than the allowed maximum.
	MaxLength ErrorCode = "maxLength"
	// Required error indicates that a field must be specified.
	Required ErrorCode = "required"
	// Pattern error indicates the input doesn't match the required regex pattern
	Pattern ErrorCode = "pattern"
	// MustExist error indicates that the named reference must exist
	MustExist ErrorCode = "mustExist"
	// ReadOnly error indicates that the given value can not be changed from a pre-defined value
	ReadOnly ErrorCode = "readOnly"
	// InvalidType error indicates that we've expected a different type
	InvalidType ErrorCode = "invalidType"
	// InvalidValue error indicates that the given value is invalid
	InvalidValue ErrorCode = "invalidValue"
)
