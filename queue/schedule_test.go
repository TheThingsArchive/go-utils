// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package queue

import (
	"sync"
	"testing"

	"time"

	. "github.com/smartystreets/assertions"
)

func TestSchedulingConflicts(t *testing.T) {
	a := New(t)

	base := new(scheduleItemWithTimestamp)
	base.time = time.Now()
	base.timestamp = base.time.UnixNano()
	base.duration = 10 * time.Second

	var testTime = func(startOffset, duration time.Duration) bool {
		obj := new(scheduleItem)
		obj.time = base.time.Add(startOffset * time.Second)
		obj.duration = duration * time.Second
		return conflict(base, obj)
	}

	var testTimestamp = func(startOffset, duration time.Duration) bool {
		obj := new(scheduleItemWithTimestamp)
		obj.time = base.time.Add(startOffset * time.Second)
		obj.timestamp = obj.time.UnixNano()
		obj.duration = duration * time.Second
		return conflict(base, obj)
	}

	for _, test := range []func(time.Duration, time.Duration) bool{
		testTime, testTimestamp,
	} {
		a.So(test(-10, 20), ShouldBeTrue)
		a.So(test(-10, 1), ShouldBeFalse)
		a.So(test(-1, 2), ShouldBeTrue)
		a.So(test(4, 2), ShouldBeTrue)
		a.So(test(9, 2), ShouldBeTrue)
		a.So(test(20, 1), ShouldBeFalse)
	}

	{
		q := NewSchedule()
		a.So(q.Schedule("0", base.time, base.duration), ShouldHaveLength, 0)
		q.Next()
		a.So(q.Conflicts(base.time.Add(base.duration+1), base.duration), ShouldHaveLength, 0)
		a.So(q.Schedule("1", base.time.Add(base.duration+1), base.duration), ShouldHaveLength, 0)
		a.So(q.Conflicts(base.time.Add(base.duration-1), base.duration), ShouldHaveLength, 2)
		a.So(q.Schedule("2", base.time.Add(base.duration-1), base.duration), ShouldHaveLength, 2)
	}

	{
		q := NewSchedule()
		a.So(q.ScheduleWithTimestamp("0", base.time, 0, 10), ShouldHaveLength, 0)
		a.So(q.ConflictsForTimestamp(11, 10), ShouldHaveLength, 0)
		a.So(q.ScheduleWithTimestamp("1", base.time, 11, 10), ShouldHaveLength, 0)
		a.So(q.ConflictsForTimestamp(9, 10), ShouldHaveLength, 2)
		a.So(q.ScheduleWithTimestamp("2", base.time, 9, 10), ShouldHaveLength, 2)
	}

}

func TestScheduleQueue(t *testing.T) {
	waitTime := 20 * time.Millisecond

	a := New(t)
	q := NewSchedule()

	a.So(q.IsEmpty(), ShouldBeTrue)

	now := time.Now()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		a.So(q.Next(), ShouldEqual, "first")
		a.So(time.Now(), ShouldHappenWithin, 5*time.Millisecond, now.Add(waitTime*2))
		wg.Done()
	}()

	time.Sleep(waitTime / 2)
	q.Schedule("last", now.Add(waitTime*3), waitTime)
	time.Sleep(waitTime)
	q.Schedule("first", now.Add(waitTime*2), waitTime)
	time.Sleep(waitTime / 2)
	q.Schedule("second", now.Add(waitTime*2), waitTime)

	wg.Wait()
	wg.Add(2)

	go func() {
		q.Next()
		wg.Done()
	}()
	go func() {
		q.Next()
		wg.Done()
	}()

	wg.Wait()
	wg.Add(1)

	go func() {
		a.So(q.Next(), ShouldBeNil)
		wg.Done()
	}()

	q.ScheduleWithTimestamp("useless", now.Add(waitTime*4), now.Add(waitTime*4).UnixNano(), waitTime)

	q.Destroy()

	q.Schedule("too late", time.Now(), waitTime) // nothing should happen
}

func TestScheduleASAP(t *testing.T) {
	block := 2 * time.Second

	a := New(t)

	{
		q := NewSchedule()
		a.So(q.ScheduleASAP("1", block), ShouldHappenWithin, time.Millisecond, time.Now())
		a.So(q.ScheduleASAP("2", block), ShouldHappenWithin, time.Millisecond, time.Now().Add(block))
	}

	{
		q := NewSchedule()
		q.Schedule("0", time.Now().Add(-1*block), block*2)                                                    // from -1 to 1 block
		q.Next()                                                                                              // removed from schedule, but still a conflict
		q.Schedule("1", time.Now(), block)                                                                    // from 0 to 1 block
		q.Schedule("2", time.Now().Add(block/2), block)                                                       // from 0.5 to 1.5 block
		q.Schedule("4", time.Now().Add(block*3), block)                                                       // from 3 to 4 block
		a.So(q.ScheduleASAP("3", block), ShouldHappenWithin, time.Millisecond, time.Now().Add(block/2+block)) // from 1.5 to 2.5 block
	}

}
