package errors

// Error is the interface og grpc errors
type Error interface {
	error

	// Code returns the error code
	Code() Code

	// Type returns the error type
	Type() Type

	// Attributes returns the error attributes
	Attributes() Attributes
}

// Code represents a unique error code
type Code uint32

// Attributes
type Attributes map[string]interface{}
