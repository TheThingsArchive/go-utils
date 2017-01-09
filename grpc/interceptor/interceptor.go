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

func Unary(fn func(req interface{}, info *grpc.UnaryServerInfo) (log.Interface, string)) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log, reqStr := fn(req, info)

		fields := fieldsFromContext(ctx)
		fields["method"] = info.FullMethod
		log = log.WithFields(fields)

		log.Debugf("received %s", reqStr)

		start := time.Now()
		resp, err = handler(ctx, req)
		log = withDurationSince(log, start)

		grpcErr := errors.BuildGRPCError(err)
		code := grpc.Code(grpcErr)
		log = log.WithField("code", code)

		if grpcErr != nil {
			log.WithError(err).Errorf("%s failed", reqStr)
		} else {
			log.Debugf("%s completed", reqStr)
		}
		return resp, grpcErr
	}
}

func Stream(fn func(srv interface{}, info *grpc.StreamServerInfo) (log.Interface, string)) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		log, streamStr := fn(srv, info)

		fields := fieldsFromContext(ss.Context())
		fields["method"] = info.FullMethod
		log = log.WithFields(fields)

		log.Debugf("opening a new %s", streamStr)

		start := time.Now()
		err = handler(srv, ss)
		log = withDurationSince(log, start)

		grpcErr := errors.BuildGRPCError(err)
		code := grpc.Code(grpcErr)
		log = log.WithField("code", code)

		if grpcErr != nil && code != codes.Canceled {
			log.WithError(err).Errorf("%s failed", streamStr)
		} else {
			log.Debugf("%s closed", streamStr)
		}

		return grpcErr
	}
}

func fieldsFromContext(ctx context.Context) log.Fields {
	fields := log.Fields{}

	if peer, ok := peer.FromContext(ctx); ok {
		fields["ip"] = peer.Addr.String()

		if peer.AuthInfo != nil {
			fields["auth-type"] = peer.AuthInfo.AuthType()
		}
	}

	md, err := api.MetadataFromContext(ctx)
	if err != nil {
		return fields
	}

	if id, err := api.IDFromMetadata(md); err == nil {
		fields["id"] = id
	}

	if offset, err := api.OffsetFromMetadata(md); err == nil && offset != 0 {
		fields["offset"] = offset
	}

	if limit, err := api.LimitFromMetadata(md); err == nil && limit != 0 {
		fields["limit"] = limit
	}

	return fields
}

func withDurationSince(log log.Interface, start time.Time) log.Interface {
	return log.WithField("duration", time.Since(start))
}
