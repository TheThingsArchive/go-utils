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

	base := new(scheduleItem)
	base.time = time.Now()
	base.duration = 10 * time.Second

	var test = func(startOffset, duration time.Duration) bool {
		obj := new(scheduleItem)
		obj.time = base.time.Add(startOffset * time.Second)
		obj.duration = duration * time.Second
		return conflict(base, obj)
	}

	a.So(test(-10, 20), ShouldBeTrue)
	a.So(test(-10, 1), ShouldBeFalse)
	a.So(test(-1, 2), ShouldBeTrue)
	a.So(test(4, 2), ShouldBeTrue)
	a.So(test(9, 2), ShouldBeTrue)
	a.So(test(20, 1), ShouldBeFalse)

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

	q.Clean()

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
		q.Schedule("1", time.Now(), block)
		q.Schedule("2", time.Now().Add(block/2), block)
		q.Schedule("4", time.Now().Add(block*3), block)
		a.So(q.ScheduleASAP("3", block), ShouldHappenWithin, time.Millisecond, time.Now().Add(block/2+block))
	}

}
