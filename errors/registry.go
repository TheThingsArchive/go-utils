package errors

import (
	"fmt"
	"sync"
)

// registry represents an error type registry
type registry struct {
	sync.RWMutex
	byCode map[Code]*ErrDescriptor
}

// Register registers a new error type
func (r *registry) Register(err *ErrDescriptor) {
	r.Lock()
	defer r.Unlock()

	if err.Code == 0 {
		panic(fmt.Errorf("No code defined in error descriptor (message: `%s`)", err.MessageFormat))
	}

	if r.byCode[err.Code] != nil {
		panic(fmt.Errorf("errors: Duplicate error code %v registered", err.Code))
	}

	err.registered = true
	r.byCode[err.Code] = err
}

// Get returns the descriptor if it exists or nil otherwise
func (r *registry) Get(code Code) *ErrDescriptor {
	r.RLock()
	defer r.RUnlock()
	return r.byCode[code]
}

// reg is a global registry to be shared by packages
var reg = &registry{
	byCode: make(map[Code]*ErrDescriptor),
}

// Register registers a new error descriptor
func Register(descriptor *ErrDescriptor) {
	reg.Register(descriptor)
}

// Get returns an error descriptor based on an error code
func Get(code Code) *ErrDescriptor {
	return reg.Get(code)
}

// Descriptor returns the error descriptor from any error
func Descriptor(err error) (desc *ErrDescriptor) {
	var code Code

	// let's hope it's an Error
	e, ok := err.(Error)
	if ok {
		code = e.Code()
		desc = Get(code)
	}

	// TODO: try to get from http or grpc errors

	// if the descriptor was found, return it
	if desc != nil {
		return desc
	}

	// return a new error descriptor with sane defaults
	return &ErrDescriptor{
		MessageFormat: err.Error(),
		Type:          Unknown,
		Code:          code,
	}
}

// GetCode infers the error code from the error
func GetCode(err error) Code {
	return Descriptor(err).Code
}

// GetMessageFormat infers the message format from the error
// or falls back to the error message
func GetMessageFormat(err error) string {
	return Descriptor(err).MessageFormat
}

// GetType infers the error type from the error
// or falls back to Unknown
func GetType(err error) Type {
	return Descriptor(err).Type
}

// GetAttributes returns the error attributes or falls back
// to empty attributes
func GetAttributes(err error) Attributes {
	e, ok := err.(Error)
	if ok {
		return e.Attributes()
	}

	return Attributes{}
}
