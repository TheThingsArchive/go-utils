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
func Register(descriptors ...*ErrDescriptor) {
	for _, descriptor := range descriptors {
		reg.Register(descriptor)
	}
}

// Get returns an error descriptor based on an error code
func Get(code Code) *ErrDescriptor {
	return reg.Get(code)
}

// From lifts an error to be and Error
func From(in error) Error {
	if err, ok := in.(Error); ok {
		return err
	}

	return FromGRPC(in)
}

// Descriptor returns the error descriptor from any error
func Descriptor(in error) (desc *ErrDescriptor) {
	err := From(in)
	descriptor := Get(err.Code())
	if descriptor != nil {
		return descriptor
	}

	// return a new error descriptor with sane defaults
	return &ErrDescriptor{
		MessageFormat: err.Error(),
		Type:          err.Type(),
		Code:          err.Code(),
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
