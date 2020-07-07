// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1beta1AlertStatus v1beta1 alert status
//
// swagger:model v1beta1.AlertStatus
type V1beta1AlertStatus struct {

	// archived at
	ArchivedAt string `json:"archivedAt,omitempty"`

	// detail
	Detail string `json:"detail,omitempty"`

	// rule
	Rule *V1beta1AlertRule `json:"rule,omitempty"`

	// silenced until
	SilencedUntil string `json:"silencedUntil,omitempty"`

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this v1beta1 alert status
func (m *V1beta1AlertStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateRule(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1AlertStatus) validateRule(formats strfmt.Registry) error {

	if swag.IsZero(m.Rule) { // not required
		return nil
	}

	if m.Rule != nil {
		if err := m.Rule.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("rule")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1beta1AlertStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1beta1AlertStatus) UnmarshalBinary(b []byte) error {
	var res V1beta1AlertStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
