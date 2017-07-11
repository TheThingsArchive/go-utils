// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"testing"

	"github.com/smartystreets/assertions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestGRPC(t *testing.T) {
	a := assertions.New(t)
	d := &ErrDescriptor{
		MessageFormat: "You do not have access to app with id {app_id}",
		Code:          77,
		Type:          PermissionDenied,
		registered:    true,
	}

	attributes := Attributes{
		"app_id": "foo",
		"count":  42,
	}

	err := d.New(attributes)

	code := GRPCCode(err)
	a.So(code, assertions.ShouldEqual, codes.PermissionDenied)

	// other errors should be unknown
	other := fmt.Errorf("Foo")
	code = GRPCCode(other)
	a.So(code, assertions.ShouldEqual, codes.Unknown)

	grpcErr := ToGRPC(err)

	got := FromGRPC(grpcErr)
	a.So(got.Code(), assertions.ShouldEqual, d.Code)
	a.So(got.Type(), assertions.ShouldEqual, d.Type)
	a.So(got.Error(), assertions.ShouldEqual, "You do not have access to app with id foo")

	a.So(got.Attributes(), assertions.ShouldNotBeEmpty)
	a.So(got.Attributes()["app_id"], assertions.ShouldResemble, attributes["app_id"])
	a.So(got.Attributes()["count"], assertions.ShouldAlmostEqual, attributes["count"])
}

func TestFromUnspecifiedGRPC(t *testing.T) {
	a := assertions.New(t)

	err := grpc.Errorf(codes.DeadlineExceeded, "This is an error")

	got := FromGRPC(err)
	a.So(got.Code(), assertions.ShouldEqual, 0)
	a.So(got.Type(), assertions.ShouldEqual, Timeout)
	a.So(got.Error(), assertions.ShouldEqual, "This is an error")
	a.So(got.Attributes(), assertions.ShouldBeNil)
}
