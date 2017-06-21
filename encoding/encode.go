// Copyright © 2017 The Things Network
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.

package encoding

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/structs"
)

type smap map[string]string

func (m smap) Set(key string, value string) {
	if _, ok := m[key]; ok {
		panic(fmt.Errorf("field names not unique (%s)", key))
	}
	m[key] = value
}

type imap map[string]interface{}

func (m imap) Set(key string, value interface{}) {
	if _, ok := m[key]; ok {
		panic(fmt.Errorf("field names not unique (%s)", key))
	}
	m[key] = value
}

func stringInSlice(search string, slice []string) bool {
	for _, i := range slice {
		if i == search {
			return true
		}
	}
	return false
}

type tagOptions []string

// Has returns true if opt is one of the options
func (t tagOptions) Has(opt string) bool {
	for _, hasOpt := range t {
		if hasOpt == opt || strings.HasPrefix(hasOpt, opt+"=") {
			return true
		}
	}
	return false
}

// Value returns value of tag option
func (t tagOptions) Value(opt string) string {
	for _, hasOpt := range t {
		if strings.HasPrefix(hasOpt, opt+"=") {
			return strings.TrimLeft(hasOpt, opt+"=")
		}
	}
	return ""
}

func parseTag(tag string) (string, tagOptions) {
	res := strings.Split(tag, ",")
	return res[0], res[1:]
}

type isZeroer interface {
	IsZero() bool
}

type isEmptier interface {
	IsEmpty() bool
}

// ToStringStringMap encodes fields tagged with tagName in input into map[string]string. Optional argument properties specifies fields to encode.
func ToStringStringMap(tagName string, input interface{}, properties ...string) (map[string]string, error) {
	vmap := smap(make(map[string]string))
	s := structs.New(input)
	s.TagName = tagName
	if len(properties) == 0 {
		properties = s.Names()
	}

	for _, field := range s.Fields() {
		if !field.IsExported() {
			continue
		}

		if !stringInSlice(field.Name(), properties) {
			continue
		}

		fieldName, opts := parseTag(field.Tag(tagName))
		squash, omitempty, include := opts.Has("squash"), opts.Has("omitempty"), opts.Has("include")
		if !squash && (fieldName == "" || fieldName == "-") {
			continue
		}

		val := field.Value()

		if omitempty {
			if field.IsZero() {
				continue
			}
			if z, ok := val.(isZeroer); ok && z.IsZero() {
				continue
			}
			if z, ok := val.(isEmptier); ok && z.IsEmpty() {
				continue
			}
		}

		kind := field.Kind()
		if kind == reflect.Ptr {
			v := reflect.ValueOf(val)
			if v.IsNil() {
				if fieldName != "" && fieldName != "-" {
					vmap.Set(fieldName, "")
				}
				continue
			}
			elem := v.Elem()
			kind = elem.Kind()
			val = elem.Interface()
		}

		if z, ok := val.(isZeroer); ok && z.IsZero() {
			vmap.Set(fieldName, "")
			continue
		}

		if z, ok := val.(isEmptier); ok && z.IsEmpty() {
			vmap.Set(fieldName, "")
			continue
		}

		if z, ok := val.(time.Time); ok {
			if z.Unix() == 0 {
				vmap.Set(fieldName, "")
				continue
			}
			val = z.UTC()
		}

		if (squash || include) && kind == reflect.Struct {
			var newProperties []string
			for _, prop := range properties {
				if strings.HasPrefix(prop, field.Name()+".") {
					newProperties = append(newProperties, strings.TrimPrefix(prop, field.Name()+"."))
				}
			}
			m, err := ToStringStringMap(tagName, val, newProperties...)
			if err != nil {
				return nil, err
			}

			var prefix string
			if !squash {
				prefix = fieldName + "."
			}

			for k, v := range m {
				vmap.Set(prefix+k, v)
			}
			continue
		}

		if v, ok := val.(string); ok {
			vmap.Set(fieldName, v)
			continue
		} else if v, ok := val.(*string); ok {
			vmap.Set(fieldName, *v)
			continue
		}

		if !field.IsZero() {
			if m, ok := val.(encoding.TextMarshaler); ok {
				txt, err := m.MarshalText()
				if err != nil {
					return nil, err
				}
				vmap.Set(fieldName, string(txt))
				continue
			}
			if m, ok := val.(json.Marshaler); ok {
				txt, err := m.MarshalJSON()
				if err != nil {
					return nil, err
				}
				vmap.Set(fieldName, string(txt))
				continue
			}
		}

		if kind == reflect.String {
			vmap.Set(fieldName, fmt.Sprint(val))
			continue
		}

		if txt, err := json.Marshal(val); err == nil {
			txt := string(txt)
			if txt == `""` || txt == "null" {
				vmap.Set(fieldName, "")
			} else {
				vmap.Set(fieldName, string(txt))
			}
			continue
		}

		vmap.Set(fieldName, fmt.Sprintf("%v", val))
	}
	return vmap, nil
}

