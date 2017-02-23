// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package queue

import (
	"sync"
	"testing"

	"time"

	. "github.com/smartystreets/assertions"
)

func TestJITQueue(t *testing.T) {
	waitTime := 20 * time.Millisecond

	a := New(t)
	q := NewJIT()

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
	q.Schedule("last", now.Add(waitTime*3))
	time.Sleep(waitTime)
	q.Schedule("first", now.Add(waitTime*2))
	time.Sleep(waitTime / 2)
	q.Schedule("second", now.Add(waitTime*2))

	wg.Wait() // Wait for the first

	a.So(q.Next(), ShouldEqual, "second")
	a.So(q.Next(), ShouldEqual, "last")

	q.Schedule("concurrent A", now.Add(waitTime*4))
	q.Schedule("concurrent B", now.Add(waitTime*4))

	wg.Add(2)
	go func() {
		a.So(q.Next(), ShouldStartWith, "concurrent")
		wg.Done()
	}()
	go func() {
		a.So(q.Next(), ShouldStartWith, "concurrent")
		wg.Done()
	}()
	wg.Wait() // wait for the concurrent ones

	wg.Add(1)
	go func() {
		a.So(q.Next(), ShouldBeNil)
		wg.Done()
	}()

	q.Destroy()

	wg.Wait()

	q.Schedule("too late", time.Now()) // nothing should happen
}
