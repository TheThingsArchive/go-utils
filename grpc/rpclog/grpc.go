package rpclog

import (
	"time"

	"github.com/TheThingsNetwork/go-utils/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ServerOptions for logging RPCs
func ServerOptions(log log.Interface) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(UnaryServerInterceptor(log)),
		grpc.StreamInterceptor(StreamServerInterceptor(log)),
	}
}

// ClientOptions for logging RPCs
func ClientOptions(log log.Interface) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(log)),
		grpc.WithStreamInterceptor(StreamClientInterceptor(log)),
	}
}

// UnaryServerInterceptor logs unary RPCs on the server side
func UnaryServerInterceptor(log log.Interface) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log := log.WithField("method", info.FullMethod)
		log = log.WithFields(FieldsFromContext(ctx))
		start := time.Now()
		resp, err = handler(ctx, req)
		log = log.WithField("duration", time.Since(start))
		log.Debug("Request done")
		return
	}
}

// StreamServerInterceptor logs streaming RPCs on the server side
func StreamServerInterceptor(log log.Interface) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log := log.WithField("method", info.FullMethod)
		log = log.WithFields(FieldsFromContext(ss.Context()))
		start := time.Now()
		log.Debug("Server stream starting")
		err = handler(srv, ss)
		log = log.WithField("duration", time.Since(start))
		log.Debug("Server stream done")
		return
	}
}

// UnaryClientInterceptor logs unary RPCs on the client side
func UnaryClientInterceptor(log log.Interface) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		log := log.WithField("method", method)
		log = log.WithFields(FieldsFromContext(ctx))
		start := time.Now()
		err = invoker(ctx, method, req, reply, cc, opts...)
		log = log.WithField("duration", time.Since(start))
		log.Debug("Request done")
		return
	}
}

// StreamClientInterceptor logs streaming RPCs on the client side
func StreamClientInterceptor(log log.Interface) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		log := log.WithField("method", method)
		log = log.WithFields(FieldsFromContext(ctx))
		log.Debug("Client stream starting")
		stream, err = streamer(ctx, desc, cc, method, opts...)
		go func() {
			<-stream.Context().Done()
			log.Debug("Client stream done")
		}()
		return
	}
}
