// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package rpcerror

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ConvertFunc converts error types
type ConvertFunc func(error) error

// UnaryServerInterceptor applies fn to errors returned by server.
func UnaryServerInterceptor(fn ConvertFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, fn(err)
	}
}

// StreamServerInterceptor applies fn to errors returned by server.
func StreamServerInterceptor(fn ConvertFunc) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		return fn(handler(srv, ss))
	}
}

// UnaryClientInterceptor applies fn to errors recieved by client.
func UnaryClientInterceptor(fn ConvertFunc) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		return fn(invoker(ctx, method, req, reply, cc, opts...))
	}
}

// StreamClientInterceptor applies fn to errors recieved by client.
func StreamClientInterceptor(fn ConvertFunc) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		stream, err = streamer(ctx, desc, cc, method, opts...)
		return stream, fn(err)
	}
}
