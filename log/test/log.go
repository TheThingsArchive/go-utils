// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package test

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"

	"github.com/TheThingsNetwork/go-utils/log"
	wrap "github.com/TheThingsNetwork/go-utils/log/apex"
	apex "github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

// Logger for testing
type Logger struct {
	logs bytes.Buffer
	log.Interface
}

// Print the logs and reset the buffer
func (l *Logger) Print(t *testing.T) {
	var loc string
	if _, file, line, ok := runtime.Caller(1); ok {
		loc = fmt.Sprintf("%s:%d", file, line)
	}
	len := l.logs.Len()
	if len > 0 {
		logs := l.logs.Next(len)
		t.Log("Logs " + loc + ": \n" + string(logs))
	}
}

// NewLogger creates a new test logger
func NewLogger() *Logger {
	l := new(Logger)
	l.Interface = wrap.Wrap(&apex.Logger{
		Handler: text.New(&l.logs),
		Level:   apex.DebugLevel,
	})
	return l
}
