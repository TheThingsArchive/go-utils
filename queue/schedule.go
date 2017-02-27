// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package queue

import "time"

// Schedule is an extension of the JIT Queue that allows setting a duration next to the time of an item.
// This allows to calculate the conflicts that an item has
type Schedule interface {
	Base

	// Add an Item to the Schedule, queued to be returned at item.Time(),
	// this func returns the conflicts based on item.Time() and item.Duration()
	Add(item ScheduleItem) []ScheduleItem

	// Conflicts based on time and duration
	Conflicts(time time.Time, duration time.Duration) []ScheduleItem

	// Conflicts based on timestamp and duration
	ConflictsForTimestamp(timestamp int64, duration time.Duration) []ScheduleItem

	// Schedule an item at the given time, with the given duration
	// this func returns the conflicts based on item time and duration
	Schedule(i interface{}, time time.Time, duration time.Duration) []ScheduleItem

	// Schedule an item at the given time+timestamp, with the given duration
	// this func returns the conflicts based on item timestamp and duration
	ScheduleWithTimestamp(i interface{}, time time.Time, timestamp int64, duration time.Duration) []ScheduleItem

	// ScheduleASAP schedules an item as soon as possible, given its duration and considering the existing Schedule
	// this func returns the time at which the item is scheduled
	ScheduleASAP(i interface{}, duration time.Duration) time.Time
}

type scheduleItem struct {
	jitItem
	duration time.Duration
}

func (s scheduleItem) Duration() time.Duration {
	return s.duration
}

// ScheduleItem has a Time() and Duration()
type ScheduleItem interface {
	JITItem
	Duration() time.Duration
}

type scheduleItemWithTimestamp struct {
	scheduleItem
	timestamp int64
}

func (s scheduleItemWithTimestamp) Timestamp() int64 {
	return s.timestamp
}

// ScheduleItemWithTimestamp has a Timestamp() and Duration()
type ScheduleItemWithTimestamp interface {
	Timestamp() int64
	Duration() time.Duration
}

// returns true if i before j
func before(i, j ScheduleItem) bool {
	iEnd := i.Time().UnixNano() + i.Duration().Nanoseconds()
	jStart := j.Time().UnixNano()
	if i, ok := i.(ScheduleItemWithTimestamp); ok {
		if j, ok := j.(ScheduleItemWithTimestamp); ok {
			iEnd = i.Timestamp() + i.Duration().Nanoseconds()
			jStart = j.Timestamp()
		}
	}
	return iEnd < jStart
}

func conflict(i, j ScheduleItem) bool {
	if before(i, j) {
		return false
	}
	if before(j, i) {
		return false
	}
	return true
}

type schedule struct {
	*jitQueue
	last ScheduleItem
}

// NewSchedule returns a new Schedule (see Schedule interface)
func NewSchedule() Schedule {
	return &schedule{
		jitQueue: NewJIT().(*jitQueue),
	}
}

func (q *schedule) Conflicts(time time.Time, duration time.Duration) []ScheduleItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.conflicts(&scheduleItem{jitItem: jitItem{time: time}, duration: duration})
}

func (q *schedule) ConflictsForTimestamp(timestamp int64, duration time.Duration) []ScheduleItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.conflicts(&scheduleItemWithTimestamp{scheduleItem: scheduleItem{duration: duration}, timestamp: timestamp})
}

func (q *schedule) conflicts(i ScheduleItem) (conflicts []ScheduleItem) {
	queue := q.queue
	if q.last != nil && conflict(i, q.last) {
		conflicts = append(conflicts, q.last)
	}
	for _, qd := range queue {
		if qd, ok := qd.(ScheduleItem); ok {
			if conflict(i, qd) {
				conflicts = append(conflicts, qd)
			}
		}
	}
	return
}

func (q *schedule) Add(i ScheduleItem) (conflicts []ScheduleItem) {
	q.mu.Lock()
	defer q.mu.Unlock()
	conflicts = q.conflicts(i)
	q.jitQueue.add(i)
	return
}

func (q *schedule) Schedule(i interface{}, time time.Time, duration time.Duration) []ScheduleItem {
	return q.Add(&scheduleItem{jitItem: jitItem{item: i, time: time}, duration: duration})
}

func (q *schedule) ScheduleWithTimestamp(i interface{}, time time.Time, timestamp int64, duration time.Duration) []ScheduleItem {
	return q.Add(&scheduleItemWithTimestamp{scheduleItem: scheduleItem{jitItem: jitItem{item: i, time: time}, duration: duration}, timestamp: timestamp})
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
		if len(q.conflicts(candidate)) == 0 {
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
	if next == nil {
		return nil
	}
	if next, ok := next.(ScheduleItem); ok && (q.last == nil || before(q.last, next)) {
		q.last = next
	}
	if next, ok := next.(hasItem); ok {
		return next.getItem()
	}
	return next
}
