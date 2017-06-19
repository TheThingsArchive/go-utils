package configuration

import "reflect"

// InvalidBindError is an error that occurs when trying to Bind
// an invalid struct
type InvalidBindError struct {
	typ reflect.Type
}

// Error implements error
func (e *InvalidBindError) Error() string {
	if e.typ == nil {
		return "configuration: Bind(nil)"
	}

	if e.typ.Kind() != reflect.Ptr {
		return "configuration: Bind(non-pointer " + e.typ.String() + ")"
	}

	return "configuration: Bind(nil " + e.typ.String() + ")"
}

// UnsupportedTypeError occurs when defining a configuration
// variable with a type that is not supported
type UnsupportedTypeError struct {
	Type reflect.Type
}

// Error implements error
func (e *UnsupportedTypeError) Error() string {
	return "configuration: unsupported type " + e.Type.String()
}
