// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"runtime"

	"google.golang.org/grpc/codes"
)

type base struct {
	source  string
	cause   error
	code    codes.Code
	message string
	fields  map[string]interface{}
}

func newBase(text string) base {
	err := base{message: text}
	if _, file, line, ok := runtime.Caller(2); ok {
		err.source = fmt.Sprintf("%s:%d", file, line)
	}
	return err
}

// New error
func New(text string) error {
	return newBase(text)
}

func (b base) Cause() error {
	if b.cause != nil {
		return FindCause(b.cause)
	}
	return b
}

func (b base) Code() codes.Code {
	if b.code != 0 {
		return b.code
	}
	if b.cause != nil {
		return FindCode(b.cause)
	}
	return codes.Unknown
}

func (b base) Error() string {
	msg := b.message
	if b.code != 0 {
		msg += fmt.Sprintf(" [%s]", b.code)
	}
	if b.cause != nil {
		msg += ": " + b.cause.Error()
	}
	return msg
}

func (b base) Fields() (fields Fields) {
	if b.cause != nil {
		if err, ok := b.cause.(hasFields); ok {
			fields = err.Fields()
		}
	}
	if fields == nil {
		fields = make(Fields)
	}
	for k, v := range b.fields {
		fields[k] = v
	}
	return
}
