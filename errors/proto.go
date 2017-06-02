// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes/any"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

type hasProto interface {
	Proto() *spb.Status
}

// Proto transforms the error into a Status protobuf
func Proto(err error) *spb.Status {
	if err, ok := err.(hasProto); ok {
		return err.Proto()
	}
	status := new(spb.Status)
	status.Code = int32(FindCode(err))
	status.Message = err.Error()
	if err, ok := err.(hasFields); ok && len(err.Fields()) > 0 {
		if fields, err := err.Fields().Proto(); err == nil {
			status.Details = append(status.Details, &any.Any{
				TypeUrl: fields.TypeUrl,
				Value:   fields.Value,
			})
		}
	}
	return status
}

// FromProto transforms the Status protobuf into an error
func FromProto(status *spb.Status) error {
	base := newBase(status.Message)
	base.code = codes.Code(status.Code)
	if status.Details != nil {
		for _, fields := range status.Details {
			fields, err := FieldsFromProto(&types.Any{
				TypeUrl: fields.TypeUrl,
				Value:   fields.Value,
			})
			if err == nil && fields != nil {
				base.fields = fields
				break
			}
		}
	}
	switch codes.Code(status.Code) {
	case codes.AlreadyExists:
		if err, ok := newErrAlreadyExistsFrom(status.Message, base.fields); ok {
			return err
		}
	case codes.InvalidArgument:
		if err, ok := newErrInvalidArgumentFrom(status.Message, base.fields); ok {
			return err
		}
	case codes.NotFound:
		if err, ok := newErrNotFoundFrom(status.Message, base.fields); ok {
			return err
		}
	case codes.PermissionDenied:
		if err, ok := newErrPermissionDeniedFrom(status.Message, base.fields); ok {
			return err
		}
	}
	return base
}
