// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package streambuffer_test

import (
	"context"
	"net"
	"testing"
	"time"

	. "github.com/TheThingsNetwork/go-utils/grpc/internal/test"
	"github.com/TheThingsNetwork/go-utils/grpc/streambuffer"
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/test"
	s "github.com/smartystreets/assertions"
	"google.golang.org/grpc"
)

const sleepTime = 20 * time.Millisecond

var (
	opts              []grpc.DialOption
	addr              string
	errIsNotRetryable bool
	backoffTime       time.Duration
)

func Example() {
	// Set up gRPC as usual
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	testClient := NewTestClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	syncBuffer := streambuffer.New(5, func() (grpc.ClientStream, error) {
		// maybe extend the context?
		return testClient.Push(ctx)
	})

	// From now on the syncBuffer starts buffering
	syncBuffer.SendMsg(&Foo{})

	// The Run func actually starts flushing the buffer out to the stream
	errCh := make(chan error)
	go func() {
		for {
			err := syncBuffer.Run()
			if err == nil || err == context.Canceled || errIsNotRetryable {
				errCh <- err
				close(errCh)
				return
			}
			time.Sleep(backoffTime)
		}
	}()

	// Just keep sending those messages
	syncBuffer.SendMsg(&Foo{})
}

func TestStreamBuffer(t *testing.T) {
	a := s.New(t)

	testLogger := test.NewLogger()
	log.Set(testLogger)
	defer testLogger.Print(t)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	server := NewTestServer()
	rpc := grpc.NewServer()
	RegisterTestServer(rpc, server)
	go func() {
		if err := rpc.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli := NewTestClient(conn)

	ctx, cancel := context.WithCancel(context.Background())

	var setupCalled bool
	buf := streambuffer.New(5, func() (grpc.ClientStream, error) {
		setupCalled = true
		return cli.Sync(ctx)
	})

	recv := buf.Recv(func() interface{} {
		return new(Bar)
	})
	defer buf.CloseRecv()

	var runReturned bool
	var runErr error
	go func() {
		if err := buf.Run(); err != nil {
			runReturned = true
			runErr = err
		}
	}()

	time.Sleep(sleepTime)
	a.So(setupCalled, s.ShouldBeTrue)

	buf.SendMsg(&Foo{Foo: "foo"})

	time.Sleep(sleepTime)
	a.So(server.SyncFoo, s.ShouldNotBeNil)
	a.So(server.SyncFoo.Foo, s.ShouldEqual, "foo")
	a.So(recv, s.ShouldNotBeEmpty)
	a.So((<-recv).(*Bar).Bar, s.ShouldEqual, "foo")

	cancel()

	time.Sleep(sleepTime)
	a.So(runReturned, s.ShouldBeTrue)
	a.So(runErr, s.ShouldEqual, context.Canceled)

	sent, dropped := buf.Stats()
	a.So(sent, s.ShouldEqual, 1)
	a.So(dropped, s.ShouldEqual, 0)
}
