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

// V1alpha1OrganizationSpec v1alpha1 organization spec
//
// swagger:model v1alpha1.OrganizationSpec
type V1alpha1OrganizationSpec struct {

	// billing account
	// Required: true
	BillingAccount *string `json:"billingAccount"`

	// credentials ref
	// Required: true
	CredentialsRef *V1SecretReference `json:"credentialsRef"`

	// parent ID
	// Required: true
	ParentID *string `json:"parentID"`

	// parent type
	// Required: true
	ParentType *string `json:"parentType"`

	// service account
	// Required: true
	ServiceAccount *string `json:"serviceAccount"`

	// token ref
	TokenRef *V1SecretReference `json:"tokenRef,omitempty"`
}

// Validate validates this v1alpha1 organization spec
func (m *V1alpha1OrganizationSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBillingAccount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCredentialsRef(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParentID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParentType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateServiceAccount(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTokenRef(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1OrganizationSpec) validateBillingAccount(formats strfmt.Registry) error {

	if err := validate.Required("billingAccount", "body", m.BillingAccount); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1OrganizationSpec) validateCredentialsRef(formats strfmt.Registry) error {

	if err := validate.Required("credentialsRef", "body", m.CredentialsRef); err != nil {
		return err
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

func (m *V1alpha1OrganizationSpec) validateParentID(formats strfmt.Registry) error {

	if err := validate.Required("parentID", "body", m.ParentID); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1OrganizationSpec) validateParentType(formats strfmt.Registry) error {

	if err := validate.Required("parentType", "body", m.ParentType); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1OrganizationSpec) validateServiceAccount(formats strfmt.Registry) error {

	if err := validate.Required("serviceAccount", "body", m.ServiceAccount); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1OrganizationSpec) validateTokenRef(formats strfmt.Registry) error {

	if swag.IsZero(m.TokenRef) { // not required
		return nil
	}

	if m.TokenRef != nil {
		if err := m.TokenRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("tokenRef")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1OrganizationSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1OrganizationSpec) UnmarshalBinary(b []byte) error {
	var res V1alpha1OrganizationSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
