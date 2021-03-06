// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1ServiceCredentialsStatus v1 service credentials status
//
// swagger:model v1.ServiceCredentialsStatus
type V1ServiceCredentialsStatus struct {

	// components
	Components []*V1Component `json:"components"`

	// message
	Message string `json:"message,omitempty"`

	// provider data
	ProviderData interface{} `json:"providerData,omitempty"`

	// provider ID
	ProviderID string `json:"providerID,omitempty"`

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this v1 service credentials status
func (m *V1ServiceCredentialsStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateComponents(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1ServiceCredentialsStatus) validateComponents(formats strfmt.Registry) error {

	if swag.IsZero(m.Components) { // not required
		return nil
	}

	for i := 0; i < len(m.Components); i++ {
		if swag.IsZero(m.Components[i]) { // not required
			continue
		}

		if m.Components[i] != nil {
			if err := m.Components[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("components" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1ServiceCredentialsStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1ServiceCredentialsStatus) UnmarshalBinary(b []byte) error {
	var res V1ServiceCredentialsStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
