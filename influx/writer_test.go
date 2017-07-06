// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package influx

import (
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/TheThingsNetwork/go-utils/handlers/cli"
	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	ttnapex "github.com/TheThingsNetwork/go-utils/log/apex"
	"github.com/TheThingsNetwork/ttn/utils/random"
	apex "github.com/apex/log"
	"github.com/fortytw2/leaktest"
	influxdb "github.com/influxdata/influxdb/client/v2"
	s "github.com/smartystreets/assertions"
)

func init() {
	ttnlog.Set(ttnapex.Wrap(&apex.Logger{
		Level:   apex.DebugLevel,
		Handler: cli.New(os.Stdout),
	}))
}

type MockBatchPointWriter struct {
	assertion *s.Assertion

	sync.RWMutex
	results map[*influxdb.Point]error
}

func newMockBatchPointWriter(a *s.Assertion) *MockBatchPointWriter {
	return &MockBatchPointWriter{
		assertion: a,
		results:   make(map[*influxdb.Point]error),
	}
}

func (w *MockBatchPointWriter) Write(bp influxdb.BatchPoints) error {
	time.Sleep(ScalingInterval)
	var err error
	if random.Bool() {
		err = errors.New("test")
	}
	for _, p := range bp.Points() {
		w.Lock()
		if w.assertion.So(w.results, s.ShouldNotContainKey, p) {
			w.results[p] = err
		}
		w.Unlock()
	}
	return err
}

const (
	ScalingInterval = time.Millisecond
	NumEntries      = 200
)

func TestBatchWriter(t *testing.T) {
	a := s.New(t)
	for _, mw := range []int{
		-1, 0, 100,
	} {
		mock := newMockBatchPointWriter(a)

		var w *BatchingWriter
		if mw < 0 {
			w = NewBatchingWriter(ttnlog.Get(), mock, WithScalingInterval(ScalingInterval))
		} else {
			w = NewBatchingWriter(ttnlog.Get(), mock, WithScalingInterval(ScalingInterval), WithInstanceLimit(uint(mw)))
			a.So(w.limit, s.ShouldEqual, mw)
		}

		checkCh := make(chan struct{})

		wg := &sync.WaitGroup{}
		expected := make(map[*influxdb.Point]bool)
		for i := 0; i < NumEntries; i++ {
			p := &influxdb.Point{}
			expected[p] = true
			if i == 0 {
				// one goroutine is spawned after the inital write, it should stay alive forever
				err := w.Write(influxdb.BatchPointsConfig{}, p)
				a.So(err, s.ShouldEqual, mock.results[p])
				go func() {
					defer leaktest.Check(t)()
					checkCh <- struct{}{}
					if mw < 0 {
						<-checkCh
						return
					}

					for {
						select {
						case <-time.After(ScalingInterval):
							if mw == 0 {
								a.So(w.active, s.ShouldEqual, 1)
							}
							if mw > 0 {
								a.So(w.active, s.ShouldBeBetweenOrEqual, 1, mw+1)
							}
						case <-checkCh:
							return
						}
					}
				}()
				// wait for leaktest to count active goroutines
				<-checkCh
				continue
			}

			wg.Add(1)
			go func() {
				err := w.Write(influxdb.BatchPointsConfig{}, p)

				mock.RLock()
				a.So(err, s.ShouldEqual, mock.results[p])
				mock.RUnlock()
				wg.Done()
			}()
		}
		wg.Wait()
		close(checkCh)

		a.So(mock.results, s.ShouldHaveLength, len(expected))
		for p := range expected {
			a.So(mock.results, s.ShouldContainKey, p)
		}
	}
}

func TestSinglePointWriter(t *testing.T) {
	a := s.New(t)
	mock := newMockBatchPointWriter(a)
	w := NewSinglePointWriter(ttnlog.Get(), mock)
	wg := &sync.WaitGroup{}
	expected := make(map[*influxdb.Point]bool)
	for i := 0; i < NumEntries; i++ {
		wg.Add(1)
		p := &influxdb.Point{}
		expected[p] = true
		go func() {
			err := w.Write(influxdb.BatchPointsConfig{}, p)
			mock.RLock()
			a.So(err, s.ShouldEqual, mock.results[p])
			mock.RUnlock()
			wg.Done()
		}()
	}
	wg.Wait()

	a.So(mock.results, s.ShouldHaveLength, len(expected))
	for p := range expected {
		a.So(mock.results, s.ShouldContainKey, p)
	}
}
