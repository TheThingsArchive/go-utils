package errors

// impl implements Error
type impl struct {
	message    string
	code       Code
	typ        Type
	attributes Attributes
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
