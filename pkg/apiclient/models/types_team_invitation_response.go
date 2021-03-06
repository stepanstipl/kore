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

// TypesTeamInvitationResponse types team invitation response
//
// swagger:model types.TeamInvitationResponse
type TypesTeamInvitationResponse struct {

	// team
	// Required: true
	Team *string `json:"team"`
}

// Validate validates this types team invitation response
func (m *TypesTeamInvitationResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTeam(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *TypesTeamInvitationResponse) validateTeam(formats strfmt.Registry) error {

	if err := validate.Required("team", "body", m.Team); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *TypesTeamInvitationResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *TypesTeamInvitationResponse) UnmarshalBinary(b []byte) error {
	var res TypesTeamInvitationResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
