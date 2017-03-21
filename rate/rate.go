// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

// Package rate implements rate counting and rate limiting.
package rate

import (
	"strconv"
	"sync"
	"time"

	redis "gopkg.in/redis.v5"
)

// Counter interface used in rate limiter
type Counter interface {
	Add(timestamp time.Time, events uint64) error
	Get(timestamp time.Time, past time.Duration) (events uint64, err error)
}

// NewCounter returns a new rate counter with the given bucket size and retention
func NewCounter(bucketSize, retention time.Duration) Counter {
	return &counter{
		bucketSize: bucketSize,
		retention:  retention,
		buckets:    make([]uint64, 2*retention/bucketSize),
	}
}

type counter struct {
	bucketSize time.Duration
	retention  time.Duration

	mu      sync.Mutex
	cleared time.Time
	buckets []uint64
}

func (c *counter) expire(timestamp time.Time) {
	// Don't need to clear if already cleared
	if c.cleared.After(timestamp) || c.bucket(timestamp) == c.bucket(c.cleared) {
		return
	}

	defer func() { c.cleared = timestamp }()

	// Initialized empty
	if c.cleared.IsZero() {
		return
	}

	// Fully expired
	if timestamp.Sub(c.cleared) > c.retention {
		for i := range c.buckets {
			c.buckets[i] = 0
		}
		return
	}

	// Partially expired
	for t := c.cleared.Add(c.retention); t.Before(timestamp.Add(c.retention)); t = t.Add(c.bucketSize) {
		c.buckets[c.bucket(t)] = 0
	}
}

func (c *counter) bucket(timestamp time.Time) int {
	return int(timestamp.UnixNano() % int64(2*c.retention) / int64(c.bucketSize))
}

func (c *counter) Add(now time.Time, events uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.expire(now)
	c.buckets[c.bucket(now)] += events
	return nil
}

func (c *counter) Get(now time.Time, past time.Duration) (events uint64, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if past > c.retention || past == 0 {
		past = c.retention
	}
	c.expire(now)
	for t := now; t.After(now.Add(-1 * past)); t = t.Add(-1 * c.bucketSize) {
		events += c.buckets[c.bucket(t)]
	}
	return events, nil
}

// NewRedisCounter returns a new redis-based counter
func NewRedisCounter(client *redis.Client, key string, bucketSize, retention time.Duration) Counter {
	return &redisCounter{
		client:     client,
		key:        key,
		bucketSize: bucketSize,
		retention:  retention,
	}
}

type redisCounter struct {
	client     *redis.Client
	key        string
	bucketSize time.Duration
	retention  time.Duration
}

func (c *redisCounter) bucket(timestamp time.Time) int {
	return int(timestamp.UnixNano() % int64(2*c.retention) / int64(c.bucketSize))
}

func (c *redisCounter) redisBuckets(from, to time.Time) []string {
	buckets := make([]string, 0, to.Sub(from)/c.bucketSize)
	for t := from; t.Before(to) || t.Equal(to); t = t.Add(c.bucketSize) {
		buckets = append(buckets, strconv.Itoa(c.bucket(t)))
	}
	return buckets
}

func (c *redisCounter) Add(now time.Time, events uint64) (err error) {
	bucket := c.bucket(now)
	pipe := c.client.TxPipeline()
	pipe.HDel(c.key, c.redisBuckets(now.Add(c.bucketSize), now.Add(c.retention))...)
	pipe.HIncrBy(c.key, strconv.Itoa(bucket), int64(events))
	pipe.Expire(c.key, c.retention)
	_, err = pipe.Exec()
	return err
}

func (c *redisCounter) Get(now time.Time, past time.Duration) (events uint64, err error) {
	if past > c.retention || past == 0 {
		past = c.retention
	}
	pipe := c.client.TxPipeline()
	pipe.HDel(c.key, c.redisBuckets(now.Add(c.bucketSize), now.Add(c.retention))...)
	buckets := pipe.HMGet(c.key, c.redisBuckets(now.Add(-1*past), now)...)
	_, err = pipe.Exec()
	if err != nil {
		return events, err
	}
	res, err := buckets.Result()
	for _, bucket := range res {
		if bucket == nil {
			continue
		}
		if bucket, ok := bucket.(string); ok {
			if i, err := strconv.ParseUint(string(bucket), 10, 64); err == nil {
				events += i
			}
		}
	}
	return events, err
}

// Limiter limits events
type Limiter interface {
	Limit() (limited bool, err error)
}

// NewLimiter returns a new limiter
func NewLimiter(counter Counter, duration time.Duration, limit uint64) Limiter {
	return &limiter{
		Counter:  counter,
		duration: duration,
		limit:    limit,
	}
}

type limiter struct {
	Counter
	duration time.Duration
	limit    uint64
}

func (l *limiter) Limit() (bool, error) {
	now := time.Now()
	events, err := l.Get(now, l.duration)
	if err != nil {
		return true, err
	}
	if events >= l.limit {
		return true, nil
	}
	return false, l.Add(now, 1)
}
