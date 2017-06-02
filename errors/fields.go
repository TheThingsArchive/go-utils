// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package errors

import (
	"reflect"

	"github.com/gogo/protobuf/types"
)

// Fields to attach to the error
type Fields map[string]interface{}

type hasFields interface {
	Fields() Fields
}

// Proto returns the fields as a protobuf
func (f Fields) Proto() (*types.Any, error) {
	s := &types.Struct{
		Fields: make(map[string]*types.Value),
	}
	for k, v := range f {
		switch v := v.(type) {
		case string:
			s.Fields[k] = &types.Value{Kind: &types.Value_StringValue{StringValue: v}}
		case bool:
			s.Fields[k] = &types.Value{Kind: &types.Value_BoolValue{BoolValue: v}}
		case int, int8, int16, int32, int64:
			s.Fields[k] = &types.Value{Kind: &types.Value_NumberValue{NumberValue: float64(reflect.ValueOf(v).Int())}}
		case uint, uint8, uint16, uint32, uint64:
			s.Fields[k] = &types.Value{Kind: &types.Value_NumberValue{NumberValue: float64(reflect.ValueOf(v).Uint())}}
		case float32, float64:
			s.Fields[k] = &types.Value{Kind: &types.Value_NumberValue{NumberValue: reflect.ValueOf(v).Float()}}
		}
	}
	p, err := types.MarshalAny(s)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// FieldsFromProto returns the fields from a protobuf
func FieldsFromProto(in *types.Any) (fields Fields, err error) {
	s := new(types.Struct)
	if !types.Is(in, s) {
		return nil, nil
	}
	err = types.UnmarshalAny(in, s)
	if err != nil {
		return nil, err
	}
	fields = make(Fields)
	for k, v := range s.Fields {
		switch v := v.Kind.(type) {
		case *types.Value_StringValue:
			fields[k] = v.StringValue
		case *types.Value_NumberValue:
			fields[k] = v.NumberValue
		case *types.Value_BoolValue:
			fields[k] = v.BoolValue
		}
	}
	return
}
