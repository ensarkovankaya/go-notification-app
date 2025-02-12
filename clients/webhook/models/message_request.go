// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// MessageRequest message request
//
// swagger:model MessageRequest
type MessageRequest struct {

	// content
	// Example: Hello, World!
	Content string `json:"content,omitempty"`

	// to
	// Example: +9055511111111
	To string `json:"to,omitempty"`
}

// Validate validates this message request
func (m *MessageRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this message request based on context it is used
func (m *MessageRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *MessageRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *MessageRequest) UnmarshalBinary(b []byte) error {
	var res MessageRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
