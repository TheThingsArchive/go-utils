// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package rpclog

import (
	"time"

	ttnlog "github.com/TheThingsNetwork/go-utils/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func getLog(log ttnlog.Interface) ttnlog.Interface {
	if log == nil {
		return ttnlog.Get()
	}
	return log
}

// ServerOptions for logging RPCs
func ServerOptions(log ttnlog.Interface) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(UnaryServerInterceptor(log)),
		grpc.StreamInterceptor(StreamServerInterceptor(log)),
	}
}

// ClientOptions for logging RPCs
func ClientOptions(log ttnlog.Interface) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(log)),
		grpc.WithStreamInterceptor(StreamClientInterceptor(log)),
	}
}

// UnaryServerInterceptor logs unary RPCs on the server side
func UnaryServerInterceptor(log ttnlog.Interface) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := getLog(log).WithField("method", info.FullMethod)
		log = log.WithFields(FieldsFromIncomingContext(ctx))
		start := time.Now()
		resp, err = handler(ctx, req)
		log = log.WithField("duration", time.Since(start))
		if err != nil {
			log.WithError(err).Debug("rpc-server: call failed")
			return
		}
		log.Debug("rpc-server: call done")
		return
	}
}

// StreamServerInterceptor logs streaming RPCs on the server side
func StreamServerInterceptor(log ttnlog.Interface) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log := getLog(log).WithField("method", info.FullMethod)
		log = log.WithFields(FieldsFromIncomingContext(ss.Context()))
		start := time.Now()
		log.Debug("rpc-server: stream starting")
		err = handler(srv, ss)
		log = log.WithField("duration", time.Since(start))
		if err != nil {
			if err == context.Canceled || grpc.Code(err) == codes.Canceled {
				log.Debug("rpc-server: stream canceled")
				return
			}
			log.WithError(err).Debug("rpc-server: stream failed")
			return
		}
		log.Debug("rpc-server: stream done")
		return
	}
}

// UnaryClientInterceptor logs unary RPCs on the client side
func UnaryClientInterceptor(log ttnlog.Interface) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		log := getLog(log).WithField("method", method)
		log = log.WithFields(FieldsFromOutgoingContext(ctx))
		start := time.Now()
		err = invoker(ctx, method, req, reply, cc, opts...)
		log = log.WithField("duration", time.Since(start))
		if err != nil {
			log.WithError(err).Debug("rpc-client: call failed")
			return
		}
		log.Debug("rpc-client: call done")
		return
	}
}

// StreamClientInterceptor logs streaming RPCs on the client side
func StreamClientInterceptor(log ttnlog.Interface) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		log := getLog(log).WithField("method", method)
		log = log.WithFields(FieldsFromOutgoingContext(ctx))
		log.Debug("rpc-client: stream starting")
		stream, err = streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			if err == context.Canceled || grpc.Code(err) == codes.Canceled {
				log.Debug("rpc-client: stream canceled")
				return
			}
			log.WithError(err).Debug("rpc-client: stream failed")
			return
		}
		go func() {
			<-stream.Context().Done()
			if err := stream.Context().Err(); err != nil {
				log = log.WithError(err)
			}
			log.Debug("rpc-client: stream done")
		}()
		return
	}
}
