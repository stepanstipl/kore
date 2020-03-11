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

package kore

import "fmt"

// ErrValidation is a specific error returned when the input provided by
// the user has failed validation somehow.
type ErrValidation struct {
	FieldErrors []ValidationFieldError `json:"fieldErrors"`
}

// newErrValidation returns a new api validation error
func newErrValidation() *ErrValidation {
	return &ErrValidation{}
}

// Error returns the details of the validation error.
func (e ErrValidation) Error() string {
	msg := ""
	for ind, fe := range e.FieldErrors {
		if ind > 0 {
			msg += "\n"
		}
		msg += fmt.Sprintf("Validation error - field: %s error: %s message: %s", fe.Field, fe.ErrCode, fe.Message)
	}
	return msg
}

// WithFieldError adds an error for a specific field to a validation error.
func (e *ErrValidation) WithFieldError(field string, errCode ValidationErrorCode, message string) *ErrValidation {
	e.FieldErrors = append(e.FieldErrors, ValidationFieldError{Field: field, ErrCode: errCode, Message: message})
	return e
}

// ValidationFieldError provides information about a validation error on a specific field.
type ValidationFieldError struct {
	// Field causing the error, in format x.y.z
	Field string `json:"field"`
	// ErrCode is the type of constraint which has been broken.
	ErrCode ValidationErrorCode `json:"errCode"`
	// Message is a human-readable description of the validation error.
	Message string `json:"message"`
}

// ValidationErrorCode is the type of validation error detected.
type ValidationErrorCode int

const (
	// MaxLength error indicates the supplied value is longer than the allowed maximum.
	MaxLength ValidationErrorCode = iota
)

func (v ValidationErrorCode) String() string {
	names := [...]string{
		"maxLength",
	}
	if v < MaxLength || v > MaxLength {
		return "Unknown"
	}
	return names[v]
}
