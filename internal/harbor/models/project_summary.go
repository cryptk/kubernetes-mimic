// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ProjectSummary project summary
//
// swagger:model ProjectSummary
type ProjectSummary struct {

	// The total number of charts under this project.
	ChartCount int64 `json:"chart_count,omitempty"`

	// The total number of developer members.
	DeveloperCount int64 `json:"developer_count,omitempty"`

	// The total number of guest members.
	GuestCount int64 `json:"guest_count,omitempty"`

	// The total number of limited guest members.
	LimitedGuestCount int64 `json:"limited_guest_count,omitempty"`

	// The total number of maintainer members.
	MaintainerCount int64 `json:"maintainer_count,omitempty"`

	// The total number of project admin members.
	ProjectAdminCount int64 `json:"project_admin_count,omitempty"`

	// quota
	Quota *ProjectSummaryQuota `json:"quota,omitempty"`

	// registry
	Registry *Registry `json:"registry,omitempty"`

	// The number of the repositories under this project.
	RepoCount int64 `json:"repo_count,omitempty"`
}

// Validate validates this project summary
func (m *ProjectSummary) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateQuota(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRegistry(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProjectSummary) validateQuota(formats strfmt.Registry) error {
	if swag.IsZero(m.Quota) { // not required
		return nil
	}

	if m.Quota != nil {
		if err := m.Quota.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("quota")
			}
			return err
		}
	}

	return nil
}

func (m *ProjectSummary) validateRegistry(formats strfmt.Registry) error {
	if swag.IsZero(m.Registry) { // not required
		return nil
	}

	if m.Registry != nil {
		if err := m.Registry.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("registry")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this project summary based on the context it is used
func (m *ProjectSummary) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateQuota(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRegistry(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProjectSummary) contextValidateQuota(ctx context.Context, formats strfmt.Registry) error {

	if m.Quota != nil {
		if err := m.Quota.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("quota")
			}
			return err
		}
	}

	return nil
}

func (m *ProjectSummary) contextValidateRegistry(ctx context.Context, formats strfmt.Registry) error {

	if m.Registry != nil {
		if err := m.Registry.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("registry")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProjectSummary) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProjectSummary) UnmarshalBinary(b []byte) error {
	var res ProjectSummary
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ProjectSummaryQuota project summary quota
//
// swagger:model ProjectSummaryQuota
type ProjectSummaryQuota struct {

	// The hard limits of the quota
	Hard ResourceList `json:"hard,omitempty"`

	// The used status of the quota
	Used ResourceList `json:"used,omitempty"`
}

// Validate validates this project summary quota
func (m *ProjectSummaryQuota) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHard(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUsed(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProjectSummaryQuota) validateHard(formats strfmt.Registry) error {
	if swag.IsZero(m.Hard) { // not required
		return nil
	}

	if m.Hard != nil {
		if err := m.Hard.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("quota" + "." + "hard")
			}
			return err
		}
	}

	return nil
}

func (m *ProjectSummaryQuota) validateUsed(formats strfmt.Registry) error {
	if swag.IsZero(m.Used) { // not required
		return nil
	}

	if m.Used != nil {
		if err := m.Used.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("quota" + "." + "used")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this project summary quota based on the context it is used
func (m *ProjectSummaryQuota) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateHard(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUsed(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProjectSummaryQuota) contextValidateHard(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Hard.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("quota" + "." + "hard")
		}
		return err
	}

	return nil
}

func (m *ProjectSummaryQuota) contextValidateUsed(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Used.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("quota" + "." + "used")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProjectSummaryQuota) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProjectSummaryQuota) UnmarshalBinary(b []byte) error {
	var res ProjectSummaryQuota
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
