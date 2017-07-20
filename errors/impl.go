// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

// impl implements Error
type impl struct {
	message    string     `json:"error"`
	code       Code       `json:"error_code,omitempty"`
	typ        Type       `json:"error_type,omitempty"`
	attributes Attributes `json:"attributes,omitempty"`
}

// Error returns the formatted error message
func (i *impl) Error() string {
	return i.message
}

// Code returns the error code
func (i *impl) Code() Code {
	return i.code
}

// Type returns the error type
func (i *impl) Type() Type {
	return i.typ
}

// Attributes returns the error attributes
func (i *impl) Attributes() Attributes {
	return i.attributes
}

// toImpl creates an equivalent impl for any Error
func toImpl(err Error) *impl {
	if i, ok := err.(*impl); ok {
		return i
	}

	return &impl{
		message:    err.Error(),
		code:       err.Code(),
		typ:        err.Type(),
		attributes: err.Attributes(),
	}
}
