// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package influx

import (
	"fmt"
	"sync"
	"time"

	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	influxdb "github.com/influxdata/influxdb/client/v2"
)

// DefaultScalingInterval represents default scaling interval(time to wait before a new writer is started/killed) used by batching writer.
const DefaultScalingInterval = 500 * time.Millisecond

// newBatchPoints creates new influxdb.BatchPoints with specified bpConf.
// Panics on errors.
func newBatchPoints(bpConf influxdb.BatchPointsConfig) influxdb.BatchPoints {
	bp, err := influxdb.NewBatchPoints(bpConf)
	if err != nil {
		// Can only happen if there's an error in the code
		panic(fmt.Errorf("Invalid batch point configuration: %s", err))
	}
	return bp
}

// BatchPointsWriter writes influxdb.BatchPoints to Influx database.
type BatchPointsWriter interface {
	Write(bp influxdb.BatchPoints) error
}

// PointWriter writes *influxdb.Point to Influx database.
type PointWriter interface {
	Write(bpConf influxdb.BatchPointsConfig, p *influxdb.Point) error
}

// SinglePointWriter is a PointWriter, which writes points one-by-one
type SinglePointWriter struct {
	log    ttnlog.Interface
	writer BatchPointsWriter
}

// NewSinglePointWriter creates new SinglePointWriter
func NewSinglePointWriter(log ttnlog.Interface, w BatchPointsWriter) *SinglePointWriter {
	return &SinglePointWriter{
		log:    log,
		writer: w,
	}
}

// Write creates new influxdb.BatchPoints containing p and delegates that to the writer
func (w *SinglePointWriter) Write(bpConf influxdb.BatchPointsConfig, p *influxdb.Point) error {
	bp := newBatchPoints(bpConf)
	bp.AddPoint(p)
	return w.writer.Write(bp)
}

// batchPoint is a respresenation of a point written by BatchingWriter.
// Result of writing the batch of points containing the wrapped *influxdb.Point must be reported on errch.
type batchPoint struct {
	*influxdb.Point
	errch chan error
}

// pushError reports err(can be nil) to the waiter.
// It does not block and must only be executed once.
func (p *batchPoint) pushError(err error) {
	p.errch <- err
	close(p.errch)
}

func writeInBatches(log ttnlog.Interface, w BatchPointsWriter, bpConf influxdb.BatchPointsConfig, scalingInterval time.Duration, ch <-chan *batchPoint) {
	log = log.WithField("config", bpConf)

	var points []*batchPoint
	for {
		select {
		case p := <-ch:
			points = append(points, p)
		default:
			if len(points) == 0 {
				select {
				case p := <-ch:
					points = append(points, p)
					continue
				case <-time.After(scalingInterval):
					log.Info("Removing batch writer instance")
					return
				}
			}

			bp := newBatchPoints(bpConf)
			for _, p := range points {
				bp.AddPoint(p.Point)
			}

			log.WithField("num", len(points)).Debug("Writing batch of points to Influx")
			err := w.Write(bp)
			for _, p := range points {
				go p.pushError(err)
			}
			points = points[:0]
		}
	}
}

// BatchingWriter is a PointWriter, which writes points in batches.
// BatchingWriter scales automatically once it notices a delay of scalingInterval to write a batch of points and downscales if no points are supplied to an instance for a duration of scalingInterval.
// maxWriters specify the maximum amount of additional instances spawned by the batching writer.
// BatchingWriter spawns an instance for each unique BatchPointsConfig specified and up to maxWriters additional instances on top of that.
// Each instance of the writer is spawned in a separate goroutine.
type BatchingWriter struct {
	log             ttnlog.Interface
	writer          BatchPointsWriter
	scalingInterval time.Duration

	activeWriterMutex sync.RWMutex
	activeWriters     uint
	maxWriterMutex    sync.RWMutex
	maxWriters        uint

	pointChanMutex sync.RWMutex
	pointChans     map[influxdb.BatchPointsConfig]chan *batchPoint
}

// NewBatchingWriter creates new BatchingWriter.
func NewBatchingWriter(log ttnlog.Interface, w BatchPointsWriter, scalingInterval time.Duration, maxWriters uint) *BatchingWriter {
	return &BatchingWriter{
		log:             log,
		writer:          w,
		maxWriters:      maxWriters,
		scalingInterval: scalingInterval,
		pointChans:      make(map[influxdb.BatchPointsConfig]chan *batchPoint),
	}
}

// Write delegates p to a running instance of BatchingWriter and spawns new instances as required.
func (w *BatchingWriter) Write(bpConf influxdb.BatchPointsConfig, p *influxdb.Point) error {
	w.pointChanMutex.RLock()
	ch, ok := w.pointChans[bpConf]
	w.pointChanMutex.RUnlock()
	if !ok {
		w.pointChanMutex.Lock()
		ch, ok = w.pointChans[bpConf]
		if !ok {
			w.activeWriterMutex.Lock()
			w.activeWriters++
			w.activeWriterMutex.Unlock()

			w.maxWriterMutex.Lock()
			w.maxWriters++
			w.maxWriterMutex.Unlock()

			ch = make(chan *batchPoint)
			w.pointChans[bpConf] = ch
			go writeInBatches(w.log, w.writer, bpConf, w.scalingInterval, ch)
		}
		w.pointChanMutex.Unlock()
	}

	point := &batchPoint{
		Point: p,
		errch: make(chan error, 1),
	}
	select {
	case ch <- point:
	case <-time.After(w.scalingInterval):
		w.maxWriterMutex.RLock()

		w.activeWriterMutex.RLock()
		spawnNew := w.maxWriters != 0 && w.activeWriters < w.maxWriters
		w.activeWriterMutex.RUnlock()

		if spawnNew {
			w.activeWriterMutex.Lock()
			if w.activeWriters < w.maxWriters {
				w.activeWriters++
				w.log.WithFields(ttnlog.Fields{
					"config":  bpConf,
					"writers": w.activeWriters,
				}).Info("Creating additional batch writer instance")
				go writeInBatches(w.log, w.writer, bpConf, w.scalingInterval, ch)
			}
			w.activeWriterMutex.Unlock()
		}
		w.maxWriterMutex.RUnlock()

		ch <- point
	}
	return <-point.errch
}
