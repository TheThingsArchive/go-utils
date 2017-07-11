// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

type jsonError struct {
	Message    string     `json:"error"`
	Code       Code       `json:"error_code,omitempty"`
	Type       Type       `json:"error_type,omitempty"`
	Attributes Attributes `json:"attributes,omitempty"`
}

func toJson(err Error) *jsonError {
	return &jsonError{
		Message:    err.Error(),
		Code:       err.Code(),
		Type:       err.Type(),
		Attributes: err.Attributes(),
	}
}
