// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"errors"
	"testing"

	s "github.com/smartystreets/assertions"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

func TestProto(t *testing.T) {
	a := s.New(t)

	var err error
	var p *spb.Status

	// Regular error
	err = errors.New("err")
	p = Proto(err)
	a.So(p.Message, s.ShouldEqual, "err")
	a.So(p.Code, s.ShouldEqual, codes.Unknown)

	// Rich error
	err = NewErrInvalidArgument("foo", "missing")
	p = Proto(err)
	a.So(p.Message, s.ShouldEqual, `invalid argument "foo": missing`)
	a.So(p.Code, s.ShouldEqual, codes.InvalidArgument)
	a.So(p.Details, s.ShouldNotBeEmpty)

	// Wrapped error
	err = Wrap(err, "unable to do something")
	p = Proto(err)
	a.So(p.Message, s.ShouldEqual, `unable to do something: invalid argument "foo": missing`)
	a.So(p.Code, s.ShouldEqual, codes.InvalidArgument)
	a.So(p.Details, s.ShouldNotBeEmpty)

	newErr := FromProto(p)
	a.So(newErr.Error(), s.ShouldEqual, `invalid argument "foo": missing`)

}
