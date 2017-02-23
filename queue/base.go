// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

// Package queue implements different kinds of queues
package queue

// Base interface for Queue
type Base interface {
	// Next item in the Queue, this function blocks until the next item is available.
	// It returns <nil> to all callers when Destroy() is called.
	Next() interface{}

	IsEmpty() bool

	// Destroy the queue
	Destroy()
}
