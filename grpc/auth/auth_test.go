// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package auth

import (
	"testing"

	"github.com/TheThingsNetwork/go-utils/grpc/ttnctx"
	. "github.com/smartystreets/assertions"
	"golang.org/x/net/context"
)

func TestAuth(t *testing.T) {
	a := New(t)
	var err error

	c := WithStaticToken("token")
	md, err := c.GetRequestMetadata(context.Background())
	a.So(err, ShouldBeNil)
	a.So(md, ShouldContainKey, "token")
	a.So(md["token"], ShouldEqual, "token")

	md, err = c.GetRequestMetadata(ttnctx.OutgoingContextWithToken(context.Background(), "existingtoken"))
	a.So(err, ShouldBeNil)
	a.So(md, ShouldContainKey, "token")
	a.So(md["token"], ShouldEqual, "existingtoken")

	a.So(c.RequireTransportSecurity(), ShouldBeTrue)
	a.So(c.WithInsecure().RequireTransportSecurity(), ShouldBeFalse)

	c = WithTokenFunc("id", func(id string) string {
		return id
	})

	md, err = c.GetRequestMetadata(ttnctx.OutgoingContextWithID(context.Background(), "id"))
	a.So(err, ShouldBeNil)
	a.So(md, ShouldContainKey, "token")
	a.So(md["token"], ShouldEqual, "id")
}
