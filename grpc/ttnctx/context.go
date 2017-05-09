// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package ttnctx

import (
	"errors"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// MetadataFromIncomingContext gets the metadata from the given context
func MetadataFromIncomingContext(ctx context.Context) metadata.MD {
	md, _ := metadata.FromIncomingContext(ctx)
	return md
}

// MetadataFromOutgoingContext gets the metadata from the given context
func MetadataFromOutgoingContext(ctx context.Context) metadata.MD {
	md, _ := metadata.FromOutgoingContext(ctx)
	return md
}

// TokenFromMetadata gets the token from the metadata or returns ErrNoToken
func TokenFromMetadata(md metadata.MD) (string, error) {
	token, ok := md["token"]
	if !ok || len(token) == 0 {
		return "", ErrNoToken
	}
	return token[0], nil
}

func outgoingContextWithMergedMetadata(ctx context.Context, kv ...string) context.Context {
	md := MetadataFromOutgoingContext(ctx)
	md = metadata.Join(metadata.Pairs(kv...), md)
	return metadata.NewOutgoingContext(ctx, md)
}

// TokenFromIncomingContext gets the token from the incoming context or returns ErrNoToken
func TokenFromIncomingContext(ctx context.Context) (string, error) {
	md := MetadataFromIncomingContext(ctx)
	return TokenFromMetadata(md)
}

// OutgoingContextWithToken returns an outgoing context with the token
func OutgoingContextWithToken(ctx context.Context, token string) context.Context {
	return outgoingContextWithMergedMetadata(ctx, "token", token)
}

// KeyFromMetadata gets the key from the metadata or returns ErrNoKey
func KeyFromMetadata(md metadata.MD) (string, error) {
	key, ok := md["key"]
	if !ok || len(key) == 0 {
		return "", ErrNoKey
	}
	return key[0], nil
}

// KeyFromIncomingContext gets the key from the incoming context or returns ErrNoKey
func KeyFromIncomingContext(ctx context.Context) (string, error) {
	md := MetadataFromIncomingContext(ctx)
	return KeyFromMetadata(md)
}

// OutgoingContextWithKey returns an outgoing context with the key
func OutgoingContextWithKey(ctx context.Context, key string) context.Context {
	return outgoingContextWithMergedMetadata(ctx, "key", key)
}

// IDFromMetadata gets the key from the metadata or returns ErrNoID
func IDFromMetadata(md metadata.MD) (string, error) {
	id, ok := md["id"]
	if !ok || len(id) == 0 {
		return "", ErrNoID
	}
	return id[0], nil
}

// IDFromIncomingContext gets the key from the incoming context or returns ErrNoID
func IDFromIncomingContext(ctx context.Context) (string, error) {
	md := MetadataFromIncomingContext(ctx)
	return IDFromMetadata(md)
}

// OutgoingContextWithID returns an outgoing context with the id
func OutgoingContextWithID(ctx context.Context, id string) context.Context {
	return outgoingContextWithMergedMetadata(ctx, "id", id)
}

// ServiceInfoFromMetadata gets the service information from the metadata or returns empty strings
func ServiceInfoFromMetadata(md metadata.MD) (serviceName, serviceVersion, netAddress string, err error) {
	serviceNameL, ok := md["service-name"]
	if ok && len(serviceNameL) > 0 {
		serviceName = serviceNameL[0]
	}
	serviceVersionL, ok := md["service-version"]
	if ok && len(serviceVersionL) > 0 {
		serviceVersion = serviceVersionL[0]
	}
	netAddressL, ok := md["net-address"]
	if ok && len(netAddressL) > 0 {
		netAddress = netAddressL[0]
	}
	return
}

// ServiceInfoFromIncomingContext gets the service information from the incoming context or returns empty strings
func ServiceInfoFromIncomingContext(ctx context.Context) (serviceName, serviceVersion, netAddress string, err error) {
	md := MetadataFromIncomingContext(ctx)
	return ServiceInfoFromMetadata(md)
}

// OutgoingContextWithServiceInfo returns an outgoing context with the id
func OutgoingContextWithServiceInfo(ctx context.Context, serviceName, serviceVersion, netAddress string) context.Context {
	return outgoingContextWithMergedMetadata(ctx, "service-name", serviceName, "service-version", serviceVersion, "net-address", netAddress)
}

// LimitFromMetadata gets the limit from the metadata
func LimitFromMetadata(md metadata.MD) (uint64, error) {
	limit, ok := md["limit"]
	if !ok || len(limit) == 0 {
		return 0, nil
	}
	return strconv.ParseUint(limit[0], 10, 64)
}

// OffsetFromMetadata gets the offset from the metadata
func OffsetFromMetadata(md metadata.MD) (uint64, error) {
	offset, ok := md["offset"]
	if !ok || len(offset) == 0 {
		return 0, nil
	}
	return strconv.ParseUint(offset[0], 10, 64)
}

// LimitAndOffsetFromIncomingContext gets the limit and offset from the incoming context
func LimitAndOffsetFromIncomingContext(ctx context.Context) (limit, offset uint64, err error) {
	md := MetadataFromIncomingContext(ctx)
	limit, err = LimitFromMetadata(md)
	if err != nil {
		return 0, 0, err
	}
	offset, err = OffsetFromMetadata(md)
	if err != nil {
		return 0, 0, err
	}
	return limit, offset, nil
}

// OutgoingContextWithLimitAndOffset returns an outgoing context with the limit and offset
func OutgoingContextWithLimitAndOffset(ctx context.Context, limit, offset uint64) context.Context {
	var pairs []string
	if limit != 0 {
		pairs = append(pairs, "limit", strconv.FormatUint(limit, 10))
	}
	if offset != 0 {
		pairs = append(pairs, "offset", strconv.FormatUint(offset, 10))
	}
	if len(pairs) == 0 {
		return ctx
	}
	return outgoingContextWithMergedMetadata(ctx, pairs...)
}

// Errors that are returned when an item could not be retrieved
// TODO: use go-utils/errors when ready
var (
	ErrNoToken = errors.New("context: Metadata does not contain token")
	ErrNoKey   = errors.New("context: Metadata does not contain key")
	ErrNoID    = errors.New("context: Metadata does not contain id")
)
