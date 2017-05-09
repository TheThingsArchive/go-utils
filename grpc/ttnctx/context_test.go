// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package ttnctx

import (
	"testing"

	. "github.com/smartystreets/assertions"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func TestContext(t *testing.T) {
	a := New(t)
	var err error

	// Errors if the context doesn't have valid metadata
	{
		ctx := context.Background()

		md := MetadataFromIncomingContext(ctx)
		a.So(md, ShouldHaveLength, 0)

		ctx = metadata.NewContext(ctx, metadata.Pairs())

		md = MetadataFromIncomingContext(ctx)
		a.So(md, ShouldHaveLength, 0)

		_, err = TokenFromIncomingContext(ctx)
		a.So(err, ShouldNotBeNil)

		_, err = KeyFromIncomingContext(ctx)
		a.So(err, ShouldNotBeNil)

		_, err = IDFromIncomingContext(ctx)
		a.So(err, ShouldNotBeNil)

		limit, offset, err := LimitAndOffsetFromIncomingContext(ctx)
		a.So(err, ShouldBeNil)
		a.So(limit, ShouldEqual, 0)
		a.So(offset, ShouldEqual, 0)

		serviceName, serviceVersion, netAddress, err := ServiceInfoFromIncomingContext(ctx)
		a.So(err, ShouldBeNil)
		a.So(serviceName, ShouldEqual, "")
		a.So(serviceVersion, ShouldEqual, "")
		a.So(netAddress, ShouldEqual, "")
	}

	// Errors if the context has wrong metadata
	{
		_, _, err := LimitAndOffsetFromIncomingContext(metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"limit", "wut",
		)))
		a.So(err, ShouldNotBeNil)

		_, _, err = LimitAndOffsetFromIncomingContext(metadata.NewIncomingContext(context.Background(), metadata.Pairs(
			"offset", "wut",
		)))
		a.So(err, ShouldNotBeNil)
	}

	{
		ctx := context.Background()

		ctx = OutgoingContextWithToken(ctx, "token")
		token, err := TokenFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(token, ShouldEqual, "token")

		ctx = OutgoingContextWithKey(ctx, "key")
		key, err := KeyFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(key, ShouldEqual, "key")

		ctx = OutgoingContextWithID(ctx, "id")
		id, err := IDFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(id, ShouldEqual, "id")

		ctx = OutgoingContextWithServiceInfo(ctx, "name", "version", "addr")
		serviceName, serviceVersion, netAddress, err := ServiceInfoFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(serviceName, ShouldEqual, "name")
		a.So(serviceVersion, ShouldEqual, "version")
		a.So(netAddress, ShouldEqual, "addr")

		ctx = OutgoingContextWithLimitAndOffset(ctx, 2, 4)
		limit, err := LimitFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(limit, ShouldEqual, 2)
		offset, err := OffsetFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(offset, ShouldEqual, 4)

		// Try the token again
		token, err = TokenFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(token, ShouldEqual, "token")
	}

	{
		ctx := OutgoingContextWithLimitAndOffset(metadata.NewContext(context.Background(), metadata.Pairs()), 0, 0)
		limit, err := LimitFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(limit, ShouldEqual, 0)
		offset, err := OffsetFromMetadata(MetadataFromOutgoingContext(ctx))
		a.So(err, ShouldBeNil)
		a.So(offset, ShouldEqual, 0)
	}

}
