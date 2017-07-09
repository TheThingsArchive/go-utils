// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package rpcerror

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/smartystreets/assertions"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var inputErr = errors.New("input")
var outputErr = errors.New("output")

func convert(err error) error {
	if err != inputErr {
		panic(fmt.Sprintf("Wrong error specified: %s", err))
	}
	return outputErr
}

func TestUnaryServerInterceptor(t *testing.T) {
	a := New(t)
	ctx := context.Background()
	req := &struct{}{}
	_, err := UnaryServerInterceptor(convert)(ctx, req, nil, func(fnCtx context.Context, fnReq interface{}) (interface{}, error) {
		a.So(fnCtx, ShouldEqual, ctx)
		a.So(fnReq, ShouldEqual, req)
		return nil, inputErr
	})
	a.So(err, ShouldEqual, outputErr)
}

func TestStreamServerInterceptor(t *testing.T) {
	a := New(t)
	ctx := context.Background()
	var ss grpc.ServerStream = nil
	err := StreamServerInterceptor(convert)(ctx, ss, nil, func(fnCtx interface{}, fnSs grpc.ServerStream) error {
		a.So(fnCtx, ShouldEqual, ctx)
		a.So(fnSs, ShouldEqual, ss)
		return inputErr
	})
	a.So(err, ShouldEqual, outputErr)
}

func TestUnaryClientInterceptor(t *testing.T) {
	a := New(t)
	ctx := context.Background()
	method := "test"
	req := &struct{}{}
	reply := &struct{}{}
	cc := &grpc.ClientConn{}
	opts := []grpc.CallOption{&grpc.EmptyCallOption{}, &grpc.EmptyCallOption{}}
	err := UnaryClientInterceptor(convert)(ctx, method, req, reply, cc, func(fnCtx context.Context, fnMethod string, fnReq, fnReply interface{}, fnCc *grpc.ClientConn, fnOpts ...grpc.CallOption) error {
		a.So(fnCtx, ShouldEqual, ctx)
		a.So(fnMethod, ShouldEqual, method)
		a.So(fnReq, ShouldEqual, req)
		a.So(fnReply, ShouldEqual, reply)
		a.So(fnCc, ShouldEqual, cc)
		a.So(fnOpts, ShouldResemble, opts)
		return inputErr
	}, opts...)
	a.So(err, ShouldEqual, outputErr)
}

func TestStreamClientInterceptor(t *testing.T) {
	a := New(t)
	ctx := context.Background()
	desc := &grpc.StreamDesc{}
	cc := &grpc.ClientConn{}
	method := "test"
	opts := []grpc.CallOption{&grpc.EmptyCallOption{}, &grpc.EmptyCallOption{}}
	_, err := StreamClientInterceptor(convert)(ctx, desc, cc, method, func(fnCtx context.Context, fnDesc *grpc.StreamDesc, fnCc *grpc.ClientConn, fnMethod string, fnOpts ...grpc.CallOption) (grpc.ClientStream, error) {
		a.So(fnCtx, ShouldEqual, ctx)
		a.So(fnDesc, ShouldEqual, desc)
		a.So(fnCc, ShouldEqual, cc)
		a.So(fnMethod, ShouldEqual, method)
		a.So(fnOpts, ShouldResemble, opts)
		return nil, inputErr
	}, opts...)
	a.So(err, ShouldEqual, outputErr)
}
