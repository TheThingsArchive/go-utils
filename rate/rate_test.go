// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package rate

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	redis "gopkg.in/redis.v5"
)

func getRedisClient() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", host),
		Password: "",
		DB:       2,
	})
}

var keys int

func getKey() string {
	keys++
	return fmt.Sprintf("test-counter:%d", keys)
}

func TestCounter(t *testing.T) {
	retention := 10 * time.Second
	counters := map[string]func() Counter{
		"*counter":      func() Counter { return NewCounter(time.Second, retention) },
		"*redisCounter": func() Counter { return NewRedisCounter(getRedisClient(), getKey(), time.Second, retention) },
	}

	defer func() {
		for key := 1; key <= keys; key++ {
			getRedisClient().Del(fmt.Sprintf("test-counter:%d", key))
		}
	}()

	for typ, counterFunc := range counters {
		Convey(fmt.Sprintf("Given a new %s Counter", typ), t, func(c C) {
			l := counterFunc()
			now := time.Now()
			Convey("When adding an event", func() {
				l.Add(now, 1)
				Convey("Then getting the events should return 1", func() {
					events, err := l.Get(now, retention)
					So(err, ShouldBeNil)
					So(events, ShouldEqual, 1)
				})
			})
			Convey("When adding \"i\" events per second for 20 seconds", func() {
				for i := 1; i <= 20; i++ {
					l.Add(now.Add(time.Duration(i)*time.Second), uint64(i))
				}
				Convey("Then getting the events after 20 seconds should return the last 10", func() {
					events, err := l.Get(now.Add(20*time.Second), retention)
					So(err, ShouldBeNil)
					So(events, ShouldEqual, 20+19+18+17+16+15+14+13+12+11)
				})
				Convey("Then getting the events after 25 seconds should return the last 5", func() {
					events, err := l.Get(now.Add(25*time.Second), retention)
					So(err, ShouldBeNil)
					So(events, ShouldEqual, 20+19+18+17+16)
				})
				Convey("Then getting the events of the past 20 seconds after 25 seconds should still return the last 5", func() {
					events, err := l.Get(now.Add(25*time.Second), retention*2)
					So(err, ShouldBeNil)
					So(events, ShouldEqual, 20+19+18+17+16)
				})
				Convey("Then getting the events after 35 seconds should return 0", func() {
					if redis, ok := l.(*redisCounter); ok {
						redis.client.Del(redis.key) // We have to manually expire here
					}
					events, err := l.Get(now.Add(35*time.Second), retention)
					So(err, ShouldBeNil)
					So(events, ShouldEqual, 0)
				})
			})
		})
	}
}

func TestLimiter(t *testing.T) {
	Convey("Given a new Limiter", t, func(c C) {
		retention := 10 * time.Second
		l := NewLimiter(NewCounter(time.Second, retention), retention, 10)
		Convey("The first 10 calls to Limit() should return false", func() {
			for i := 1; i <= 10; i++ {
				limit, err := l.Limit()
				So(err, ShouldBeNil)
				So(limit, ShouldBeFalse)
			}
			Convey("The next call should return true", func() {
				limit, err := l.Limit()
				So(err, ShouldBeNil)
				So(limit, ShouldBeTrue)
			})
		})
	})
}

func BenchmarkCounter(b *testing.B) {
	l := NewCounter(time.Second, 10*time.Second)
	now := time.Now()
	for i := 0; i < b.N; i++ {
		t := now.Add(time.Duration(i) * 100 * time.Millisecond)
		l.Add(t, 1)
		l.Get(t, 10*time.Second)
	}
}
