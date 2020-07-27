// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// V1alpha1AKSSpec v1alpha1 a k s spec
//
// swagger:model v1alpha1.AKSSpec
type V1alpha1AKSSpec struct {

	// agent pool profiles
	// Required: true
	AgentPoolProfiles []*V1alpha1AgentPoolProfile `json:"agentPoolProfiles"`

	// authorized IP ranges
	AuthorizedIPRanges []string `json:"authorizedIPRanges"`

	// cluster
	Cluster *V1Ownership `json:"cluster,omitempty"`

	// credentials
	// Required: true
	Credentials *V1Ownership `json:"credentials"`

	// description
	// Required: true
	Description *string `json:"description"`

	// dns prefix
	// Required: true
	DNSPrefix *string `json:"dnsPrefix"`

	// enable pod security policy
	EnablePodSecurityPolicy bool `json:"enablePodSecurityPolicy,omitempty"`

	// enable private cluster
	EnablePrivateCluster bool `json:"enablePrivateCluster,omitempty"`

	// kubernetes version
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`

	// linux profile
	LinuxProfile *V1alpha1LinuxProfile `json:"linuxProfile,omitempty"`

	// location
	// Required: true
	Location *string `json:"location"`

	// network plugin
	// Required: true
	NetworkPlugin *string `json:"networkPlugin"`

	// network policy
	NetworkPolicy string `json:"networkPolicy,omitempty"`

	// tags
	Tags map[string]string `json:"tags,omitempty"`

	// windows profile
	WindowsProfile *V1alpha1WindowsProfile `json:"windowsProfile,omitempty"`
}

// Validate validates this v1alpha1 a k s spec
func (m *V1alpha1AKSSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAgentPoolProfiles(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCluster(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCredentials(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDescription(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDNSPrefix(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLinuxProfile(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLocation(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateNetworkPlugin(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWindowsProfile(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1AKSSpec) validateAgentPoolProfiles(formats strfmt.Registry) error {

	if err := validate.Required("agentPoolProfiles", "body", m.AgentPoolProfiles); err != nil {
		return err
	}

	for i := 0; i < len(m.AgentPoolProfiles); i++ {
		if swag.IsZero(m.AgentPoolProfiles[i]) { // not required
			continue
		}

		if m.AgentPoolProfiles[i] != nil {
			if err := m.AgentPoolProfiles[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("agentPoolProfiles" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1alpha1AKSSpec) validateCluster(formats strfmt.Registry) error {

	if swag.IsZero(m.Cluster) { // not required
		return nil
	}

	if m.Cluster != nil {
		if err := m.Cluster.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cluster")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1AKSSpec) validateCredentials(formats strfmt.Registry) error {

	if err := validate.Required("credentials", "body", m.Credentials); err != nil {
		return err
	}

	if m.Credentials != nil {
		if err := m.Credentials.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("credentials")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1AKSSpec) validateDescription(formats strfmt.Registry) error {

	if err := validate.Required("description", "body", m.Description); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1AKSSpec) validateDNSPrefix(formats strfmt.Registry) error {

	if err := validate.Required("dnsPrefix", "body", m.DNSPrefix); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1AKSSpec) validateLinuxProfile(formats strfmt.Registry) error {

	if swag.IsZero(m.LinuxProfile) { // not required
		return nil
	}

	if m.LinuxProfile != nil {
		if err := m.LinuxProfile.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("linuxProfile")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1AKSSpec) validateLocation(formats strfmt.Registry) error {

	if err := validate.Required("location", "body", m.Location); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1AKSSpec) validateNetworkPlugin(formats strfmt.Registry) error {

	if err := validate.Required("networkPlugin", "body", m.NetworkPlugin); err != nil {
		return err
	}

	return nil
}

func (m *V1alpha1AKSSpec) validateWindowsProfile(formats strfmt.Registry) error {

	if swag.IsZero(m.WindowsProfile) { // not required
		return nil
	}

	if m.WindowsProfile != nil {
		if err := m.WindowsProfile.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("windowsProfile")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1AKSSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1AKSSpec) UnmarshalBinary(b []byte) error {
	var res V1alpha1AKSSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
