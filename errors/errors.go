package errors

import "fmt"

const cause = "cause"

type Typeable interface {
	Type() Type
}

type Error interface {
	error
	Typeable
	WithField(string, interface{}) Error
	WithCause(error) Error
	Fields() map[string]interface{}
}

type ErrorWithFields struct {
	typ     Type                     `json:"type"`
	message string                   `json:"message"`
	fields  []map[string]interface{} `json:"fields"`
}

func (e ErrorWithFields) Type() Type {
	return e.typ
}

func (e ErrorWithFields) Error() string {
	return e.message
}

func (e ErrorWithFields) WithField(name string, val interface{}) Error {
	flds := make([]map[string]interface{}, 0, len(e.fields)+1)
	flds = append(flds, map[string]interface{}{
		name: val,
	})
	flds = append(flds, e.fields...)
	return &ErrorWithFields{
		typ:    e.typ,
		fields: flds,
	}
}

func (e ErrorWithFields) WithCause(err error) Error {
	return e.WithField(cause, err)
}

func (e ErrorWithFields) Cause() error {
	// return the first error we find
	for _, fld := range e.fields {
		err, ok := fld[cause]
		if ok {
			if ce, ok := err.(error); ok {
				return ce
			}
		}
	}

	return e
}

func Err(t Type, message string) Error {
	return &ErrorWithFields{
		typ:     t,
		message: message,
	}
}

func Errf(t Type, format string, v ...interface{}) Error {
	return Err(t, fmt.Sprintf(format, v...))
}

func (e ErrorWithFields) Fields() map[string]interface{} {
	res := make(map[string]interface{})

	for _, flds := range e.fields {
		for key, val := range flds {
			res[key] = val
		}
	}

	return res
}
