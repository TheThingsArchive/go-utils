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
	NumEntries      = 100
	MaxWriters      = 2
)

func TestBatchWriter(t *testing.T) {
	a := s.New(t)
	mock := newMockBatchPointWriter(a)
	w := NewBatchingWriter(ttnlog.Get(), mock, ScalingInterval, MaxWriters)

	a.So(w.maxWriters, s.ShouldEqual, MaxWriters)
	closeCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(ScalingInterval):
				a.So(w.activeWriters, s.ShouldBeLessThanOrEqualTo, MaxWriters+1)
			case <-closeCh:
				return
			}
		}
	}()
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
	close(closeCh)

	a.So(mock.results, s.ShouldHaveLength, len(expected))
	for p := range expected {
		a.So(mock.results, s.ShouldContainKey, p)
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