// ToStringInterfaceMap encodes fields tagged with tagName in input into map[string]interface{}. Optional argument properties specifies fields to encode.
func ToStringInterfaceMap(tagName string, input interface{}, properties ...string) (map[string]interface{}, error) {
	vmap := imap(make(map[string]interface{}))

	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	t := val.Type()

	if len(properties) == 0 {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if _, ok := field.Tag.Lookup(tagName); ok {
				properties = append(properties, field.Name)
			}
		}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue
		}
		if !stringInSlice(field.Name, properties) {
			continue
		}

		fieldName, opts := parseTag(field.Tag.Get(tagName))
		squash, omitempty, include := opts.Has("squash"), opts.Has("omitempty"), opts.Has("include")
		if !squash && (fieldName == "" || fieldName == "-") {
			continue
		}

		fieldVal := val.Field(i)

		iface := fieldVal.Interface()
		if omitempty {
			if reflect.DeepEqual(iface, reflect.Zero(field.Type).Interface()) {
				continue
			}
			if z, ok := iface.(isZeroer); ok && z.IsZero() {
				continue
			}
			if z, ok := iface.(isEmptier); ok && z.IsEmpty() {
				continue
			}
		}

		kind := fieldVal.Kind()
		if kind == reflect.Ptr {
			if fieldVal.IsNil() {
				if fieldName != "" && fieldName != "-" {
					vmap.Set(fieldName, nil)
				}
				continue
			}
			fieldVal = fieldVal.Elem()
			kind = fieldVal.Kind()
		}

		if opts.Has("cast") {
			switch opts.Value("cast") {
			case "int64":
				var v int64
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			case "int32":
				var v int32
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			case "int16":
				var v int16
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			case "int8":
				var v int8
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			case "uint64":
				var v uint64
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			case "uint32":
				var v uint32
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			case "uint16":
				var v uint16
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			case "uint8":
				var v uint8
				fieldVal = fieldVal.Convert(reflect.TypeOf(v))
			default:
				panic(fmt.Errorf("Wrong cast type specified: %d", opts.Value("cast")))
			}
		}

		iface = fieldVal.Interface()

		if (squash || include) && kind == reflect.Struct {
			var newProperties []string
			for _, prop := range properties {
				if strings.HasPrefix(prop, field.Name+".") {
					newProperties = append(newProperties, strings.TrimPrefix(prop, field.Name+"."))
				}
			}
			m, err := ToStringInterfaceMap(tagName, iface, newProperties...)
			if err != nil {
				return nil, err
			}

			var prefix string
			if !squash {
				prefix = fieldName + "."
			}

			for k, v := range m {
				vmap.Set(prefix+k, v)
			}
			continue
		}

		vmap.Set(fieldName, iface)
	}
	return vmap, nil
}
