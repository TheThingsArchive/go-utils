// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package queue

import (
	"sort"
	"sync"
	"time"
)

// JITItem has a Time()
type JITItem interface {
	Time() time.Time
}

type jitItem struct {
	item interface{}
	time time.Time
}

func (i *jitItem) Time() time.Time {
	return i.time
}

// JIT is a just-in-time implementation of the Queue. It allows setting a time for each item in the queue. The item will
// be returned by Next() immediately after this time.
type JIT interface {
	Base

	// Add an Item to the JIT Queue, will be returned by Next() at item.Time()
	Add(item JITItem)

	// Schedule an Item to the JIT Queue, will be returned by Next() at time
	Schedule(i interface{}, time time.Time)
}

type jitSlice []JITItem

func (s jitSlice) Len() int           { return len(s) }
func (s jitSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s jitSlice) Less(i, j int) bool { return s[i].Time().Before(s[j].Time()) }

type jitQueue struct {
	nextMu sync.Mutex

	mu    sync.Mutex
	queue jitSlice

	changed chan struct{}
}

// NewJIT returns a new JIT Queue
func NewJIT() JIT {
	return &jitQueue{
		queue:   make([]JITItem, 0),
		changed: make(chan struct{}),
	}
}

func (q *jitQueue) Add(i JITItem) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.add(i)
}

func (q *jitQueue) add(i JITItem) {
	q.queue = append(q.queue, i)
	sort.Sort(q.queue)
	close(q.changed)
	q.changed = make(chan struct{})
}

func (q *jitQueue) Schedule(i interface{}, time time.Time) {
	q.Add(&jitItem{item: i, time: time})
}

func (q *jitQueue) Next() interface{} {
	q.nextMu.Lock()
	defer q.nextMu.Unlock()
	return q.next()
}

func (q *jitQueue) next() interface{} {
	var i JITItem
	for {
		q.mu.Lock()
		if q.changed == nil {
			q.mu.Unlock()
			return nil
		}
		changed := q.changed
		if !q.isEmpty() {
			i = q.queue[0]
		}
		q.mu.Unlock()

		if i == nil {
			<-changed
			continue
		}

		select {
		case <-changed:
			continue
		case <-time.After(time.Until(i.Time())):
			q.mu.Lock()
			if !q.isEmpty() && i == q.queue[0] {
				defer q.mu.Unlock()
				q.queue = q.queue[1:]
				if i, ok := i.(*jitItem); ok {
					return i.item
				}
				return i
			}
			// this is highly unlikely
			q.mu.Unlock()
		}
	}
}

func (q *jitQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.isEmpty()
}

func (q *jitQueue) isEmpty() bool {
	return len(q.queue) == 0
}

func (q *jitQueue) Clean() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = make([]JITItem, 0)
	close(q.changed)
	q.changed = nil
}
