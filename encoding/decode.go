package encoding

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func decodeToType(typ reflect.Kind, value string) interface{} {
	switch typ {
	case reflect.String:
		return value
	case reflect.Bool:
		v, _ := strconv.ParseBool(value)
		return v
	case reflect.Int:
		v, _ := strconv.ParseInt(value, 10, 64)
		return int(v)
	case reflect.Int8:
		return int8(decodeToType(reflect.Int, value).(int))
	case reflect.Int16:
		return int16(decodeToType(reflect.Int, value).(int))
	case reflect.Int32:
		return int32(decodeToType(reflect.Int, value).(int))
	case reflect.Int64:
		return int64(decodeToType(reflect.Int, value).(int))
	case reflect.Uint:
		v, _ := strconv.ParseUint(value, 10, 64)
		return uint(v)
	case reflect.Uint8:
		return uint8(decodeToType(reflect.Uint, value).(uint))
	case reflect.Uint16:
		return uint16(decodeToType(reflect.Uint, value).(uint))
	case reflect.Uint32:
		return uint32(decodeToType(reflect.Uint, value).(uint))
	case reflect.Uint64:
		return uint64(decodeToType(reflect.Uint, value).(uint))
	case reflect.Float64:
		v, _ := strconv.ParseFloat(value, 64)
		return v
	case reflect.Float32:
		return float32(decodeToType(reflect.Float64, value).(float64))
	}
	return nil
}

func unmarshalToType(typ reflect.Type, value string) (val interface{}, err error) {
	// If we get a pointer in, we'll return a pointer out
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	val = reflect.New(typ).Interface()
	defer func() {
		if err == nil && typ.Kind() != reflect.Ptr {
			val = reflect.Indirect(reflect.ValueOf(val)).Interface()
		}
	}()

	// Try Unmarshalers
	if um, ok := val.(encoding.TextUnmarshaler); ok {
		if err = um.UnmarshalText([]byte(value)); err == nil {
			return val, nil
		}
	}
	if um, ok := val.(json.Unmarshaler); ok {
		if err = um.UnmarshalJSON([]byte(value)); err == nil {
			return val, nil
		}
	}

	// Try JSON
	if err = json.Unmarshal([]byte(value), val); err == nil {
		return val, nil
	}

	// Return error if we have one
	if err != nil {
		return nil, err
	}

	return val, fmt.Errorf("No way to unmarshal \"%s\" to %s", value, typ.Name())
}

// FromStringStringMap decodes input into output with the same type as base. Only fields tagged by tagName get decoded. Optional argument properties specifies fields to decode.
func FromStringStringMap(tagName string, base interface{}, input map[string]string) (output interface{}, err error) {
	baseType := reflect.TypeOf(base)

	valType := baseType
	if baseType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}

	// If we get a pointer in, we'll return a pointer out
	valPtr := reflect.New(valType)
	val := valPtr.Elem()
	output = valPtr.Interface()

	defer func() {
		if err == nil && baseType.Kind() != reflect.Ptr {
			output = reflect.Indirect(reflect.ValueOf(output)).Interface()
		}
	}()

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		if field.Anonymous {
			continue
		}

		tagField, _ := parseTag(field.Tag.Get(tagName))
		if tagField == "" || tagField == "-" {
			continue
		}

		str, ok := input[tagField]
		if !ok || str == "" {
			continue
		}

		fieldType := field.Type
		fieldKind := field.Type.Kind()

		isPointerField := fieldKind == reflect.Ptr

		if isPointerField {
			if str == "null" {
				continue
			}
			fieldType = fieldType.Elem()
			fieldKind = fieldType.Kind()
		}

		var fieldVal interface{}

		switch fieldKind {
		case reflect.Struct, reflect.Array, reflect.Interface, reflect.Slice, reflect.Map:
			fieldVal, err = unmarshalToType(fieldType, str)
			if err != nil {
				return nil, err
			}
		default:
			fieldVal = decodeToType(fieldKind, str)
		}

		if isPointerField {
			fieldValPtr := reflect.New(fieldType)
			fieldValPtr.Elem().Set(reflect.ValueOf(fieldVal))
			val.Field(i).Set(fieldValPtr)
		} else {
			val.Field(i).Set(reflect.ValueOf(fieldVal))
		}
	}
	return output, nil
}
