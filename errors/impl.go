// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

// impl implements Error
type impl struct {
	Imessage    string     `json:"error"`
	Icode       Code       `json:"error_code,omitempty"`
	Ityp        Type       `json:"error_type,omitempty"`
	Iattributes Attributes `json:"attributes,omitempty"`
}

// Error returns the formatted error message
func (i *impl) Error() string {
	return i.Imessage
}

// Code returns the error code
func (i *impl) Code() Code {
	return i.Icode
}

// Type returns the error type
func (i *impl) Type() Type {
	return i.Ityp
}

// Attributes returns the error attributes
func (i *impl) Attributes() Attributes {
	return i.Iattributes
}

// toImpl creates an equivalent impl for any Error
func toImpl(err Error) *impl {
	if i, ok := err.(*impl); ok {
		return i
	}

	return &impl{
		Imessage:    err.Error(),
		Icode:       err.Code(),
		Ityp:        err.Type(),
		Iattributes: err.Attributes(),
	}
}
