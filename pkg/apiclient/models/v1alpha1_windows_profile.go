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

// V1alpha1WindowsProfile v1alpha1 windows profile
//
// swagger:model v1alpha1.WindowsProfile
type V1alpha1WindowsProfile struct {

	// admin password
	// Required: true
	AdminPassword *string `json:"adminPassword"`

	// admin username
	// Required: true
	AdminUsername *string `json:"adminUsername"`
}

// Validate validates this v1alpha1 windows profile
func (m *V1alpha1WindowsProfile) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAdminPassword(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAdminUsername(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1WindowsProfile) validateAdminPassword(formats strfmt.Registry) error {

	if err := validate.Required("adminPassword", "body", m.AdminPassword); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1WindowsProfile) validateAdminUsername(formats strfmt.Registry) error {

	if err := validate.Required("adminUsername", "body", m.AdminUsername); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1WindowsProfile) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1WindowsProfile) UnmarshalBinary(b []byte) error {
	var res V1alpha1WindowsProfile
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}