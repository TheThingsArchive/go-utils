// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package restartstream

import (
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	. "github.com/TheThingsNetwork/go-utils/grpc/internal/test"
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/htdvisser/grpc-testing/test"
	grpc_middleware "github.com/mwitkow/go-grpc-middleware"
	. "github.com/smartystreets/assertions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const sleepTime = 20 * time.Millisecond

func TestReconnect(t *testing.T) {
	a := New(t)

	testLogger := test.NewLogger()
	log.Set(testLogger)
	defer testLogger.Print(t)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	_, port, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		t.Fatalf("Failed to parse listener address: %v", err)
	}
	s := grpc.NewServer()

	server := NewTestServer()

	RegisterTestServer(s, server)
	go s.Serve(lis)

	addr := "localhost:" + port
	breakStream := NewCancel()
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(Interceptor, breakStream.Interceptor)))
	if err != nil {
		t.Fatalf("Dial(%q) = %v", addr, err)
	}
	cli := NewTestClient(conn)

	{
		res, err := cli.Get(context.Background(), &Foo{Foo: "ok"})
		a.So(err, ShouldBeNil)
		a.So(res, ShouldNotBeNil)
		a.So(server.GetFoo, ShouldNotBeNil)
		a.So(server.GetFoo.Foo, ShouldEqual, "ok")

		testLogger.Print(t)
	}

	var testPush = func(doCancel bool) {
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := cli.Push(ctx)
		a.So(err, ShouldBeNil)

		var done bool
		go func() {
			for {
				bar := new(Bar)
				err := stream.RecvMsg(bar)
				if err == io.EOF || grpc.Code(err) == codes.Canceled {
					log.Get().WithField("Method", "Push").WithError(err).Debugf("[TEST] EOF")
					done = true
					return
				}
				if err == nil {
					log.Get().WithField("Method", "Push").WithField("Bar", bar).Debugf("[TEST] Recv Ok")
				} else {
					log.Get().WithField("Method", "Push").WithField("Bar", bar).WithError(err).Debugf("[TEST] Recv Err")
				}
			}
		}()

		err = stream.Send(&Foo{Foo: "ok"})
		a.So(err, ShouldBeNil)
		time.Sleep(sleepTime)
		a.So(server.PushFoo, ShouldNotBeNil)
		a.So(server.PushFoo.Foo, ShouldEqual, "ok")

		// not ok breaks the stream
		err = stream.Send(&Foo{Foo: "not ok"})
		a.So(err, ShouldBeNil)
		time.Sleep(sleepTime)

		err = stream.Send(&Foo{Foo: "ok again"})
		a.So(err, ShouldBeNil)
		time.Sleep(sleepTime)
		a.So(server.PushFoo, ShouldNotBeNil)
		a.So(server.PushFoo.Foo, ShouldEqual, "ok again")
		time.Sleep(sleepTime)

		// break the stream
		breakStream.Cancel()
		time.Sleep(sleepTime)

		err = stream.Send(&Foo{Foo: "and again"})
		a.So(err, ShouldBeNil)
		time.Sleep(sleepTime)
		a.So(server.PushFoo, ShouldNotBeNil)
		a.So(server.PushFoo.Foo, ShouldEqual, "and again")

		if doCancel {
			cancel()
		} else {
			stream.CloseSend()
		}

		time.Sleep(sleepTime)
		a.So(done, ShouldBeTrue)

		testLogger.Print(t)
	}

	testPush(false)

	testPush(true)

	var testPull = func(foo string, doCancel bool) {
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := cli.Pull(ctx, &Foo{Foo: foo})
		a.So(err, ShouldBeNil)

		var wg sync.WaitGroup
		wg.Add(5)

		go func() {
			for {
				bar, err := stream.Recv()
				if err == io.EOF || grpc.Code(err) == codes.Canceled {
					log.Get().WithField("Method", "Pull").WithError(err).Debugf("[TEST] EOF")
					return
				}
				if err == nil {
					log.Get().WithField("Method", "Pull").WithField("Bar", bar).Debugf("[TEST] Recv Ok")
				} else {
					log.Get().WithField("Method", "Pull").WithField("Bar", bar).WithError(err).Debugf("[TEST] Recv Err")
				}
				wg.Done()
			}
		}()

		time.Sleep(2 * sleepTime)

		if doCancel {
			cancel()
		} else {
			wg.Wait()
		}

		time.Sleep(sleepTime)

		testLogger.Print(t)
	}

	testPull("not ok", true)

	testPull("not ok after 1", true)

	testPull("ok", false)

	testPull("ok", true)

	var testSync = func(doCancel bool) {
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := cli.Sync(ctx)
		a.So(err, ShouldBeNil)

		go func() {
			for {
				bar, err := stream.Recv()
				if err == io.EOF || grpc.Code(err) == codes.Canceled {
					log.Get().WithField("Method", "Sync").WithError(err).Debugf("[TEST] EOF")
					return
				}
				if err == nil {
					log.Get().WithField("Method", "Sync").WithField("Bar", bar).Debugf("[TEST] Recv Ok")
				} else {
					log.Get().WithField("Method", "Sync").WithField("Bar", bar).WithError(err).Debugf("[TEST] Recv Err")
				}
			}
		}()

		err = stream.Send(&Foo{Foo: "ok"})
		a.So(err, ShouldBeNil)
		time.Sleep(sleepTime)
		a.So(server.SyncFoo, ShouldNotBeNil)
		a.So(server.SyncFoo.Foo, ShouldEqual, "ok")

		// not ok breaks the stream
		err = stream.Send(&Foo{Foo: "not ok"})
		a.So(err, ShouldBeNil)
		time.Sleep(sleepTime)

		err = stream.Send(&Foo{Foo: "ok again"})
		a.So(err, ShouldBeNil)
		time.Sleep(sleepTime)
		a.So(server.SyncFoo, ShouldNotBeNil)
		a.So(server.SyncFoo.Foo, ShouldEqual, "ok again")
		time.Sleep(sleepTime)

		time.Sleep(2 * sleepTime)

		if doCancel {
			cancel()
		} else {
			stream.CloseSend()
		}

		time.Sleep(sleepTime)

		testLogger.Print(t)
	}

	testSync(true)
	testSync(false)

}
