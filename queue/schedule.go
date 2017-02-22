// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package queue

import "time"

// Schedule is an extension of the JIT Queue that allows setting a duration next to the time of an item.
// This allows to calculate the number of conflicts that an item has
type Schedule interface {
	Base

	// Add an Item to the Schedule, queued to be returned at item.Time(),
	// this func returns the number of conflicts based on item.Time() and item.Duration()
	Add(item ScheduleItem) int

	// Schedule an item at the given time, with the given duration
	// this func returns the number of conflicts based on item time and duration
	Schedule(i interface{}, time time.Time, duration time.Duration) int

	// ScheduleASAP schedules an item as soon as possible, given its duration and considering the existing Schedule
	// this func returns the time at which the item is scheduled
	ScheduleASAP(i interface{}, duration time.Duration) time.Time
}

type scheduleItem struct {
	jitItem
	duration time.Duration
}

func (s *scheduleItem) Duration() time.Duration {
	return s.duration
}

// ScheduleItem has a Time() and Duration()
type ScheduleItem interface {
	JITItem
	Duration() time.Duration
}

func conflict(i, j ScheduleItem) bool {
	// i ended before j started
	if i.Time().Add(i.Duration()).Before(j.Time()) {
		return false
	}
	// j ended before i started
	if j.Time().Add(j.Duration()).Before(i.Time()) {
		return false
	}
	return true
}

type schedule struct {
	*jitQueue
	lastFinished time.Time
}

// NewSchedule returns a new Schedule
func NewSchedule() Schedule {
	return &schedule{
		jitQueue:     NewJIT().(*jitQueue),
		lastFinished: time.Now(),
	}
}

func (q *schedule) conflicts(i ScheduleItem) (conflicts int) {
	queue := q.queue
	if i.Time().Before(q.lastFinished) {
		conflicts++
	}
	for _, qd := range queue {
		if qd, ok := qd.(ScheduleItem); ok {
			if conflict(i, qd) {
				conflicts++
			}
		}
	}
	return
}

func (q *schedule) Add(i ScheduleItem) (conflicts int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	conflicts = q.conflicts(i)
	q.jitQueue.add(i)
	return
}

func (q *schedule) Schedule(i interface{}, time time.Time, duration time.Duration) int {
	return q.Add(&scheduleItem{jitItem: jitItem{item: i, time: time}, duration: duration})
}

func (q *schedule) ScheduleASAP(i interface{}, duration time.Duration) time.Time {
	q.mu.Lock()
	defer q.mu.Unlock()
	candidate := &scheduleItem{jitItem: jitItem{item: i, time: time.Now()}, duration: duration}
	for _, qd := range q.queue {
		if qd, ok := qd.(ScheduleItem); ok {
			candidate.time = qd.Time().Add(qd.Duration() + 1)
		} else {
			continue
		}
		if q.conflicts(candidate) == 0 {
			break
		}
	}
	q.add(candidate)
	return candidate.time
}

func (q *schedule) Next() interface{} {
	q.nextMu.Lock()
	defer q.nextMu.Unlock()
	next := q.jitQueue.next()
	if next, ok := next.(ScheduleItem); ok && next.Time().After(q.lastFinished) {
		q.lastFinished = next.Time().Add(next.Duration())
	}
	if next, ok := next.(*scheduleItem); ok {
		return next.item
	}
	return next
}
