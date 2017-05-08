// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package rpclog

import (
	"net"
	"testing"

	. "github.com/smartystreets/assertions"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestFieldsFromContext(t *testing.T) {
	a := New(t)

	ctx := context.Background()

	addr, _ := net.ResolveTCPAddr("", "127.0.0.1:4242")
	ctx = peer.NewContext(ctx, &peer.Peer{
		Addr:     addr,
		AuthInfo: credentials.TLSInfo{},
	})

	fields := FieldsFromContext(ctx)
	a.So(fields, ShouldContainKey, "caller-ip")
	a.So(fields, ShouldContainKey, "auth-type")

	{
		ctx := metadata.NewIncomingContext(ctx, metadata.Pairs(
			"id", "id",
			"key", "key",
			"token", "",
			"service-version", "",
		))
		fields := FieldsFromContext(ctx)
		a.So(fields, ShouldContainKey, "id")
		a.So(fields, ShouldContainKey, "auth-type")
		a.So(fields["auth-type"], ShouldEqual, "tls,key")
		a.So(fields, ShouldNotContainKey, "service-version")
	}

	{
		ctx := metadata.NewOutgoingContext(ctx, metadata.Pairs(
			"id", "id",
			"token", "token",
		))
		fields := FieldsFromContext(ctx)
		a.So(fields, ShouldContainKey, "id")
		a.So(fields, ShouldContainKey, "auth-type")
		a.So(fields["auth-type"], ShouldEqual, "tls,token")
	}

}
