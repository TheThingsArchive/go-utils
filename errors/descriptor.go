// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import "fmt"

// ErrDescriptor is a helper struct to easily build new Errors from and to be
// the authoritive information about error codes.
//
// The descriptor can be used to find out information about the error after it
// has been handed over between components
type ErrDescriptor struct {
	// MessageFormat is the format of the error message. Attributes will be filled
	// in when an error is created using New(). For example:
	//
	//   "This is an error about user {username}"
	//
	// when passed an atrtributes map with "username" set to "john" would interpolate to
	//
	//   "This is an error about user john"
	//
	// The idea about this message format is that is is localizable
	MessageFormat string

	// Code is the code of errors that are created by this descriptor
	Code Code

	// Type is the type of errors created by this descriptor
	Type Type

	// registered denotes wether or not the error has been registered
	// (by a call to Register)
	registered bool
}

// Format formats the attributes into an Error
func (err *ErrDescriptor) New(attributes Attributes) Error {
	if err.Code != 0 && !err.registered {
		panic(fmt.Errorf("Error descriptor with code %v was not registered", err.Code))
	}

	return &impl{
		Imessage:    Format(err.MessageFormat, attributes),
		Icode:       err.Code,
		Ityp:        err.Type,
		Iattributes: attributes,
	}
}

// New creates a new Error from a descriptor and some attributes
func New(descriptor *ErrDescriptor, attributes Attributes) Error {
	return descriptor.New(attributes)
}

// Register registers the descriptor
func (err *ErrDescriptor) Register() {
	Register(err)
}
