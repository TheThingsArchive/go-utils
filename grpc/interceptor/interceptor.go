package interceptor

import (
	"time"

	"github.com/TheThingsNetwork/ttn/api"
	"github.com/TheThingsNetwork/ttn/utils/errors"
	"github.com/apex/log"
	context "golang.org/x/net/context" //TODO change to "context", when protoc supports it
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
)

// Unary generates a function usable as an argument to grpc.UnaryServerInterceptor
// fn should return appropriate log.Interface and request name to use in logging
func Unary(fn func(req interface{}, info *grpc.UnaryServerInfo) (log.Interface, string)) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log, reqStr := fn(req, info)

		fields := fieldsFromContext(ctx)
		fields["Method"] = info.FullMethod
		log = log.WithFields(fields)

		log.Debugf("%s started", reqStr)

		start := time.Now()
		resp, err = handler(ctx, req)
		log = withDurationSince(log, start)

		grpcErr := errors.BuildGRPCError(err)
		code := grpc.Code(grpcErr)
		log = log.WithField("Code", code)

		if grpcErr != nil {
			log = log.WithError(err)
		}
		log.Debugf("%s completed", reqStr)

		return resp, grpcErr
	}
}

// Stream generates a function usable as an argument to grpc.StreamServerInterceptor
// fn should return appropriate log.Interface and stream name to use in logging
func Stream(fn func(srv interface{}, info *grpc.StreamServerInfo) (log.Interface, string)) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log, streamStr := fn(srv, info)

		fields := fieldsFromContext(ss.Context())
		fields["Method"] = info.FullMethod
		log = log.WithFields(fields)

		log.Debugf("%s opened", streamStr)

		start := time.Now()
		err = handler(srv, ss)
		log = withDurationSince(log, start)

		grpcErr := errors.BuildGRPCError(err)
		code := grpc.Code(grpcErr)
		log = log.WithField("Code", code)

		if grpcErr != nil && code != codes.Canceled {
			log = log.WithError(err)
		}
		log.Debugf("%s closed", streamStr)

		return grpcErr
	}
}

func fieldsFromContext(ctx context.Context) log.Fields {
	fields := log.Fields{}

	if peer, ok := peer.FromContext(ctx); ok {
		fields["CallerIP"] = peer.Addr.String()

		if peer.AuthInfo != nil {
			fields["Auth-Type"] = peer.AuthInfo.AuthType()
		}
	}

	md, err := api.MetadataFromContext(ctx)
	if err != nil {
		return fields
	}

	if id, err := api.IDFromMetadata(md); err == nil {
		fields["ID"] = id
	}

	if offset, err := api.OffsetFromMetadata(md); err == nil && offset != 0 {
		fields["Offset"] = offset
	}

	if limit, err := api.LimitFromMetadata(md); err == nil && limit != 0 {
		fields["Limit"] = limit
	}

	return fields
}

func withDurationSince(log log.Interface, start time.Time) log.Interface {
	return log.WithField("Duration", time.Since(start))
}
