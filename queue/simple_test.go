// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package queue

import (
	"sync"
	"testing"

	"time"

	. "github.com/smartystreets/assertions"
)

const sleepTime = 5 * time.Millisecond

func TestSimpleQueue(t *testing.T) {
	a := New(t)
	q := NewSimple()

	a.So(q.IsEmpty(), ShouldBeTrue)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		a.So(q.Next(), ShouldEqual, 1)
		wg.Done()
	}()

	time.Sleep(sleepTime)

	q.Add(1)

	wg.Wait()

	q.Add(2)
	q.Add(3)

	a.So(q.IsEmpty(), ShouldBeFalse)
	a.So(q.Next(), ShouldEqual, 2)
	a.So(q.Next(), ShouldEqual, 3)

	wg.Add(1)
	go func() {
		a.So(q.Next(), ShouldBeNil)
		wg.Done()
	}()

	time.Sleep(sleepTime)

	q.Clean()

	wg.Wait()

}
