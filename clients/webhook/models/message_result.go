// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MessageResult message result
//
// swagger:model MessageResult
type MessageResult struct {

	// message
	// Example: Accepted
	Message string `json:"message,omitempty"`

	// message Id
	// Example: 67f2f8a8-ea58-4ed0-a6f9-ff217df4d849
	MessageID string `json:"messageId,omitempty"`
}

// Validate validates this message result
func (m *MessageResult) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this message result based on context it is used
func (m *MessageResult) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MessageResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MessageResult) UnmarshalBinary(b []byte) error {
	var res MessageResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
