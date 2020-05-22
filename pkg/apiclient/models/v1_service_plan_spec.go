// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// V1ServicePlanSpec v1 service plan spec
//
// swagger:model v1.ServicePlanSpec
type V1ServicePlanSpec struct {

	// configuration
	Configuration interface{} `json:"configuration,omitempty"`

	// credential schema
	CredentialSchema string `json:"credentialSchema,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// display name
	DisplayName string `json:"displayName,omitempty"`

	// kind
	// Required: true
	Kind *string `json:"kind"`

	// labels
	Labels map[string]string `json:"labels,omitempty"`

	// provider data
	ProviderData string `json:"providerData,omitempty"`

	// schema
	Schema string `json:"schema,omitempty"`

	// summary
	// Required: true
	Summary *string `json:"summary"`
}

// Validate validates this v1 service plan spec
func (m *V1ServicePlanSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateKind(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSummary(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1ServicePlanSpec) validateKind(formats strfmt.Registry) error {

	if err := validate.Required("kind", "body", m.Kind); err != nil {
		return err
	}

	return nil
}

func (m *V1ServicePlanSpec) validateSummary(formats strfmt.Registry) error {

	if err := validate.Required("summary", "body", m.Summary); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1ServicePlanSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1ServicePlanSpec) UnmarshalBinary(b []byte) error {
	var res V1ServicePlanSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
