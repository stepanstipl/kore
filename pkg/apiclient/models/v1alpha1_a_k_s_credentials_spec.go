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

// V1alpha1AKSCredentialsSpec v1alpha1 a k s credentials spec
//
// swagger:model v1alpha1.AKSCredentialsSpec
type V1alpha1AKSCredentialsSpec struct {

	// client ID
	// Required: true
	ClientID *string `json:"clientID"`

	// credentials ref
	CredentialsRef *V1SecretReference `json:"credentialsRef,omitempty"`

	// subscription ID
	// Required: true
	SubscriptionID *string `json:"subscriptionID"`

	// tenant ID
	// Required: true
	TenantID *string `json:"tenantID"`
}

// Validate validates this v1alpha1 a k s credentials spec
func (m *V1alpha1AKSCredentialsSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateClientID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCredentialsRef(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSubscriptionID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTenantID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1AKSCredentialsSpec) validateClientID(formats strfmt.Registry) error {

	if err := validate.Required("clientID", "body", m.ClientID); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1AKSCredentialsSpec) validateCredentialsRef(formats strfmt.Registry) error {

	if swag.IsZero(m.CredentialsRef) { // not required
		return nil
	}

	if m.CredentialsRef != nil {
		if err := m.CredentialsRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("credentialsRef")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1AKSCredentialsSpec) validateSubscriptionID(formats strfmt.Registry) error {

	if err := validate.Required("subscriptionID", "body", m.SubscriptionID); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1AKSCredentialsSpec) validateTenantID(formats strfmt.Registry) error {

	if err := validate.Required("tenantID", "body", m.TenantID); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1AKSCredentialsSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1AKSCredentialsSpec) UnmarshalBinary(b []byte) error {
	var res V1alpha1AKSCredentialsSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}