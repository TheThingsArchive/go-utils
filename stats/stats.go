// Copyright © 2016 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package stats

import (
	"runtime"
	"time"

	"github.com/TheThingsNetwork/go-utils/log"
)

var megaByte = float64(1024 * 1024)

// Start starts the stat process that will log relevant memory-related stats
// to ctx, at an interval determined by interval.
func Start(ctx log.Interface, interval time.Duration) {
	ctx.WithField("interval", interval).Debug("starting stats loop")
	go func() {
		memstats := new(runtime.MemStats)
		for range time.Tick(interval) {
			runtime.ReadMemStats(memstats)
			ctx.WithFields(log.Fields{
				"goroutines": runtime.NumGoroutine(),
				"memory":     float64(memstats.Alloc) / megaByte, // MegaBytes allocated and not yet freed
			}).Debugf("memory stats")
		}
	}()
}
