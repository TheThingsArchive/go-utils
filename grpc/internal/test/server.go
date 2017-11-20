// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

//go:generate protoc --gogoslick_out=plugins=grpc:. test.proto

package test

import (
	"io"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/TheThingsNetwork/go-utils/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

type TestServerImplementation struct {
	log     log.Interface
	Token   string
	GetFoo  *Foo
	PushFoo *Foo
	PullFoo *Foo
	SyncFoo *Foo
}

func NewTestServer() *TestServerImplementation {
	return &TestServerImplementation{
		log:     log.Get(),
		GetFoo:  &Foo{Foo: "none"},
		PushFoo: &Foo{Foo: "none"},
		PullFoo: &Foo{Foo: "none"},
		SyncFoo: &Foo{Foo: "none"},
	}
}

func (s *TestServerImplementation) req(ctx context.Context) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if token, ok := md["token"]; ok && len(token) > 0 {
			s.Token = token[0]
			return
		}
	}
	s.Token = ""
}

func (s *TestServerImplementation) Get(ctx context.Context, foo *Foo) (*Bar, error) {
	s.log.WithField("Method", "Get").WithField("Foo", foo).Debugf("[SERVER] Request")
	s.req(ctx)

	if foo.GetFoo() == "not ok" {
		return nil, grpc.Errorf(codes.InvalidArgument, "Foo not ok")
	}
	s.GetFoo = foo
	return &Bar{Bar: "ok"}, nil
}

func (s *TestServerImplementation) Push(stream Test_PushServer) error {
	s.log.WithField("Method", "Push").Debugf("[SERVER] Start")
	s.req(stream.Context())

	var streamErr atomic.Value
	go func() {
		<-stream.Context().Done()
		streamErr.Store(stream.Context().Err())
	}()

	for {
		streamErr := streamErr.Load()
		if streamErr != nil {
			return streamErr.(error)
		}

		foo, err := stream.Recv()
		if err == io.EOF {
			s.log.WithField("Method", "Push").Debugf("[SERVER] EOF")
			return stream.SendAndClose(&Bar{Bar: s.PushFoo.Foo})
		}
		if err != nil {
			s.log.WithField("Method", "Push").WithError(err).Debugf("[SERVER] Error")
			return err
		}
		s.log.WithField("Method", "Push").WithField("Foo", foo).Debugf("[SERVER] Recv")
		if foo.GetFoo() == "not ok" {
			s.log.WithField("Method", "Push").WithField("Foo", foo).Debugf("[SERVER] Foo not ok")
			return grpc.Errorf(codes.InvalidArgument, "Foo not ok")
		}
		s.PushFoo = foo
	}
}

func (s *TestServerImplementation) Pull(foo *Foo, stream Test_PullServer) (err error) {
	s.log.WithField("Method", "Pull").WithField("Foo", foo).Debugf("[SERVER] Start")
	s.req(stream.Context())

	if foo.GetFoo() == "not ok" {
		return grpc.Errorf(codes.InvalidArgument, "Foo not ok")
	}

	var streamErr atomic.Value
	go func() {
		<-stream.Context().Done()
		streamErr.Store(stream.Context().Err())
	}()

	var i int
	for i < 5 {
		streamErr := streamErr.Load()
		if streamErr != nil {
			return streamErr.(error)
		}

		err = stream.Send(&Bar{Bar: strconv.Itoa(i)})
		if err != nil {
			s.log.WithField("Method", "Pull").WithError(err).Debugf("[SERVER] Error")
			return err
		}
		s.log.WithField("Method", "Pull").WithField("Bar", i).Debugf("[SERVER] Send")
		i++
		if foo.GetFoo() == "not ok after 1" {
			return grpc.Errorf(codes.InvalidArgument, "Foo not ok after 1")
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func (s *TestServerImplementation) Sync(stream Test_SyncServer) error {
	s.log.WithField("Method", "Sync").Debugf("[SERVER] Start")
	s.req(stream.Context())

	var streamErr atomic.Value
	go func() {
		<-stream.Context().Done()
		streamErr.Store(stream.Context().Err())
	}()

	var lastFoo string
	for {
		streamErr := streamErr.Load()
		if streamErr != nil {
			return streamErr.(error)
		}

		foo, err := stream.Recv()
		if err == io.EOF {
			s.log.WithField("Method", "Sync").Debugf("[SERVER] EOF")
			return nil
		}
		if err != nil {
			s.log.WithField("Method", "Sync").WithError(err).Debugf("[SERVER] Error")
			return err
		}
		s.log.WithField("Method", "Sync").WithField("Foo", foo).Debugf("[SERVER] Recv")
		if foo.GetFoo() == "not ok" {
			s.log.WithField("Method", "Sync").WithField("Foo", foo).Debugf("[SERVER] Foo not ok")
			return grpc.Errorf(codes.InvalidArgument, "Foo not ok")
		}
		lastFoo = foo.GetFoo()

		err = stream.Send(&Bar{Bar: lastFoo})
		if err != nil {
			s.log.WithField("Method", "Sync").WithError(err).Debugf("[SERVER] Error")
			return err
		}
		s.SyncFoo = foo
	}
}
