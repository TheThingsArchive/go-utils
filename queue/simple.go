// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package queue

import "sync"

// Simple Queue implementation
type Simple interface {
	Base

	// Add an item to the Queue
	Add(interface{})
}

type simpleQueue struct {
	mu        sync.Mutex
	queue     []interface{}
	available *sync.Cond
}

// NewSimple returns a new Simple Queue
func NewSimple() Simple {
	q := &simpleQueue{
		queue: make([]interface{}, 0),
	}
	q.available = sync.NewCond(&q.mu)
	return q
}

func (q *simpleQueue) Add(i interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	wasEmpty := q.isEmpty()
	q.queue = append(q.queue, i)
	if wasEmpty {
		q.available.Signal()
	}
}

func (q *simpleQueue) Next() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.isEmpty() {
		q.available.Wait()
	}
	if q.isEmpty() {
		return nil
	}
	i := q.queue[0]
	q.queue = q.queue[1:]

	return i
}

func (q *simpleQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.isEmpty()
}

func (q *simpleQueue) isEmpty() bool {
	return len(q.queue) == 0
}

func (q *simpleQueue) Destroy() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = make([]interface{}, 0)
	q.available.Broadcast()
}
